package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/solo-io/ext-auth-service/pkg/server"

	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"

	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/gloo/pkg/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauthpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"
	ratelimitservice "github.com/solo-io/solo-projects/test/services/ratelimit"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Rate Limit Local E2E", func() {

	var (
		ctx             context.Context
		cancel          context.CancelFunc
		testClients     services.TestClients
		redisSession    *gexec.Session
		isServerHealthy func() (bool, error)
		glooSettings    *gloov1.Settings
		cache           memory.InMemoryResourceCache
		rlAddr          string
	)
	const (
		redisaddr = "127.0.0.1"
		redisport = uint32(6379)
		rladdr    = "127.0.0.1"
		rlport    = uint32(18081)
	)
	BeforeEach(func() {
		glooSettings = &gloov1.Settings{}
	})

	runAllTests := func() {

		Context("With envoy", func() {

			var (
				envoyInstance *services.EnvoyInstance
				testUpstream  *v1helpers.TestUpstream
				envoyPort     = uint32(8080)
			)

			BeforeEach(func() {
				var err error
				envoyInstance, err = envoyFactory.NewEnvoyInstance()
				Expect(err).NotTo(HaveOccurred())

				envoyInstance.RatelimitAddr = rladdr
				envoyInstance.RatelimitPort = rlport
				rlAddr = envoyInstance.LocalAddr()

				err = envoyInstance.Run(testClients.GlooPort)
				Expect(err).NotTo(HaveOccurred())

				testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
				// drain channel as we dont care about it
				go func() {
					for range testUpstream.C {
					}
				}()
				var opts clients.WriteOpts
				up := testUpstream.Upstream
				_, err = testClients.UpstreamClient.Write(up, opts)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				if envoyInstance != nil {
					_ = envoyInstance.Clean()
				}
			})

			It("should rate limit envoy", func() {

				hosts := map[string]bool{"host1": true}
				proxy := getAuthEnabledProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())
				EventuallyRateLimited("host1", envoyPort)
			})

			It("should rate limit two vhosts", func() {

				hosts := map[string]bool{"host1": true, "host2": true}
				proxy := getAuthEnabledProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())

				EventuallyRateLimited("host1", envoyPort)
				EventuallyRateLimited("host2", envoyPort)
			})

			It("should rate limit one of two vhosts", func() {

				hosts := map[string]bool{"host1": false, "host2": true}
				proxy := getAuthEnabledProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts)

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())

				ConsistentlyNotRateLimited("host1", envoyPort)
				EventuallyRateLimited("host2", envoyPort)
			})

			It("should ratelimit on route", func() {
				ingressRateLimit := &ratelimit.IngressRateLimit{
					AnonymousLimits: &rlv1alpha1.RateLimit{
						RequestsPerUnit: 1,
						Unit:            rlv1alpha1.RateLimit_SECOND,
					},
				}
				rlb := RlProxyBuilder{
					envoyPort:                    envoyPort,
					upstream:                     testUpstream.Upstream.Metadata.Ref(),
					hostsToVirtualHostRateLimits: map[string]bool{"host1": false},
					hostsToRouteRateLimits:       map[string]bool{"host1": true},
					ingressRateLimit:             ingressRateLimit,
				}
				proxy := rlb.getProxy()
				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				Eventually(isServerHealthy, "5s").Should(BeTrue())
				EventuallyRateLimited("host1/noauth", envoyPort)
			})

			Context("with auth", func() {

				BeforeEach(func() {
					// start the ext auth server
					extauthport := uint32(9100)

					extauthserver := &gloov1.Upstream{
						Metadata: core.Metadata{
							Name:      "extauth-server",
							Namespace: "default",
						},
						UseHttp2: &types.BoolValue{Value: true},
						UpstreamType: &gloov1.Upstream_Static{
							Static: &gloov1static.UpstreamSpec{
								Hosts: []*gloov1static.Host{{
									Addr: envoyInstance.LocalAddr(),
									Port: extauthport,
								}},
							},
						},
					}

					_, err := testClients.AuthConfigClient.Write(&extauthpb.AuthConfig{
						Metadata: core.Metadata{
							Name:      GetBasicAuthExtension().GetConfigRef().Name,
							Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauthpb.AuthConfig_Config{{
							AuthConfig: &extauthpb.AuthConfig_Config_BasicAuth{
								BasicAuth: getBasicAuthConfig(),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					_, err = testClients.UpstreamClient.Write(extauthserver, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					ref := extauthserver.Metadata.Ref()
					extauthSettings := &extauthpb.Settings{
						ExtauthzServerRef: &ref,
					}
					glooSettings.Extauth = extauthSettings

					settings := extauthrunner.Settings{
						GlooAddress: fmt.Sprintf("localhost:%d", testClients.GlooPort),
						ExtAuthSettings: server.Settings{
							DebugPort:    0,
							ServerPort:   int(extauthport),
							SigningKey:   "hello",
							UserIdHeader: "X-User-Id",
						},
					}
					go func(testctx context.Context) {
						defer GinkgoRecover()
						err := extauthrunner.RunWithSettings(testctx, settings)
						if testctx.Err() == nil {
							Expect(err).NotTo(HaveOccurred())
						}
					}(ctx)
				})

				It("should ratelimit authorized users", func() {
					ingressRateLimit := &ratelimit.IngressRateLimit{
						AuthorizedLimits: &rlv1alpha1.RateLimit{
							RequestsPerUnit: 1,
							Unit:            rlv1alpha1.RateLimit_SECOND,
						},
					}
					rlb := RlProxyBuilder{
						envoyPort:                    envoyPort,
						upstream:                     testUpstream.Upstream.Metadata.Ref(),
						hostsToVirtualHostRateLimits: map[string]bool{"host1": true},
						ingressRateLimit:             ingressRateLimit,
					}
					proxy := rlb.getProxy()
					vhost := proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.VirtualHosts[0]
					vhost.Options.Extauth = GetBasicAuthExtension()
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					// do the eventually first to give envoy a chance to start
					EventuallyRateLimited("user:password@host1", envoyPort)
					ConsistentlyNotRateLimited("host1/noauth", envoyPort)
				})
			})

			Context("reserved keyword rules (i.e., weighted and applyAlways rules)", func() {
				BeforeEach(func() {
					glooSettings.Ratelimit = &ratelimit.ServiceSettings{
						Descriptors: []*rlv1alpha1.Descriptor{
							{
								Key:   "generic_key",
								Value: "unprioritized",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
							},
							{
								Key:   "generic_key",
								Value: "prioritized",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_SECOND,
									RequestsPerUnit: 1000,
								},
								Weight: 1,
							},
							{
								Key:   "generic_key",
								Value: "always",
								RateLimit: &rlv1alpha1.RateLimit{
									Unit:            rlv1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 2,
								},
								AlwaysApply: true,
							},
						},
					}
				})

				It("should honor weighted rate limit rules", func() {
					hosts := map[string]bool{"host1": true}
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "unprioritized"},
							}},
						}}}

					proxy := getCustomProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts, rateLimits)
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					EventuallyRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// add a new rate limit action that points to a weighted rule with generous limit
					weightedAction := &rlv1alpha1.RateLimitActions{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "prioritized"},
							}},
						}}
					rateLimits = append(rateLimits, weightedAction)

					proxy = getCustomProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts, rateLimits)
					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// weighted rule has generous limit that will not be hit, however its larger weight trumps
					// the previous rule (that returned 429 before). we do not expect this to rate limit anymore
					ConsistentlyNotRateLimited("host1", envoyPort)
				})

				It("should honor alwaysApply rate limit rules", func() {
					hosts := map[string]bool{"host1": true}
					// add a prioritized rule to match against (has largest weight)
					rateLimits := []*rlv1alpha1.RateLimitActions{{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "prioritized"},
							}},
						}}}

					proxy := getCustomProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts, rateLimits)
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(isServerHealthy, "5s").Should(BeTrue())
					ConsistentlyNotRateLimited("host1", envoyPort)

					err = testClients.ProxyClient.Delete(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.DeleteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// add a new rate limit action that points to a "concurrent" rule, i.e. always evaluated
					weightedAction := &rlv1alpha1.RateLimitActions{
						Actions: []*rlv1alpha1.Action{{
							ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
								GenericKey: &rlv1alpha1.Action_GenericKey{DescriptorValue: "always"},
							}},
						}}
					rateLimits = append(rateLimits, weightedAction)

					proxy = getCustomProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts, rateLimits)
					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					// we added a ratelimit action that points to a rule with alwaysApply: true. Even though the rule
					// has zero weight, we will still evaluate the rule. the original request matched a weighted rule
					// that was too generous to return a 429, but the new rule should trigger and return a 429
					EventuallyRateLimited("host1", envoyPort)
				})
			})

		})
	}
	justBeforeEach := func() {
		// add the rl service as a static upstream
		rlserver := &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "rl-server",
				Namespace: "default",
			},
			UseHttp2: &types.BoolValue{Value: true},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: rlAddr,
						Port: rlport,
					}},
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(rlserver, clients.WriteOpts{})
		Expect(err).ToNot(HaveOccurred())

		ref := rlserver.Metadata.Ref()
		rlSettings := &ratelimit.Settings{
			RatelimitServerRef: &ref,
		}

		isServerHealthy = ratelimitservice.RunRateLimitServer(ctx, rladdr, testClients.GlooPort)

		glooSettings.RatelimitServer = rlSettings

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, glooSettings)
	}

	Context("Redis-backed rate limiting", func() {

		BeforeEach(func() {
			var err error
			os.Setenv("REDIS_URL", fmt.Sprintf("%s:%d", redisaddr, redisport))
			os.Setenv("REDIS_SOCKET_TYPE", "tcp")

			command := exec.Command(getRedisPath(), "--port", "6379")
			redisSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			// give redis a chance to start
			Eventually(redisSession.Out, "5s").Should(gbytes.Say("Ready to accept connections"))

			ctx, cancel = context.WithCancel(context.Background())
			cache = memory.NewInMemoryResourceCache()

			testClients = services.GetTestClients(cache)
			testClients.GlooPort = int(services.AllocateGlooPort())

		})
		JustBeforeEach(justBeforeEach)

		AfterEach(func() {
			cancel()
			redisSession.Kill()
		})

		runAllTests()
	})

	Context("DynamoDb-backed rate limiting", func() {

		BeforeEach(func() {
			// By setting these environment variables to non-empty values we signal we want to use DynamoDb
			// instead of Redis as our rate limiting backend. Local DynamoDB requires any non-empty creds to work
			os.Setenv("AWS_ACCESS_KEY_ID", "fakeMyKeyId")
			os.Setenv("AWS_SECRET_ACCESS_KEY", "fakeSecretAccessKey")

			awsEndpoint := "http://" + services.GetDynamoDbHost() + ":" + services.DynamoDbPort
			// Set AWS session to use local DynamoDB instead of defaulting to live AWS web services
			os.Setenv("AWS_ENDPOINT", awsEndpoint)

			services.RunDynamoDbContainer()
			Eventually(services.DynamoDbHealthCheck(awsEndpoint), "5s", "100ms").Should(BeEquivalentTo(services.HealthCheck{IsHealthy: true}))

			ctx, cancel = context.WithCancel(context.Background())
			cache = memory.NewInMemoryResourceCache()

			testClients = services.GetTestClients(cache)
			testClients.GlooPort = int(services.AllocateGlooPort())
		})

		JustBeforeEach(justBeforeEach)

		AfterEach(func() {
			cancel()
			services.MustKillAndRemoveContainer(services.DynamoDbContainerName)
		})

		runAllTests()
	})
})

