package extauth_test

import (
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/extauth"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	"github.com/solo-io/gloo/pkg/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Plugin", func() {

	var (
		params        plugins.Params
		vhostParams   plugins.VirtualHostParams
		routeParams   plugins.RouteParams
		plugin        *Plugin
		virtualHost   *v1.VirtualHost
		upstream      *v1.Upstream
		secret        *v1.Secret
		route         *v1.Route
		authConfig    *extauthv1.AuthConfig
		authExtension *extauthv1.ExtAuthExtension
		clientSecret  *extauthv1.OauthSecret
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		err := plugin.Init(plugins.InitParams{})
		Expect(err).ToNot(HaveOccurred())

		upstream = &v1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth",
				Namespace: "default",
			},
			UpstreamType: &v1.Upstream_Static{
				Static: &static.UpstreamSpec{
					Hosts: []*static.Host{{
						Addr: "test",
						Port: 1234,
					}},
				},
			},
		}
		route = &v1.Route{
			Matchers: []*matchers.Matcher{{
				PathSpecifier: &matchers.Matcher_Prefix{
					Prefix: "/",
				},
			}},
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: utils.ResourceRefPtr(upstream.Metadata.Ref()),
							},
						},
					},
				},
			},
		}

		clientSecret = &extauthv1.OauthSecret{
			ClientSecret: "1234",
		}

		secret = &v1.Secret{
			Metadata: core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &v1.Secret_Oauth{
				Oauth: clientSecret,
			},
		}
		secretRef := secret.Metadata.Ref()

		authConfig = &extauthv1.AuthConfig{
			Metadata: core.Metadata{
				Name:      "oauth",
				Namespace: "gloo-system",
			},
			Configs: []*extauthv1.AuthConfig_Config{{
				AuthConfig: &extauthv1.AuthConfig_Config_Oauth{
					Oauth: &extauthv1.OAuth{
						ClientSecretRef: &secretRef,
						ClientId:        "ClientId",
						IssuerUrl:       "IssuerUrl",
						AppUrl:          "AppUrl",
						CallbackPath:    "CallbackPath",
					},
				},
			}},
		}
		authConfigRef := authConfig.Metadata.Ref()
		authExtension = &extauthv1.ExtAuthExtension{
			Spec: &extauthv1.ExtAuthExtension_ConfigRef{
				ConfigRef: &authConfigRef,
			},
		}
	})

	JustBeforeEach(func() {

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Options: &v1.VirtualHostOptions{
				Extauth: authExtension,
			},
			Routes: []*v1.Route{route},
		}

		proxy := &v1.Proxy{
			Metadata: core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Listeners: []*v1.Listener{{
				Name: "default",
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: &v1.HttpListener{
						VirtualHosts: []*v1.VirtualHost{virtualHost},
					},
				},
			}},
		}

		params.Snapshot = &v1.ApiSnapshot{
			Proxies:     v1.ProxyList{proxy},
			Upstreams:   v1.UpstreamList{upstream},
			Secrets:     v1.SecretList{secret},
			AuthConfigs: extauthv1.AuthConfigList{authConfig},
		}
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}
		routeParams = plugins.RouteParams{
			VirtualHostParams: vhostParams,
			VirtualHost:       virtualHost,
		}
	})

	Context("no extauth settings", func() {
		It("should provide sanitize filter", func() {
			filters, err := plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))
		})
	})

	Context("with extauth server", func() {
		var (
			extAuthRef      *core.ResourceRef
			extAuthSettings *extauthv1.Settings
		)
		BeforeEach(func() {
			second := time.Second
			extAuthRef = &core.ResourceRef{
				Name:      "extauth",
				Namespace: "default",
			}
			extAuthSettings = &extauthv1.Settings{
				ExtauthzServerRef: extAuthRef,
				FailureModeAllow:  true,
				RequestBody: &extauthv1.BufferSettings{
					AllowPartialMessage: true,
					MaxRequestBytes:     54,
				},
				RequestTimeout: &second,
			}
		})
		JustBeforeEach(func() {
			err := plugin.Init(plugins.InitParams{
				Settings: &v1.Settings{Extauth: extAuthSettings},
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("should provide sanitize filter with listener overriding global", func() {
			filters, err := plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))

			goTpfc := filters[0].HttpFilter.GetTypedConfig()
			Expect(goTpfc).NotTo(BeNil())
			var sanitizeCfg extauth.Sanitize
			gogoTpfc := &types.Any{TypeUrl: goTpfc.TypeUrl, Value: goTpfc.Value}
			err = types.UnmarshalAny(gogoTpfc, &sanitizeCfg)
			Expect(err).ToNot(HaveOccurred())

			Expect(sanitizeCfg.HeadersToRemove).To(Equal([]string{DefaultAuthHeader}))

			// now provide a listener override for auth header
			extAuthSettings.UserIdHeader = "override"
			listener := &v1.HttpListener{
				VirtualHosts: []*v1.VirtualHost{virtualHost},
				Options:      &v1.HttpListenerOptions{Extauth: extAuthSettings},
			}

			filters, err = plugin.HttpFilters(params, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))
			Expect(filters[0].HttpFilter.Name).To(Equal(SanitizeFilterName))

			goTpfc = filters[0].HttpFilter.GetTypedConfig()
			Expect(goTpfc).NotTo(BeNil())
			gogoTpfc = &types.Any{TypeUrl: goTpfc.TypeUrl, Value: goTpfc.Value}
			err = types.UnmarshalAny(gogoTpfc, &sanitizeCfg)
			Expect(err).ToNot(HaveOccurred())

			Expect(sanitizeCfg.HeadersToRemove).To(Equal([]string{"override"}))
		})

		It("should not error processing vhost", func() {
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})

		It("should mark vhost with no auth as disabled", func() {
			// remove auth extension
			virtualHost.Options.Extauth = nil
			var out envoyroute.VirtualHost
			err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should mark route with extension as disabled", func() {
			disabled := &extauthv1.ExtAuthExtension{
				Spec: &extauthv1.ExtAuthExtension_Disable{Disable: true},
			}

			route.Options = &v1.RouteOptions{
				Extauth: disabled,
			}
			var out envoyroute.Route
			err := plugin.ProcessRoute(routeParams, route, &out)
			Expect(err).NotTo(HaveOccurred())
			ExpectDisabled(&out)
		})

		It("should do nothing to a route that's not explicitly disabled", func() {
			var out envoyroute.Route
			err := plugin.ProcessRoute(routeParams, route, &out)
			Expect(err).NotTo(HaveOccurred())
			Expect(IsDisabled(&out)).To(BeFalse())
		})
	})

})

