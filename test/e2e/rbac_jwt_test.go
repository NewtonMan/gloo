package e2e_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/solo-io/gloo/test/services/envoy"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/test/ginkgo/parallel"

	errors "github.com/rotisserie/eris"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/fgrosse/zaptest"
	"github.com/golang-jwt/jwt"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	jwtplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	glootransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gopkg.in/square/go-jose.v2"
)

var (
	baseJwksPort = uint32(28000)
)

const (
	issuer   = "issuer"
	audience = "thats-us"

	admin            = "admin"
	editor           = "editor"
	user             = "user"
	noDelimiterAdmin = "noDelimiterAdmin"
)

func jwks(ctx context.Context) (uint32, *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey) {
	rsaPriv, err := rsa.GenerateKey(rand.Reader, 512)
	Expect(err).NotTo(HaveOccurred())
	ecdsaPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Expect(err).NotTo(HaveOccurred())
	ed25519Pub, ed25519Priv, err := ed25519.GenerateKey(rand.Reader)
	Expect(err).NotTo(HaveOccurred())
	keySet := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{{
			Key:       rsaPriv.Public(),
			Algorithm: "RS256",
			Use:       "sig",
		}, {
			Key:       ecdsaPriv.Public(),
			Algorithm: "ES256",
			Use:       "sig",
		}, {
			Key:       ed25519Pub,
			Algorithm: "EdDSA",
			Use:       "sig",
		}},
	}

	jwksBytes, err := json.Marshal(keySet)
	Expect(err).NotTo(HaveOccurred())

	jwksPort := parallel.AdvancePortSafeListen(&baseJwksPort)
	jwtHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Write(jwksBytes)
	}
	h2s := &http2.Server{}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", jwksPort),
		Handler: h2c.NewHandler(http.HandlerFunc(jwtHandler), h2s),
	}

	go s.ListenAndServe()
	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	// serialize json and show
	return jwksPort, rsaPriv, ecdsaPriv, ed25519Priv
}

func getToken(claims jwt.Claims, key interface{}, method jwt.SigningMethod) string {
	var s string
	var err error
	switch key.(type) {
	case *rsa.PublicKey:
		s, err = jwt.NewWithClaims(method, claims).SignedString(key.(*rsa.PublicKey))
	case *ecdsa.PublicKey:
		s, err = jwt.NewWithClaims(method, claims).SignedString(key.(*ecdsa.PublicKey))
	case *ed25519.PublicKey:
		s, err = jwt.NewWithClaims(method, claims).SignedString(key.(*ed25519.PublicKey))
	default:
		err = eris.New("Unsupported token type")
	}
	s, err = jwt.NewWithClaims(method, claims).SignedString(key)
	Expect(err).NotTo(HaveOccurred())
	return s
}

func getMapToken(claims jwt.MapClaims, key *rsa.PrivateKey) string {
	s, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	ExpectWithOffset(2, err).NotTo(HaveOccurred())
	return s
}

