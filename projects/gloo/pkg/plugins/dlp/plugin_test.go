package dlp

import (
	"context"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation_ee"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	envoy_type "github.com/solo-io/solo-kit/pkg/api/external/envoy/type"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("dlp plugin", func() {
	var (
		plugin       *Plugin
		params       plugins.Params
		vhostParams  plugins.VirtualHostParams
		virtualHost  *v1.VirtualHost
		route        *v1.Route
		httpListener *v1.HttpListener
		dlpVhost     *dlp.Config
		dlpRoute     *dlp.Config
		dlpListener  *dlp.FilterConfig

		matchAll = &matchers.Matcher{
			PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
		}

		customTestAction = &dlp.Action{
			ActionType: dlp.Action_CUSTOM,
			Shadow:     true,
			CustomAction: &dlp.CustomAction{
				Name:  "test",
				Regex: []string{"regex"},
				Percent: &envoy_type.Percent{
					Value: 75,
				},
				MaskChar: "Z",
				RegexActions: []*transformation_ee.RegexAction{
					{Regex: "actionRegex", Subgroup: 1},
				},
			},
		}
	)

	BeforeEach(func() {
		dlpListener = &dlp.FilterConfig{}
	})

	JustBeforeEach(func() {
		if dlpRoute == nil {
			dlpRoute = &dlp.Config{}
		}
		route = &v1.Route{
			Matchers: []*matchers.Matcher{matchAll},
			Action: &v1.Route_DirectResponseAction{
				DirectResponseAction: &v1.DirectResponseAction{
					Status: 200,
					Body:   "test",
				},
			},
			Options: &v1.RouteOptions{
				Dlp: dlpRoute,
			},
		}

		if dlpVhost == nil {
			dlpVhost = &dlp.Config{}
		}

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Options: &v1.VirtualHostOptions{
				Dlp: dlpVhost,
			},
			Routes: []*v1.Route{route},
		}

		httpListener = &v1.HttpListener{
			VirtualHosts: []*v1.VirtualHost{virtualHost},
			Options: &v1.HttpListenerOptions{
				Dlp: dlpListener,
			},
		}
		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Listeners: []*v1.Listener{{
				Name: "default",
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: httpListener,
				},
			}},
		}

		params.Snapshot = &v1.ApiSnapshot{
			Proxies: v1.ProxyList{proxy},
		}
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}

	})

	var checkAllActions = func(actions []*dlp.Action, dlpTransform *transformation_ee.DlpTransformation, transformNum int) {
		Expect(dlpTransform).NotTo(BeNil())
		Expect(dlpTransform.GetActions()).To(HaveLen(transformNum))
		relevantActions := getRelevantActions(context.Background(), actions)
		Expect(dlpTransform.GetActions()).To(Equal(relevantActions))
	}

	var checkAllDefaultActions = func(actions []*dlp.Action, dlpTransform *transformation_ee.DlpTransformation) {
		checkAllActions(actions, dlpTransform, len(transformMap)-1)
	}
	var checkAllCCActions = func(actions []*dlp.Action, dlpTransform *transformation_ee.DlpTransformation) {
		checkAllActions(actions, dlpTransform, len(transformMap)-2)
	}

	var checkCustomAction = func(dlpTransform *transformation_ee.DlpTransformation) {
		Expect(dlpTransform).NotTo(BeNil())
		Expect(dlpTransform.GetActions()).To(HaveLen(1))
		customTransform := dlpTransform.GetActions()[0]
		Expect(customTransform.Shadow).To(Equal(customTestAction.Shadow))
		Expect(customTransform.MaskChar).To(Equal(customTestAction.CustomAction.MaskChar))
		Expect(customTransform.Name).To(Equal(customTestAction.CustomAction.Name))
		Expect(customTransform.Percent.Equal(customTestAction.CustomAction.Percent)).To(BeTrue())
		Expect(customTransform.Regex).To(Equal(customTestAction.CustomAction.Regex))
		Expect(len(customTransform.RegexActions)).To(Equal(1))
		Expect(customTransform.RegexActions[0].GetRegex()).To(Equal(customTestAction.CustomAction.RegexActions[0].GetRegex()))
		Expect(customTransform.RegexActions[0].GetSubgroup()).To(Equal(customTestAction.CustomAction.RegexActions[0].GetSubgroup()))
	}

	Context("process snapshot", func() {
		var (
			outRoute   envoy_config_route_v3.Route
			outVhost   envoy_config_route_v3.VirtualHost
			outFilters []plugins.StagedHttpFilter
		)

		var translateRoute = func() *transformation_ee.RouteTransformations {
			goTpfc := outRoute.TypedPerFilterConfig[FilterName]
			Expect(goTpfc).NotTo(BeNil())
			var perRouteDlp transformation_ee.RouteTransformations
			err := ptypes.UnmarshalAny(goTpfc, &perRouteDlp)
			Expect(err).NotTo(HaveOccurred())
			return &perRouteDlp
		}

		var translateVhost = func() *transformation_ee.RouteTransformations {
			goTpfc := outVhost.TypedPerFilterConfig[FilterName]
			Expect(goTpfc).NotTo(BeNil())
			var perVhostDlp transformation_ee.RouteTransformations
			err := ptypes.UnmarshalAny(goTpfc, &perVhostDlp)
			Expect(err).NotTo(HaveOccurred())
			return &perVhostDlp
		}

		JustBeforeEach(func() {
			outVhost = envoy_config_route_v3.VirtualHost{
				Name: "test",
			}
			outRoute = envoy_config_route_v3.Route{}
			routesParams := plugins.RouteParams{
				VirtualHostParams: vhostParams,
				VirtualHost:       virtualHost,
			}
			// run it like the translator:
			err := plugin.ProcessRoute(routesParams, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, httpListener)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			plugin = NewPlugin()
			plugin.Init(plugins.InitParams{})
		})

		Context("empty extensions", func() {
			It("can create the proper filters", func() {
				Expect(outFilters).To(HaveLen(1))
				wafFilter := outFilters[0]
				Expect(wafFilter.HttpFilter.Name).To(Equal(FilterName))
				Expect(wafFilter.Stage).To(Equal(plugins.BeforeStage(plugins.WafStage)))
				st := wafFilter.HttpFilter.GetTypedConfig()
				Expect(st).To(BeNil())
			})
		})

		Context("http filters", func() {

			var (
				dlpRule *dlp.DlpRule
			)

			var checkListenerFilter = func() *transformation_ee.TransformationRule {
				Expect(outFilters).To(HaveLen(1))
				dlpFilter := outFilters[0]
				Expect(dlpFilter.HttpFilter.Name).To(Equal(FilterName))
				Expect(dlpFilter.Stage).To(Equal(plugins.BeforeStage(plugins.WafStage)))
				goTypedConfig := dlpFilter.HttpFilter.GetTypedConfig()
				if goTypedConfig == nil {
					return nil
				}
				Expect(goTypedConfig).NotTo(BeNil())
				var filterDlp transformation_ee.FilterTransformations
				err := ptypes.UnmarshalAny(goTypedConfig, &filterDlp)
				Expect(err).NotTo(HaveOccurred())
				if len(filterDlp.GetTransformations()) == 0 {
					return nil
				}
				return filterDlp.GetTransformations()[0]
			}
			Context("nil", func() {
				BeforeEach(func() {
					dlpListener = &dlp.FilterConfig{
						DlpRules: []*dlp.DlpRule{
							{
								Matcher: nil,
								Actions: nil,
							},
						},
					}
				})

				It("can create the proper nil http filters", func() {
					filterDlp := checkListenerFilter()
					Expect(filterDlp.GetRouteTransformations()).To(BeNil())
				})
			})

			Context("default filters", func() {
				BeforeEach(func() {
					dlpRule = &dlp.DlpRule{
						Matcher: nil,
						Actions: nil,
					}
					dlpRule.Actions = make([]*dlp.Action, 0, len(transformMap))
					for key := range transformMap {
						dlpRule.Actions = append(dlpRule.Actions, &dlp.Action{
							ActionType: key,
						})
					}
					dlpListener = &dlp.FilterConfig{
						DlpRules: []*dlp.DlpRule{dlpRule},
					}
				})

				It("can create the proper filled http filters", func() {
					rule := checkListenerFilter()
					filterDlp := rule.GetRouteTransformations()
					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation()).NotTo(BeNil())
					checkAllDefaultActions(dlpRule.GetActions(), filterDlp.GetResponseTransformation().GetDlpTransformation())
					mAll := translator.GlooMatcherToEnvoyMatcher(params.Ctx, matchAll)
					expected := toGlooRouteMatch(&mAll)
					Expect(rule.GetMatch()).To(BeNil())
					Expect(*(rule.GetMatchV3())).To(Equal(*expected))
				})
			})

			Context("default filters with regex matcher", func() {
				matchRegex := &matchers.Matcher{
					PathSpecifier: &matchers.Matcher_Regex{
						Regex: "foo",
					},
				}

				BeforeEach(func() {
					dlpRule = &dlp.DlpRule{
						Matcher: matchRegex,
						Actions: nil,
					}
					dlpRule.Actions = make([]*dlp.Action, 0, len(transformMap))
					for key := range transformMap {
						dlpRule.Actions = append(dlpRule.Actions, &dlp.Action{
							ActionType: key,
						})
					}
					dlpListener = &dlp.FilterConfig{
						DlpRules: []*dlp.DlpRule{dlpRule},
					}
				})

				It("can create the proper filled http filters", func() {
					rule := checkListenerFilter()
					filterDlp := rule.GetRouteTransformations()
					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation()).NotTo(BeNil())
					checkAllDefaultActions(dlpRule.GetActions(), filterDlp.GetResponseTransformation().GetDlpTransformation())
					mRegex := translator.GlooMatcherToEnvoyMatcher(params.Ctx, matchRegex)
					expected := toGlooRouteMatch(&mRegex)
					Expect(rule.GetMatch()).To(BeNil())
					Expect(*(rule.GetMatchV3())).To(Equal(*expected))
				})
			})

			Context("all filters action with shadow", func() {
				BeforeEach(func() {
					dlpRule = &dlp.DlpRule{
						Matcher: nil,
						Actions: []*dlp.Action{
							{
								ActionType: dlp.Action_ALL_CREDIT_CARDS,
								Shadow:     true,
							},
						},
					}
					dlpListener = &dlp.FilterConfig{
						DlpRules: []*dlp.DlpRule{dlpRule},
					}
				})

				It("can create the proper filled http filters", func() {
					filterDlp := checkListenerFilter().GetRouteTransformations()
					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation()).NotTo(BeNil())
					checkAllCCActions(dlpRule.GetActions(), filterDlp.GetResponseTransformation().GetDlpTransformation())
				})
			})

			Context("custom filter", func() {
				BeforeEach(func() {
					dlpRule = &dlp.DlpRule{
						Matcher: nil,
						Actions: []*dlp.Action{customTestAction},
					}
					dlpListener = &dlp.FilterConfig{
						DlpRules: []*dlp.DlpRule{dlpRule},
					}
				})

				It("can create the proper filled http filters", func() {
					filterDlp := checkListenerFilter().GetRouteTransformations()
					checkCustomAction(filterDlp.GetResponseTransformation().GetDlpTransformation())
				})
			})

			Context("only access logs", func() {
				BeforeEach(func() {
					dlpRule = &dlp.DlpRule{
						Matcher: nil,
						Actions: nil,
					}
					dlpRule.Actions = make([]*dlp.Action, 0, len(transformMap))
					for key := range transformMap {
						dlpRule.Actions = append(dlpRule.Actions, &dlp.Action{
							ActionType: key,
						})
					}
					dlpListener = &dlp.FilterConfig{
						DlpRules:   []*dlp.DlpRule{dlpRule},
						EnabledFor: dlp.FilterConfig_ACCESS_LOGS,
					}
				})

				It("can create the proper filled http filters", func() {
					rule := checkListenerFilter()
					filterDlp := rule.GetRouteTransformations()
					Expect(filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation()).NotTo(BeNil())
					Expect(filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeTrue())
					Expect(filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeTrue())
					checkAllDefaultActions(dlpRule.GetActions(), filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation())

					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation()).To(BeNil())
					mAll := translator.GlooMatcherToEnvoyMatcher(params.Ctx, matchAll)
					expected := toGlooRouteMatch(&mAll)
					Expect(rule.GetMatch()).To(BeNil())
					Expect(*(rule.GetMatchV3())).To(Equal(*expected))
				})
			})

			Context("access logs and responses", func() {
				BeforeEach(func() {
					dlpRule = &dlp.DlpRule{
						Matcher: nil,
						Actions: nil,
					}
					dlpRule.Actions = make([]*dlp.Action, 0, len(transformMap))
					for key := range transformMap {
						dlpRule.Actions = append(dlpRule.Actions, &dlp.Action{
							ActionType: key,
						})
					}
					dlpListener = &dlp.FilterConfig{
						DlpRules:   []*dlp.DlpRule{dlpRule},
						EnabledFor: dlp.FilterConfig_ALL,
					}
				})

				It("can create the proper filled http filters", func() {
					rule := checkListenerFilter()
					filterDlp := rule.GetRouteTransformations()
					Expect(filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation()).NotTo(BeNil())
					Expect(filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeTrue())
					Expect(filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeTrue())
					checkAllDefaultActions(dlpRule.GetActions(), filterDlp.GetOnStreamCompletionTransformation().GetDlpTransformation())

					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation()).NotTo(BeNil())
					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeFalse())
					Expect(filterDlp.GetResponseTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeFalse())
					checkAllDefaultActions(dlpRule.GetActions(), filterDlp.GetResponseTransformation().GetDlpTransformation())
					mAll := translator.GlooMatcherToEnvoyMatcher(params.Ctx, matchAll)
					expected := toGlooRouteMatch(&mAll)
					Expect(rule.GetMatch()).To(BeNil())
					Expect(*(rule.GetMatchV3())).To(Equal(*expected))
				})
			})

		})

		Context("per route/vhost", func() {

			Context("nil", func() {
				BeforeEach(func() {
					dlpRoute = &dlp.Config{
						Actions: nil,
					}

					dlpVhost = &dlp.Config{
						Actions: nil,
					}
				})

				It("sets disabled on route", func() {
					pfc := outRoute.TypedPerFilterConfig[FilterName]
					Expect(pfc).To(BeNil())

				})

				It("sets disabled on vhost", func() {
					pfc := outVhost.TypedPerFilterConfig[FilterName]
					Expect(pfc).To(BeNil())
				})
			})

			Context("default actions", func() {
				BeforeEach(func() {
					dlpRoute = &dlp.Config{}
					dlpVhost = &dlp.Config{}
					for key := range transformMap {
						dlpRoute.Actions = append(dlpRoute.Actions, &dlp.Action{
							ActionType: key,
						})
						dlpVhost.Actions = append(dlpVhost.Actions, &dlp.Action{
							ActionType: key,
						})
					}
				})

				It("sets default actions on route", func() {
					perRouteDlp := translateRoute()
					checkAllDefaultActions(dlpRoute.GetActions(), perRouteDlp.GetResponseTransformation().GetDlpTransformation())
				})

				It("sets default actions on vhost", func() {
					perVhostDlp := translateVhost()
					checkAllDefaultActions(dlpVhost.GetActions(), perVhostDlp.GetResponseTransformation().GetDlpTransformation())
				})
			})

			Context("All default actions with shadow", func() {
				BeforeEach(func() {
					dlpRoute = &dlp.Config{
						Actions: []*dlp.Action{
							{
								ActionType: dlp.Action_ALL_CREDIT_CARDS,
								Shadow:     true,
							},
						},
					}
					dlpVhost = &dlp.Config{
						Actions: []*dlp.Action{
							{
								ActionType: dlp.Action_ALL_CREDIT_CARDS,
								Shadow:     true,
							},
						},
					}
				})

				It("sets default actions on route", func() {
					perRouteDlp := translateRoute()
					checkAllCCActions(dlpRoute.GetActions(), perRouteDlp.GetResponseTransformation().GetDlpTransformation())
				})

				It("sets default actions on vhost", func() {
					perVhostDlp := translateVhost()
					checkAllCCActions(dlpVhost.GetActions(), perVhostDlp.GetResponseTransformation().GetDlpTransformation())
				})
			})

			Context("custom action", func() {
				BeforeEach(func() {
					dlpRoute = &dlp.Config{
						Actions: []*dlp.Action{customTestAction},
					}
					dlpVhost = &dlp.Config{
						Actions: []*dlp.Action{customTestAction},
					}
				})

				It("sets default actions on route", func() {
					perRouteDlp := translateRoute()
					checkCustomAction(perRouteDlp.GetResponseTransformation().GetDlpTransformation())
				})

				It("sets default actions on vhost", func() {
					perVhostDlp := translateVhost()
					checkCustomAction(perVhostDlp.GetResponseTransformation().GetDlpTransformation())
				})
			})
		})

		Context("only access logs", func() {
			BeforeEach(func() {
				dlpRoute = &dlp.Config{}
				dlpVhost = &dlp.Config{}

				dlpRoute.EnabledFor = dlp.Config_ACCESS_LOGS
				dlpVhost.EnabledFor = dlp.Config_ACCESS_LOGS

				for key := range transformMap {
					dlpRoute.Actions = append(dlpRoute.Actions, &dlp.Action{
						ActionType: key,
					})
					dlpVhost.Actions = append(dlpVhost.Actions, &dlp.Action{
						ActionType: key,
					})
				}
			})

			It("sets default actions on route", func() {
				perRouteDlp := translateRoute()
				Expect(perRouteDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeTrue())
				Expect(perRouteDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeTrue())
				checkAllDefaultActions(dlpRoute.GetActions(), perRouteDlp.GetOnStreamCompletionTransformation().GetDlpTransformation())

				Expect(perRouteDlp.GetResponseTransformation()).To(BeNil())
			})

			It("sets default actions on vhost", func() {
				perVhostDlp := translateVhost()
				Expect(perVhostDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeTrue())
				Expect(perVhostDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeTrue())
				checkAllDefaultActions(dlpRoute.GetActions(), perVhostDlp.GetOnStreamCompletionTransformation().GetDlpTransformation())

				Expect(perVhostDlp.GetResponseTransformation()).To(BeNil())
			})
		})

		Context("access logs and responses", func() {
			BeforeEach(func() {
				dlpRoute = &dlp.Config{}
				dlpVhost = &dlp.Config{}

				dlpRoute.EnabledFor = dlp.Config_ALL
				dlpVhost.EnabledFor = dlp.Config_ALL

				for key := range transformMap {
					dlpRoute.Actions = append(dlpRoute.Actions, &dlp.Action{
						ActionType: key,
					})
					dlpVhost.Actions = append(dlpVhost.Actions, &dlp.Action{
						ActionType: key,
					})
				}
			})

			It("sets default actions on route", func() {
				perRouteDlp := translateRoute()
				Expect(perRouteDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeTrue())
				Expect(perRouteDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeTrue())
				checkAllDefaultActions(dlpRoute.GetActions(), perRouteDlp.GetOnStreamCompletionTransformation().GetDlpTransformation())

				Expect(perRouteDlp.GetResponseTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeFalse())
				Expect(perRouteDlp.GetResponseTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeFalse())
				checkAllDefaultActions(dlpRoute.GetActions(), perRouteDlp.GetResponseTransformation().GetDlpTransformation())
			})

			It("sets default actions on vhost", func() {
				perVhostDlp := translateVhost()
				Expect(perVhostDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeTrue())
				Expect(perVhostDlp.GetOnStreamCompletionTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeTrue())
				checkAllDefaultActions(dlpRoute.GetActions(), perVhostDlp.GetOnStreamCompletionTransformation().GetDlpTransformation())

				Expect(perVhostDlp.GetResponseTransformation().GetDlpTransformation().GetEnableHeaderTransformation()).To(BeFalse())
				Expect(perVhostDlp.GetResponseTransformation().GetDlpTransformation().GetEnableDynamicMetadataTransformation()).To(BeFalse())
				checkAllDefaultActions(dlpRoute.GetActions(), perVhostDlp.GetResponseTransformation().GetDlpTransformation())
			})
		})
	})
})
