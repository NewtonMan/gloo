package advanced_http

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes/wrappers"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	envoy_core_v3_endpoint "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	envoy_advanced_http "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/advanced_http"
	envoy_type_matcher_v3_solo "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	envoy_type_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloo_advanced_http "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/advanced_http"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
)

var (
	_ plugins.Plugin         = new(plugin)
	_ plugins.UpstreamPlugin = new(plugin)
)

const (
	ExtensionName     = "advanced_http"
	HealthCheckerName = "io.solo.health_checkers.advanced_http"
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Init(_ plugins.InitParams) {
}

func (p *plugin) Name() string {
	return ExtensionName
}

func shouldProcess(in *gloov1.Upstream) bool {
	if len(in.GetHealthChecks()) == 0 {
		return false
	}

	// only do this for static upstreams with custom health path and/or method defined,
	// so that we only use new logic when we have to. this is done to minimize potential error impact.
	for _, host := range in.GetStatic().GetHosts() {
		if host.GetHealthCheckConfig().GetPath() != "" {
			return true
		}
		if host.GetHealthCheckConfig().GetMethod() != "" {
			return true
		}
	}
	// do this for failover endpoints that have custom health path and/or method defined
	for _, priority := range in.GetFailover().GetPrioritizedLocalities() {
		for _, localityEndpoints := range priority.GetLocalityEndpoints() {
			for _, endpoint := range localityEndpoints.GetLbEndpoints() {
				if endpoint.GetHealthCheckConfig().GetPath() != "" {
					return true
				}
				if endpoint.GetHealthCheckConfig().GetMethod() != "" {
					return true
				}
			}
		}
	}
	for _, hc := range in.GetHealthChecks() {
		if hc.GetHttpHealthCheck().GetResponseAssertions() != nil {
			return true
		}
	}
	return false
}

func (p *plugin) ProcessUpstream(params plugins.Params, in *gloov1.Upstream, out *envoy_config_cluster_v3.Cluster) error {

	// only do this for static upstreams with custom health path defined.
	// so that we only use new logic when we have to. this is done to minimize potential error impact.
	if !shouldProcess(in) {
		return nil
	}

	secretList := params.Snapshot.Secrets
	out.HealthChecks = make([]*envoy_config_core_v3.HealthCheck, len(in.GetHealthChecks()))
	shouldEnforceNamespaceMatch := false
	shouldEnforceStr := os.Getenv(api_conversion.MatchingNamespaceEnv)
	if shouldEnforceStr != "" {
		var err error
		shouldEnforceNamespaceMatch, err = strconv.ParseBool(shouldEnforceStr)
		if err != nil {
			return err
		}
	}
	for i, hcCfg := range in.GetHealthChecks() {
		envoyHc, err := api_conversion.ToEnvoyHealthCheck(hcCfg, &secretList, api_conversion.HeaderSecretOptions{EnforceNamespaceMatch: shouldEnforceNamespaceMatch, UpstreamNamespace: in.GetMetadata().GetNamespace()})
		if err != nil {
			return err
		}

		glooHc, err := convertEnvoyToGloo(params.Ctx, envoyHc.GetHttpHealthCheck())
		if err != nil {
			return err
		}
		advancedHttpHealthCheck := envoy_advanced_http.AdvancedHttp{
			HttpHealthCheck:    glooHc,
			ResponseAssertions: convertGlooToEnvoyRespAssertions(hcCfg.GetHttpHealthCheck().GetResponseAssertions()),
		}

		serializedAny, err := utils.MessageToAny(&advancedHttpHealthCheck)
		if err != nil {
			return err
		}

		out.HealthChecks[i] = envoyHc
		out.HealthChecks[i].HealthChecker = &envoy_config_core_v3.HealthCheck_CustomHealthCheck_{
			CustomHealthCheck: &envoy_config_core_v3.HealthCheck_CustomHealthCheck{
				Name: HealthCheckerName,
				ConfigType: &envoy_config_core_v3.HealthCheck_CustomHealthCheck_TypedConfig{
					TypedConfig: serializedAny,
				},
			},
		}
	}
	return nil
}

func convertGlooToEnvoyRespAssertions(assertions *gloo_advanced_http.ResponseAssertions) *envoy_advanced_http.ResponseAssertions {
	if assertions == nil {
		return nil
	}

	return &envoy_advanced_http.ResponseAssertions{
		ResponseMatchers: convertGlooResponseMatchersToEnvoy(assertions.ResponseMatchers),
		NoMatchHealth:    convertMatchHealthWithDefault(assertions.NoMatchHealth, envoy_advanced_http.HealthCheckResult_unhealthy),
	}
}

func convertMatchHealthWithDefault(mh gloo_advanced_http.HealthCheckResult, defaultHealth envoy_advanced_http.HealthCheckResult) envoy_advanced_http.HealthCheckResult {
	converted := defaultHealth

	switch mh {
	case gloo_advanced_http.HealthCheckResult_healthy:
		converted = envoy_advanced_http.HealthCheckResult_healthy
	case gloo_advanced_http.HealthCheckResult_degraded:
		converted = envoy_advanced_http.HealthCheckResult_degraded
	case gloo_advanced_http.HealthCheckResult_unhealthy:
		converted = envoy_advanced_http.HealthCheckResult_unhealthy
	}

	return converted
}

func convertGlooResponseMatchersToEnvoy(responseMatchers []*gloo_advanced_http.ResponseMatcher) []*envoy_advanced_http.ResponseMatcher {
	var respMatchers []*envoy_advanced_http.ResponseMatcher
	for _, rm := range responseMatchers {

		respMatcher := &envoy_advanced_http.ResponseMatcher{
			ResponseMatch: &envoy_advanced_http.ResponseMatch{
				JsonKey:            convertGlooJsonKeyToEnvoy(rm.GetResponseMatch().GetJsonKey()),
				IgnoreErrorOnParse: rm.GetResponseMatch().GetIgnoreErrorOnParse(),
				Regex:              rm.GetResponseMatch().GetRegex(),
			},
			MatchHealth: convertMatchHealthWithDefault(rm.MatchHealth, envoy_advanced_http.HealthCheckResult_healthy),
		}

		switch typed := rm.GetResponseMatch().GetSource().(type) {
		case *gloo_advanced_http.ResponseMatch_Header:
			respMatcher.ResponseMatch.Source = &envoy_advanced_http.ResponseMatch_Header{
				Header: typed.Header,
			}
		case *gloo_advanced_http.ResponseMatch_Body:
			respMatcher.ResponseMatch.Source = &envoy_advanced_http.ResponseMatch_Body{
				Body: typed.Body,
			}
		}

		respMatchers = append(respMatchers, respMatcher)
	}
	return respMatchers
}

func convertGlooJsonKeyToEnvoy(jsonKey *gloo_advanced_http.JsonKey) *envoy_advanced_http.JsonKey {
	if jsonKey == nil {
		return nil
	}

	var path []*envoy_advanced_http.JsonKey_PathSegment
	for _, ps := range jsonKey.Path {
		switch typed := ps.Segment.(type) {
		case *gloo_advanced_http.JsonKey_PathSegment_Key:
			segment := &envoy_advanced_http.JsonKey_PathSegment_Key{
				Key: typed.Key,
			}
			path = append(path, &envoy_advanced_http.JsonKey_PathSegment{
				Segment: segment,
			})
		}
	}
	return &envoy_advanced_http.JsonKey{
		Path: path,
	}
}

func convertEnvoyToGloo(ctx context.Context, httpHealth *envoy_config_core_v3.HealthCheck_HttpHealthCheck) (*envoy_core_v3_endpoint.HealthCheck_HttpHealthCheck, error) {
	if httpHealth == nil {
		return nil, errors.New("http health check is nil")
	}
	if len(httpHealth.GetPath()) == 0 {
		return nil, errors.New("http health check path is empty")
	}
	ret := &envoy_core_v3_endpoint.HealthCheck_HttpHealthCheck{
		Host: httpHealth.GetHost(), // ok if empty, defaults to name of cluster
		Path: httpHealth.GetPath(),
	}
	for _, st := range httpHealth.ExpectedStatuses {
		if st == nil {
			// slices in golang can contain nil values; protect against dev error although this should never be hit
			contextutils.LoggerFrom(ctx).DPanic("nil value in expected statuses slice")
			continue
		}
		ret.ExpectedStatuses = append(ret.ExpectedStatuses, &envoy_type_v3.Int64Range{
			Start: st.GetStart(),
			End:   st.GetEnd(),
		})
	}
	for _, rh := range httpHealth.RequestHeadersToAdd {
		if rh.GetHeader() == nil {
			// slices in golang can contain nil values; protect against dev error although this should never be hit
			contextutils.LoggerFrom(ctx).DPanic("nil value in request headers to add slice")
			continue
		}

		appendValue := &wrappers.BoolValue{Value: true}
		// The `httpHealth` parameter was generated from our internal config, so we only care about what configuration we handle during translation.
		// From Gloo to Envoy, we translate `append: true` to `append_if_exists_or_add` and `append: false` to `overwrite_if_exists_or_add`, so here we're reversing that translation.
		if rh.GetAppendAction() == envoy_config_core_v3.HeaderValueOption_APPEND_IF_EXISTS_OR_ADD {
			appendValue = &wrappers.BoolValue{Value: true}
		} else if rh.GetAppendAction() == envoy_config_core_v3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD {
			appendValue = &wrappers.BoolValue{Value: false}
		}
		ret.RequestHeadersToAdd = append(ret.RequestHeadersToAdd, &envoy_core_v3_endpoint.HeaderValueOption{
			Header: &envoy_core_v3_endpoint.HeaderValue{
				Key:   rh.GetHeader().GetKey(),
				Value: rh.GetHeader().GetValue(),
			},
			Append: appendValue,
		})
	}
	ret.RequestHeadersToRemove = httpHealth.GetRequestHeadersToRemove()
	ret.CodecClientType = envoy_type_v3.CodecClientType(httpHealth.GetCodecClientType())
	if mp := httpHealth.GetServiceNameMatcher().GetMatchPattern(); mp != nil {
		ret.ServiceNameMatcher = &envoy_type_matcher_v3_solo.StringMatcher{
			IgnoreCase: httpHealth.GetServiceNameMatcher().GetIgnoreCase(),
		}
		switch pattern := mp.(type) {
		case *envoy_type_matcher_v3.StringMatcher_Exact:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_Exact{
				Exact: pattern.Exact,
			}
		case *envoy_type_matcher_v3.StringMatcher_Prefix:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_Prefix{
				Prefix: pattern.Prefix,
			}
		case *envoy_type_matcher_v3.StringMatcher_SafeRegex:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_SafeRegex{
				SafeRegex: &envoy_type_matcher_v3_solo.RegexMatcher{
					EngineType: &envoy_type_matcher_v3_solo.RegexMatcher_GoogleRe2{GoogleRe2: &envoy_type_matcher_v3_solo.RegexMatcher_GoogleRE2{
						MaxProgramSize: pattern.SafeRegex.GetGoogleRe2().GetMaxProgramSize(),
					}},
					Regex: pattern.SafeRegex.GetRegex(),
				},
			}
		case *envoy_type_matcher_v3.StringMatcher_Suffix:
			ret.ServiceNameMatcher.MatchPattern = &envoy_type_matcher_v3_solo.StringMatcher_Suffix{
				Suffix: pattern.Suffix,
			}
		default:
			return nil, fmt.Errorf("unknown match pattern type %T", pattern)
		}
	}
	return ret, nil
}