var _ = Describe("JWT_RBAC", func() {

	var (
		ctx               context.Context
		cancel            context.CancelFunc
		testClients       services.TestClients
		jwksPort          uint32
		rsaPrivateKey     *rsa.PrivateKey
		ecdsaPrivateKey   *ecdsa.PrivateKey
		ed25519PrivateKey ed25519.PrivateKey
		jwtksServerRef    *core.ResourceRef
		envoyInstance     *envoy.Instance
		testUpstream      *v1helpers.TestUpstream
		envoyPort         uint32
	)

	BeforeEach(func() {
		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		jwksPort, rsaPrivateKey, ecdsaPrivateKey, ed25519PrivateKey = jwks(ctx)

		envoyInstance = envoyFactory.NewInstance()
		envoyPort = envoyInstance.HttpPort

		jwksServer := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "jwks-server",
				Namespace: "default",
			},
			UseHttp2: &wrappers.BoolValue{Value: true},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: envoyInstance.GlooAddr,
						Port: jwksPort,
					}},
				},
			},
		}

		testClients.UpstreamClient.Write(jwksServer, clients.WriteOpts{})
		jwtksServerRef = jwksServer.Metadata.Ref()
		rbacSettings := &rbac.Settings{
			RequireRbac: true,
		}

		settings := &gloov1.Settings{
			Gateway: &gloov1.GatewayOptions{
				Validation: &gloov1.GatewayOptions_ValidationOptions{
					// To validate transformations, we call out to an Envoy binary running in validate mode
					// https://github.com/solo-io/gloo/blob/01d04751f72c168e304977c4f67fdbcbf30232a9/projects/gloo/pkg/bootstrap/bootstrap_validation.go#L28
					// This binary is present in our CI/CD pipeline. But when running locally it is not, so we fallback to the Upstream Envoy binary
					// which doesn't have the custom Solo.io types registered with the deserializer. Therefore, when running locally tests will fail,
					// and the logs will contain:
					//	"Invalid type URL, unknown type: envoy.api.v2.filter.http.RouteTransformations for type Any)"
					// We do not perform transformation validation as part of our in memory e2e tests, so we explicitly disable this
					DisableTransformationValidation: &wrappers.BoolValue{
						Value: true,
					},
				},
			},
			Rbac: rbacSettings,
		}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, Settings: settings})

		err := envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		var opts clients.WriteOpts
		up := testUpstream.Upstream
		_, err = testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())

		helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
			return testClients.UpstreamClient.Read(up.GetMetadata().GetNamespace(), up.GetMetadata().GetName(), clients.ReadOpts{})
		})
	})

	AfterEach(func() {
		envoyInstance.Clean()
		cancel()
	})

	ExpectAccess := func(bar, fooget, foopost int, getBookRecommendations int, getVerifiedEmail int, augmentRequest func(*http.Request)) {
		// nestedFunctionLevel is for callstack reporting and should end up in the calling function.
		// here we are in expected access and then in the testquery function
		const nestedFunctionLevel = 2
		testQuery := func(method, path string, expectedStatus int, eventually bool) {
			url := fmt.Sprintf("http://%s:%d%s", "localhost", envoyPort, path)
			By("Querying " + url)
			req, err := http.NewRequest(method, url, nil)
			Expect(err).NotTo(HaveOccurred())
			augmentRequest(req)

			if eventually {
				EventuallyWithOffset(nestedFunctionLevel, func() (int, error) {
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(expectedStatus))
			} else {
				resp, err := http.DefaultClient.Do(req)
				ExpectWithOffset(nestedFunctionLevel, err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				_, _ = io.ReadAll(resp.Body)
				ExpectWithOffset(nestedFunctionLevel, resp.StatusCode).To(Equal(expectedStatus))
			}
		}

		// test public route in eventually to let the proxy time to start
		testQuery("GET", "/public_route", http.StatusOK, true)

		// No need to do eventually here as all is initialized.
		testQuery("GET", "/private_route", http.StatusForbidden, false)

		testQuery("GET", "/bar", bar, false)

		testQuery("GET", "/foo", fooget, false)

		testQuery("POST", "/foo", foopost, false)

		// These endpoints are only for those with advanced nested claims, -1 to skip
		if getBookRecommendations != -1 {
			testQuery("GET", "/book-recommendations", getBookRecommendations, false)
		}
		if getVerifiedEmail != -1 {
			testQuery("GET", "/verified-email", getVerifiedEmail, false)
		}
	}

	getClaims := func(sub string) jwt.StandardClaims {
		return jwt.StandardClaims{
			Issuer:   issuer,
			Audience: audience,
			Subject:  sub,
		}
	}
	getRsaTokenFor := func(sub string) string {
		tok := getToken(getClaims(sub), rsaPrivateKey, jwt.SigningMethodRS256)
		By("using token " + tok)
		return tok
	}

	getAdvancedClaimTokenFor := func(sub string, emailVerified bool, hobbies []string) string {
		claims := jwt.MapClaims{
			"iss": issuer,
			"aud": audience,
			"sub": sub,
			"metadata": map[string]interface{}{
				"foo": map[string]interface{}{
					"role": sub,
				},
				"email_verified": emailVerified,
				"hobbies":        hobbies,
			},
		}
		tok := getToken(claims, rsaPrivateKey, jwt.SigningMethodRS256)
		By("using token " + tok)
		return tok
	}

	getDefaultAdvancedClaimTokenFor := func(sub string) string {
		return getAdvancedClaimTokenFor(sub, true, []string{"long walks", "reading", "writing e2e tests"})
	}

	getMapTokenFor := func(sub string) string {
		claims := jwt.MapClaims{
			"iss": issuer,
			"aud": audience,
			"sub": sub,
			"data": map[string]string{
				"name": "test",
			},
		}
		tok := getMapToken(claims, rsaPrivateKey)
		By("using token " + tok)
		return tok
	}

	addBearer := func(req *http.Request, token string) {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	addToken := func(req *http.Request, sub string) {
		addBearer(req, getRsaTokenFor(sub))
	}
	addAdvancedClaimToken := func(req *http.Request, sub string, emailVerified bool, hobbies []string) {
		addBearer(req, getAdvancedClaimTokenFor(sub, emailVerified, hobbies))
	}
	addDefaultAdvancedClaimToken := func(req *http.Request, sub string) {
		addBearer(req, getDefaultAdvancedClaimTokenFor(sub))
	}

	Context("jwt tests", func() {
		BeforeEach(func() {
			proxy := getProxyJwt(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})

			// wait for key service to start
			waitForKeyService(jwksPort)
		})

		Context("forward token", func() {
			Context("should forward token upsteam", func() {
				testWithToken := func(token string) {
					Eventually(func() (int, error) {
						url := fmt.Sprintf("http://%s:%d/authnonly", "localhost", envoyPort)
						By("Querying " + url)
						req, err := http.NewRequest("GET", url, nil)
						Expect(err).NotTo(HaveOccurred())
						req.Header.Add("x-jwt", "JWT "+token)
						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return 0, err
						}
						defer resp.Body.Close()
						_, _ = io.ReadAll(resp.Body)
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))

					select {
					case received := <-testUpstream.C:
						Expect(received.Headers).To(HaveKeyWithValue("X-Jwt", []string{"JWT " + token}))
					default:
						Fail("request didnt make it upstream")
					}
				}

				// Test has the additional purpose of checking support for rsa, ecdsa, and ed25519 encodings
				It("rsa token", func() {
					testWithToken(getRsaTokenFor("user"))
				})
				It("ecdsa token", func() {
					testWithToken(getToken(getClaims("user"), ecdsaPrivateKey, jwt.SigningMethodES256))
				})
				It("ed25519 token", func() {
					testWithToken(getToken(getClaims("user"), ed25519PrivateKey, jwt.SigningMethodEdDSA))
				})
			})

		})
		Context("token source", func() {

			It("should get token from custom header", func() {
				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/authnonly", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					token := getRsaTokenFor("user")
					req.Header.Add("x-jwt", "JWT "+token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})
			It("should get token from custom query param", func() {
				Eventually(func() (int, error) {
					token := getRsaTokenFor("user")

					url := fmt.Sprintf("http://%s:%d/authnonly?jwttoken="+token, "localhost", envoyPort)
					By("Querying " + url)
					resp, err := http.Get(url)
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})
		})

		Context("claims to headers", func() {
			It("should should move the sub claim to a header", func() {
				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/authnonly", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					token := getRsaTokenFor("user")
					req.Header.Add("x-jwt", "JWT "+token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))

				select {
				case received := <-testUpstream.C:
					Expect(received.Headers).To(HaveKeyWithValue("X-Sub", []string{"user", "user"}))
				default:
					Fail("request didnt make it upstream")
				}
			})
			It("should re-route based on the new header added", func() {
				Eventually(func() (int, error) {
					// test with nested claim in token
					token := getMapTokenFor("teatime")
					url := fmt.Sprintf("http://%s:%d/authnonly?jwttoken=%s", "localhost", envoyPort, token)
					By("Querying " + url)
					resp, err := http.Get(url)
					if err != nil {
						return 0, err
					}
					defer resp.Body.Close()
					_, _ = io.ReadAll(resp.Body)
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))

				select {
				case received := <-testUpstream.C:
					Expect(received.Headers).To(HaveKeyWithValue("X-New-Header", []string{"new"}))
				default:
					Fail("request didnt make it upstream")
				}
			})
		})

	})
	Context("user access tests", func() {

		BeforeEach(func() {

			// paranoid check. We do this each time as someof the handles
			// are reset by a higher level beforeeach
			proxy := getProxyJwtRbac(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())
			_ = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})

			proxy = getProxyJwtRbac(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})

			// wait for key service to start
			waitForKeyService(jwksPort)

		})

		Context("non admin user", func() {
			It("should allow non admin user access to GET foo", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusForbidden, -1, -1,
					func(req *http.Request) { addToken(req, "user") })
			})

		})

		Context("editor user", func() {
			It("should allow most things", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusOK, -1, -1,
					func(req *http.Request) { addToken(req, "editor") })
			})
		})

		Context("admin user", func() {
			It("should allow everything", func() {
				ExpectAccess(http.StatusOK, http.StatusOK, http.StatusOK, -1, -1,
					func(req *http.Request) { addToken(req, "admin") })
			})
		})

		Context("anonymous user", func() {
			It("should only allow public route", func() {
				ExpectAccess(http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized, -1, -1,
					func(req *http.Request) {})
			})
		})

		Context("bad token user", func() {
			It("should only allow public route", func() {
				token := getRsaTokenFor("admin")
				// remove some stuff to make the signature invalid
				badToken := token[:len(token)-10]
				ExpectAccess(http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized, -1, -1,
					func(req *http.Request) { addBearer(req, badToken) })
			})
		})

	})

	Context("User access with nested claims", func() {

		BeforeEach(func() {

			// paranoid check. We do this each time as someof the handles
			// are reset by a higher level beforeeach
			proxy := getProxyJwtRbacNestedClaims(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())
			_ = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})

			proxy = getProxyJwtRbacNestedClaims(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})
			// wait for key service to start
			waitForKeyService(jwksPort)
		})

		Context("non admin user", func() {
			It("should allow non admin user access to GET foo", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusForbidden, http.StatusOK, http.StatusOK,
					func(req *http.Request) { addDefaultAdvancedClaimToken(req, "user") })
			})

		})

		Context("editor user", func() {
			It("should allow most things", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK,
					func(req *http.Request) { addDefaultAdvancedClaimToken(req, "editor") })
			})
		})

		Context("admin user", func() {
			It("should allow everything", func() {
				ExpectAccess(http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK, http.StatusOK,
					func(req *http.Request) { addDefaultAdvancedClaimToken(req, "admin") })
			})
		})

		Context("anonymous user", func() {
			It("should only allow public route", func() {
				ExpectAccess(http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized,
					func(req *http.Request) {})
			})
		})

		Context("bad token user", func() {
			It("should only allow public route", func() {
				token := getRsaTokenFor("admin")
				// remove some stuff to make the signature invalid
				badToken := token[:len(token)-10]
				ExpectAccess(http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized,
					func(req *http.Request) { addBearer(req, badToken) })
			})
		})

		Context("noDelimiterAdmin user", func() {
			// Without the nestedClaimDelimiter the "role" claim value should never be found,
			// because the matcher is not looking for a path, it's looking for a top-level
			// claim named "metadata.foo.role"
			It("should deny everything", func() {
				ExpectAccess(http.StatusForbidden, http.StatusForbidden, http.StatusForbidden,
					http.StatusOK, // there is a delimiter on the book-recommendations policy, so that one should be 200
					http.StatusOK, // there is a delimiter on the verified-email policy, so that one should be 200
					func(req *http.Request) { addDefaultAdvancedClaimToken(req, "noDelimiterAdmin") })
			})
		})

		// Tests ClaimMatcher.LIST_CONTAINS
		Context("users that don't like to read", func() {
			It("should not have access to book recommendations", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusForbidden, http.StatusForbidden, http.StatusOK,
					func(req *http.Request) {
						addAdvancedClaimToken(req, "user", true, []string{"long walks", "writing e2e tests"})
					})
			})
		})

		// Tests ClaimMatcher.BOOLEAN
		Context("non-verified emails", func() {
			It("should not have access to /verified-email", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusForbidden, http.StatusOK, http.StatusForbidden,
					func(req *http.Request) {
						addAdvancedClaimToken(req, "user", false, []string{"long walks", "reading", "writing e2e tests"})
					})
			})
		})
	})
})