type envoyTypedPerFilterConfig interface {
	GetTypedPerFilterConfig() map[string]*any.Any
}

func ExpectDisabled(e envoyTypedPerFilterConfig) {
	Expect(IsDisabled(e)).To(BeTrue())
}

// Returns true if the ext_authz filter is explicitly disabled
func IsDisabled(e envoyTypedPerFilterConfig) bool {
	if e.GetTypedPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := ptypes.UnmarshalAny(e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization], &cfg)
	Expect(err).NotTo(HaveOccurred())

	return cfg.GetDisabled()
}

// Returns true if the ext_authz filter is enabled and if the ContextExtensions have the expected number of entries.
func IsEnabled(e envoyTypedPerFilterConfig) bool {
	if e.GetTypedPerFilterConfig() == nil {
		return false
	}
	if _, ok := e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization]; !ok {
		return false
	}
	var cfg envoyauth.ExtAuthzPerRoute
	err := ptypes.UnmarshalAny(e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization], &cfg)
	Expect(err).NotTo(HaveOccurred())

	if cfg.GetCheckSettings() == nil {
		return false
	}

	return len(cfg.GetCheckSettings().ContextExtensions) == 3
}

// Returns true if no PerFilterConfig is set for the ext_authz filter
func IsNotSet(e envoyTypedPerFilterConfig) bool {
	if e.GetTypedPerFilterConfig() == nil {
		return true
	}
	_, ok := e.GetTypedPerFilterConfig()[wellknown.HTTPExternalAuthorization]
	return !ok
}
