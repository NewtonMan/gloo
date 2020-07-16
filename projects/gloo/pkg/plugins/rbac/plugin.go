package rbac

import (
	"context"
	"sort"

	envoyroutev2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoycfgauthz "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauthz "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	"github.com/gogo/protobuf/proto"

	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
)

const (
	FilterName    = "envoy.filters.http.rbac"
	ExtensionName = "rbac"
)

var (
	_           plugins.Plugin = new(Plugin)
	filterStage                = plugins.DuringStage(plugins.AuthZStage)
)

type Plugin struct {
	settings *rbac.Settings
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.settings = params.Settings.GetRbac()
	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroutev2.VirtualHost) error {
	rbacConf := in.Options.GetRbac()
	if rbacConf == nil {
		// no config found, nothing to do here
		return nil
	}

	if rbacConf.Disable == true {
		perRouteRbac := &envoyauthz.RBACPerRoute{}
		pluginutils.SetVhostPerFilterConfig(out, FilterName, perRouteRbac)
		return nil
	}

	perRouteRbac, err := translateRbac(params.Ctx, in.Name, rbacConf.GetPolicies())
	if err != nil {
		return err
	}
	pluginutils.SetVhostPerFilterConfig(out, FilterName, perRouteRbac)

	return nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroutev2.Route) error {
	rbacConfig := in.GetOptions().GetRbac()
	if rbacConfig == nil {
		// no config found, nothing to do here
		return nil
	}

	var perRouteRbac *envoyauthz.RBACPerRoute

	if rbacConfig.Disable {
		perRouteRbac = &envoyauthz.RBACPerRoute{}
	} else {
		var err error
		perRouteRbac, err = translateRbac(params.Ctx, params.VirtualHost.Name, rbacConfig.GetPolicies())
		if err != nil {
			return err
		}
	}
	if perRouteRbac != nil {
		pluginutils.SetRoutePerFilterConfig(out, FilterName, perRouteRbac)
	}
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	strict := p.settings.GetRequireRbac()

	var cfg proto.Message
	if strict {
		// add a default config that denies everything
		cfg = &envoyauthz.RBAC{
			Rules: &envoycfgauthz.RBAC{
				Action: envoycfgauthz.RBAC_ALLOW,
			},
		}
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, cfg, filterStage)
	if err != nil {
		return nil, err
	}
	var filters []plugins.StagedHttpFilter
	filters = append(filters, stagedFilter)
	return filters, nil
}

func translateRbac(ctx context.Context, vhostname string, userPolicies map[string]*rbac.Policy) (*envoyauthz.RBACPerRoute, error) {
	ctx = contextutils.WithLogger(ctx, "rbac")
	policies := make(map[string]*envoycfgauthz.Policy)
	res := &envoyauthz.RBACPerRoute{
		Rbac: &envoyauthz.RBAC{
			Rules: &envoycfgauthz.RBAC{
				Action:   envoycfgauthz.RBAC_ALLOW,
				Policies: policies,
			},
		},
	}
	if userPolicies != nil {
		for k, v := range userPolicies {
			policies[k] = translatePolicy(contextutils.WithLogger(ctx, k), vhostname, v)
		}
	}
	return res, nil
}
func translatedMethods(methods []string) *envoycfgauthz.Permission {
	var allPermissions []*envoycfgauthz.Permission
	for _, method := range methods {
		allPermissions = append(allPermissions, &envoycfgauthz.Permission{
			Rule: &envoycfgauthz.Permission_Header{
				Header: &envoyroute.HeaderMatcher{
					Name: ":method",
					HeaderMatchSpecifier: &envoyroute.HeaderMatcher_ExactMatch{
						ExactMatch: method,
					},
				},
			},
		})
	}

	if len(allPermissions) == 1 {
		return allPermissions[0]
	}

	return &envoycfgauthz.Permission{
		Rule: &envoycfgauthz.Permission_OrRules{
			OrRules: &envoycfgauthz.Permission_Set{
				Rules: allPermissions,
			},
		},
	}
}