// Essentially the same as getProxyJwtRbac, but requires a "metadata.foo.role"
// nested claim, rather than the "iss" and "sub" claims.
func getProxyJwtRbacNestedClaims(envoyPort uint32, jwtksServerRef, upstream *core.ResourceRef) *gloov1.Proxy {
	jwtCfg := &jwtplugin.VhostExtension{
		Providers: map[string]*jwtplugin.Provider{
			"testprovider": {
				Jwks: &jwtplugin.Jwks{
					Jwks: &jwtplugin.Jwks_Remote{
						Remote: &jwtplugin.RemoteJwks{
							Url:         "http://test/keys",
							UpstreamRef: jwtksServerRef,
						},
					},
				},
				Audiences: []string{audience},
				Issuer:    issuer,
			}},
	}

	rbacCfg := &rbac.ExtensionSettings{
		Policies: map[string]*rbac.Policy{
			"user": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"metadata.foo.role": user,
						},
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/foo",
					Methods:    []string{"GET"},
				},
				NestedClaimDelimiter: ".",
			},
			"editor": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"metadata.foo.role": editor,
						},
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/foo",
					Methods:    []string{"GET", "POST"},
				},
				NestedClaimDelimiter: ".",
			},
			"admin": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"metadata.foo.role": admin,
						},
					},
				}},
				Permissions:          &rbac.Permissions{},
				NestedClaimDelimiter: ".",
			},
			"noDelimiterAdmin": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"metadata.foo.role": noDelimiterAdmin,
						},
					},
				}},
				Permissions: &rbac.Permissions{},
			},
			"book-recommendations": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"metadata.hobbies": "reading",
						},
						Matcher: rbac.JWTPrincipal_LIST_CONTAINS,
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/book-recommendations",
					Methods:    []string{"GET"},
				},
				NestedClaimDelimiter: ".",
			},
			"verified-email": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"metadata.email_verified": "true",
						},
						Matcher: rbac.JWTPrincipal_BOOLEAN,
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/verified-email",
					Methods:    []string{"GET"},
				},
				NestedClaimDelimiter: ".",
			},
		},
	}

	return getProxyJwtRbacWithExtensions(envoyPort, jwtksServerRef, upstream, jwtCfg, rbacCfg)
}