func EventuallyOk(hostname string, port uint32) {
	// wait for three seconds so gloo race can be waited out
	// it's possible gloo upstreams hit after the proxy does
	// (gloo resyncs once per second)
	time.Sleep(3 * time.Second)
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func ConsistentlyNotRateLimited(hostname string, port uint32) {
	// waiting for envoy to start, so that consistently works
	EventuallyOk(hostname, port)

	ConsistentlyWithOffset(2, func() error {
		res, err := get(hostname, port)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func EventuallyRateLimited(hostname string, port uint32) {
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusTooManyRequests {
			return errors.New(fmt.Sprintf("%v is not TooManyRequests", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func get(hostname string, port uint32) (*http.Response, error) {
	parts := strings.SplitN(hostname, "/", 2)
	hostname = parts[0]
	path := "1"
	if len(parts) > 1 {
		path = parts[1]
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/"+path, "localhost", port), nil)
	Expect(err).NotTo(HaveOccurred())

	// remove password part if exists
	parts = strings.SplitN(hostname, "@", 2)
	if len(parts) > 1 {
		hostname = parts[1]
		auth := strings.Split(parts[0], ":")
		req.SetBasicAuth(auth[0], auth[1])
	}

	req.Host = hostname
	return http.DefaultClient.Do(req)
}

func getAuthEnabledProxy(envoyPort uint32, upstream core.ResourceRef, hostsToRateLimits map[string]bool) *gloov1.Proxy {
	ingressRateLimit := &ratelimit.IngressRateLimit{
		AnonymousLimits: &rlv1alpha1.RateLimit{
			RequestsPerUnit: 1,
			Unit:            rlv1alpha1.RateLimit_SECOND,
		},
	}
	rlb := RlProxyBuilder{
		envoyPort:                    envoyPort,
		upstream:                     upstream,
		hostsToVirtualHostRateLimits: hostsToRateLimits,
		ingressRateLimit:             ingressRateLimit,
	}
	return rlb.getProxy()
}

type RlProxyBuilder struct {
	ingressRateLimit             *ratelimit.IngressRateLimit
	upstream                     core.ResourceRef
	hostsToVirtualHostRateLimits map[string]bool
	hostsToRouteRateLimits       map[string]bool
	envoyPort                    uint32
}

func (b *RlProxyBuilder) getProxy() *gloov1.Proxy {

	var vhosts []*gloov1.VirtualHost

	for hostname, enableRateLimits := range b.hostsToVirtualHostRateLimits {
		routeRateLimit := b.ingressRateLimit
		if !b.hostsToRouteRateLimits[hostname] {
			routeRateLimit = nil
		}

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system_virt" + hostname,
			Domains: []string{hostname},
			Routes: []*gloov1.Route{
				{
					Name: "gloo-system_route-noauth-" + hostname,
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/noauth",
						},
					}},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(b.upstream),
									},
								},
							},
						},
					},
					Options: &gloov1.RouteOptions{
						Extauth:        &extauthpb.ExtAuthExtension{Spec: &extauthpb.ExtAuthExtension_Disable{Disable: true}},
						RatelimitBasic: routeRateLimit,
					},
				},
				{
					Name: "gloo-system_route-auth-" + hostname,
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(b.upstream),
									},
								},
							},
						},
					},
					Options: &gloov1.RouteOptions{
						RatelimitBasic: routeRateLimit,
					},
				},
			},
		}

		if enableRateLimits {
			vhost.Options = &gloov1.VirtualHostOptions{
				RatelimitBasic: b.ingressRateLimit,
			}
		}
		vhosts = append(vhosts, vhost)
	}

	p := &gloov1.Proxy{
		Metadata: core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "0.0.0.0",
			BindPort:    b.envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}