func translatePolicy(ctx context.Context, vhostname string, p *rbac.Policy) *envoycfgauthz.Policy {
	outPolicy := &envoycfgauthz.Policy{}
	for _, principal := range p.GetPrincipals() {
		outPrincipal := translateJwtPrincipal(ctx, vhostname, principal.JwtPrincipal)
		if outPrincipal != nil {
			outPolicy.Principals = append(outPolicy.Principals, outPrincipal)
		}
	}

	var allPermissions []*envoycfgauthz.Permission
	if permission := p.GetPermissions(); permission != nil {
		if permission.PathPrefix != "" {
			allPermissions = append(allPermissions, &envoycfgauthz.Permission{
				Rule: &envoycfgauthz.Permission_Header{
					Header: &envoyroute.HeaderMatcher{
						Name: ":path",
						HeaderMatchSpecifier: &envoyroute.HeaderMatcher_PrefixMatch{
							PrefixMatch: permission.PathPrefix,
						},
					},
				},
			})
		}

		if len(permission.Methods) != 0 {
			allPermissions = append(allPermissions, translatedMethods(permission.Methods))
		}
	}

	if len(allPermissions) == 0 {
		outPolicy.Permissions = []*envoycfgauthz.Permission{{
			Rule: &envoycfgauthz.Permission_Any{
				Any: true,
			},
		}}
	} else if len(allPermissions) == 1 {
		outPolicy.Permissions = []*envoycfgauthz.Permission{allPermissions[0]}
	} else {
		outPolicy.Permissions = []*envoycfgauthz.Permission{{
			Rule: &envoycfgauthz.Permission_AndRules{
				AndRules: &envoycfgauthz.Permission_Set{
					Rules: allPermissions,
				},
			},
		}}
	}

	return outPolicy
}

func getName(vhostname string, jwtPrincipal *rbac.JWTPrincipal) string {
	if vhostname == "" {
		return jwt.PayloadInMetadata
	}
	if jwtPrincipal.GetProvider() == "" {
		return jwt.PayloadInMetadata
	}
	return jwt.ProviderName(vhostname, jwtPrincipal.GetProvider())
}

func sortedKeys(m map[string]string) (keys []string) {
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

func translateJwtPrincipal(ctx context.Context, vhostname string, jwtPrincipal *rbac.JWTPrincipal) *envoycfgauthz.Principal {
	var jwtPrincipals []*envoycfgauthz.Principal
	claims := jwtPrincipal.GetClaims()
	// sort for idempotency
	for _, claim := range sortedKeys(claims) {
		value := claims[claim]
		claimPrincipal := &envoycfgauthz.Principal{
			Identifier: &envoycfgauthz.Principal_Metadata{
				Metadata: &envoymatcher.MetadataMatcher{
					Filter: "envoy.filters.http.jwt_authn",
					Path: []*envoymatcher.MetadataMatcher_PathSegment{
						{
							Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
								Key: getName(vhostname, jwtPrincipal),
							},
						},
						{
							Segment: &envoymatcher.MetadataMatcher_PathSegment_Key{
								Key: claim,
							},
						},
					},
					Value: &envoymatcher.ValueMatcher{
						MatchPattern: &envoymatcher.ValueMatcher_StringMatch{
							StringMatch: &envoymatcher.StringMatcher{
								MatchPattern: &envoymatcher.StringMatcher_Exact{
									Exact: value,
								},
							},
						},
					},
				},
			},
		}
		jwtPrincipals = append(jwtPrincipals, claimPrincipal)
	}

	if len(jwtPrincipals) == 0 {
		logger := contextutils.LoggerFrom(ctx)
		logger.Info("RBAC JWT Principal with zero claims - ignoring")
		return nil
	} else if len(jwtPrincipals) == 1 {
		return jwtPrincipals[0]
	}
	return &envoycfgauthz.Principal{
		Identifier: &envoycfgauthz.Principal_AndIds{
			AndIds: &envoycfgauthz.Principal_Set{
				Ids: jwtPrincipals,
			},
		},
	}
}