func getProxyJwtRbac(envoyPort uint32, jwtksServerRef, upstream *core.ResourceRef) *gloov1.Proxy {

	jwtCfg := &jwtplugin.VhostExtension{
		Providers: map[string]*jwtplugin.Provider{
			"testprovider": {
				Jwks: &jwtplugin.Jwks{
					Jwks: &jwtplugin.Jwks_Remote{
						Remote: &jwtplugin.RemoteJwks{
							Url:         "http://test/keys",
							UpstreamRef: jwtksServerRef,
						},
					},
				},
				Audiences: []string{audience},
				Issuer:    issuer,
			}},
	}

	rbacCfg := &rbac.ExtensionSettings{
		Policies: map[string]*rbac.Policy{
			"user": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"iss": issuer,
							"sub": user,
						},
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/foo",
					Methods:    []string{"GET"},
				},
			},
			"editor": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"iss": issuer,
							"sub": editor,
						},
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/foo",
					Methods:    []string{"GET", "POST"},
				},
			},
			"admin": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"iss": issuer,
							"sub": admin,
						},
					},
				}},
				Permissions: &rbac.Permissions{},
			},
		},
	}

	return getProxyJwtRbacWithExtensions(envoyPort, jwtksServerRef, upstream, jwtCfg, rbacCfg)
}