func getCustomProxy(envoyPort uint32, upstream core.ResourceRef, hostsToRateLimits map[string]bool, rateLimits []*rlv1alpha1.RateLimitActions) *gloov1.Proxy {
	rlVhostExt := &ratelimit.RateLimitVhostExtension{
		RateLimits: rateLimits,
	}
	rlb := CustomRlProxyBuilder{
		envoyPort:         envoyPort,
		upstream:          upstream,
		hostsToRateLimits: hostsToRateLimits,
		customRateLimit:   rlVhostExt,
	}
	return rlb.getCustomProxy()
}

type CustomRlProxyBuilder struct {
	customRateLimit   *ratelimit.RateLimitVhostExtension
	upstream          core.ResourceRef
	hostsToRateLimits map[string]bool
	envoyPort         uint32
}

func (b *CustomRlProxyBuilder) getCustomProxy() *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	for hostname, enableRateLimits := range b.hostsToRateLimits {
		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system_virt" + hostname,
			Domains: []string{hostname},
			Routes: []*gloov1.Route{
				{
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(b.upstream),
									},
								},
							},
						},
					},
				},
			},
		}

		if enableRateLimits {
			vhost.Options = &gloov1.VirtualHostOptions{
				RateLimitConfigType: &gloov1.VirtualHostOptions_Ratelimit{
					Ratelimit: b.customRateLimit,
				},
			}
		}
		vhosts = append(vhosts, vhost)
	}

	p := &gloov1.Proxy{
		Metadata: core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "0.0.0.0",
			BindPort:    b.envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}
func getRedisPath() string {
	binaryPath := os.Getenv("REDIS_BINARY")
	if binaryPath != "" {
		return binaryPath
	}
	return "redis-server"
}