func getProxyJwt(envoyPort uint32, jwtksServerRef, upstream *core.ResourceRef) *gloov1.Proxy {
	jwtCfg := &jwtplugin.VhostExtension{
		Providers: map[string]*jwtplugin.Provider{
			"provider1": {
				Jwks: &jwtplugin.Jwks{
					Jwks: &jwtplugin.Jwks_Remote{
						Remote: &jwtplugin.RemoteJwks{
							Url:         "http://test/keys",
							UpstreamRef: jwtksServerRef,
						},
					},
				},
				Issuer:    issuer,
				Audiences: []string{audience},
				KeepToken: true,
				TokenSource: &jwtplugin.TokenSource{
					Headers: []*jwtplugin.TokenSource_HeaderSource{{
						Header: "x-jwt",
						Prefix: "JWT ",
					}},
					QueryParams: []string{"jwttoken"},
				},
				ClaimsToHeaders: []*jwtplugin.ClaimToHeader{{
					Claim:  "sub",
					Header: "x-sub",
				}, {
					Claim:  "sub",
					Header: "x-sub",
					Append: true,
				},
					{
						Claim:  "data",
						Header: "x-data",
						Append: true,
					},
				},
			},
		},
	}

	return getProxyJwtRbacWithExtensions(envoyPort, jwtksServerRef, upstream, jwtCfg, nil)
}

func getProxyJwtRbacWithExtensions(envoyPort uint32, jwtksServerRef, upstream *core.ResourceRef, jwtCfg *jwtplugin.VhostExtension, rbacCfg *rbac.ExtensionSettings) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "virt1",
		Domains: []string{"*"},
		Options: &gloov1.VirtualHostOptions{
			Rbac: rbacCfg,
			JwtConfig: &gloov1.VirtualHostOptions_Jwt{
				Jwt: jwtCfg,
			},
		},
		Routes: []*gloov1.Route{
			{
				Options: &gloov1.RouteOptions{
					JwtConfig: &gloov1.RouteOptions_Jwt{
						Jwt: getDisabledJwt(),
					},
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/public_route",
					},
				}},
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: upstream,
								},
							},
						},
					},
				},
			}, {
				Options: &gloov1.RouteOptions{
					// Disable JWT and not RBAC, so that no one can get here
					JwtConfig: &gloov1.RouteOptions_Jwt{
						Jwt: getDisabledJwt(),
					},
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/private_route",
					},
				}},
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: upstream,
								},
							},
						},
					},
				},
			}, {
				Options: &gloov1.RouteOptions{
					Transformations: &glootransformation.Transformations{
						RequestTransformation: &glootransformation.Transformation{
							TransformationType: &glootransformation.Transformation_TransformationTemplate{
								TransformationTemplate: &glootransformation.TransformationTemplate{
									Headers:            map[string]*glootransformation.InjaTemplate{"x-new-header": {Text: "new"}},
									BodyTransformation: &glootransformation.TransformationTemplate_Passthrough{Passthrough: &glootransformation.Passthrough{}},
								},
							},
						},
					},
					// Disable RBAC and not JWT, for authn only tests
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					Headers: []*matchers.HeaderMatcher{{
						Name:  "x-sub",
						Value: "teatime,teatime",
					},
						{
							Name:  "x-data",
							Value: "{\"name\":\"test\"}",
						},
					},
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/authnonly",
					},
				}},
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: upstream,
								},
							},
						},
					},
				},
			}, {
				Options: &gloov1.RouteOptions{
					// Disable RBAC and not JWT, for authn only tests
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/authnonly",
					},
				}},
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: upstream,
								},
							},
						},
					},
				},
			}, {
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/",
					},
				}},
				Action: &gloov1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: upstream,
								},
							},
						},
					},
				},
			}},
	}

	vhosts = append(vhosts, vhost)

	p := &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: net.IPv4zero.String(),
			BindPort:    envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}

func getDisabledJwt() *jwtplugin.RouteExtension {
	return &jwtplugin.RouteExtension{
		Disable: true,
	}
}

func getDisabledRbac() *rbac.ExtensionSettings {
	return &rbac.ExtensionSettings{
		Disable: true,
	}
}

func waitForKeyService(jwksPort uint32) {
	Eventually(func() error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", jwksPort))
		if err != nil {
			return err
		}
		if resp == nil {
			return errors.New("Expected non-nil response from key service")
		}
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return errors.New(fmt.Sprintf("Unexpected status code from key server: %d", resp.StatusCode))
		}
		return nil
	}, "5s", "0.5s").ShouldNot(HaveOccurred())
}
