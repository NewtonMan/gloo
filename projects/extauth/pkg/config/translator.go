package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/solo-io/ext-auth-service/pkg/config/hmac"
	"github.com/solo-io/ext-auth-service/pkg/utils/cipher"

	jwtextauth "github.com/solo-io/ext-auth-service/pkg/config/jwt"

	"github.com/solo-io/ext-auth-service/pkg/config/oauth2"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys/externalstorage"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys/secrets"
	"github.com/solo-io/ext-auth-service/pkg/config/utils/jwks"
	"github.com/solo-io/ext-auth-service/pkg/controller/translation"

	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/rotisserie/eris"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/ext-auth-service/pkg/chain"
	"github.com/solo-io/ext-auth-service/pkg/config"
	"github.com/solo-io/ext-auth-service/pkg/config/apikeys"
	"github.com/solo-io/ext-auth-service/pkg/config/apr"
	"github.com/solo-io/ext-auth-service/pkg/config/ldap"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation/utils"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	grpcPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/grpc"
	httpPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/http"
	extRedis "github.com/solo-io/ext-auth-service/pkg/redis"
	"github.com/solo-io/ext-auth-service/pkg/session"
	redissession "github.com/solo-io/ext-auth-service/pkg/session/redis"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/go-utils/contextutils"
	extauthSoloApis "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	"go.uber.org/zap"
)

//go:generate mockgen -source ./translator.go -destination ./mocks/translator.go

type extAuthConfigTranslator struct {
	signingKey     []byte
	serviceFactory config.AuthServiceFactory
}

const defaultUseBearerSchemaForAuthorization = false

type ExtAuthConfigTranslator interface {
	Translate(ctx context.Context, resource *extauthv1.ExtAuthConfig) (svc api.AuthService, err error)
}

func NewTranslator(
	key []byte,
	serviceFactory config.AuthServiceFactory,
) ExtAuthConfigTranslator {
	return &extAuthConfigTranslator{
		signingKey:     key,
		serviceFactory: serviceFactory,
	}
}

func (t *extAuthConfigTranslator) Translate(ctx context.Context, resource *extauthv1.ExtAuthConfig) (svc api.AuthService, err error) {
	defer func() {
		if r := recover(); r != nil {
			svc = nil
			stack := string(debug.Stack())
			err = errors.Errorf("panicked while retrieving config for resource %v: %v %v", resource, r, stack)
		}
	}()

	contextutils.LoggerFrom(ctx).Debugw("Getting config for resource", zap.Any("resource", resource))

	if len(resource.Configs) != 0 {
		return t.getConfigs(ctx, resource.BooleanExpr.GetValue(), resource.Configs, resource.FailOnRedirect)
	}

	return nil, nil
}

func (t *extAuthConfigTranslator) getConfigs(
	ctx context.Context,
	boolLogic string,
	configs []*extauthv1.ExtAuthConfig_Config,
	failOnRedirect bool,
) (svc api.AuthService, err error) {

	services := chain.NewAuthServiceChain()
	services.SetFailOnRedirect(failOnRedirect)
	for i, cfg := range configs {
		svc, name, err := t.authConfigToService(ctx, cfg)
		if err != nil {
			return nil, err
		}
		if name == "" {
			name = fmt.Sprintf("config_%d", i)
		}
		if err := services.AddAuthService(name, svc); err != nil {
			return nil, err
		}
	}
	if strings.ContainsAny(boolLogic, "-+/*^%") {
		return nil, errors.New("auth config boolean logic contains an invalid character, do not use any of (-+/*^%) ")
	}
	if err = services.SetAuthorizer(boolLogic); err != nil {
		return nil, err
	}

	return services, nil
}

func (t *extAuthConfigTranslator) authConfigToService(
	ctx context.Context,
	eaconfig *extauthv1.ExtAuthConfig_Config,
) (svc api.AuthService, name string, err error) {
	switch cfg := eaconfig.AuthConfig.(type) {
	case *extauthv1.ExtAuthConfig_Config_Jwt:
		return &jwtextauth.JwtAuthService{}, eaconfig.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_BasicAuth:
		aprCfg := apr.Config{
			Realm:                            cfg.BasicAuth.Realm,
			SaltAndHashedPasswordPerUsername: convertAprUsers(cfg.BasicAuth.GetApr().GetUsers()),
		}

		return &aprCfg, eaconfig.GetName().GetValue(), nil

	// support deprecated config
	case *extauthv1.ExtAuthConfig_Config_Oauth:
		cb := cfg.Oauth.CallbackPath
		if cb == "" {
			cb = DefaultCallback
		}
		issuerUrl := addTrailingSlash(cfg.Oauth.IssuerUrl)

		authService, err := t.serviceFactory.NewOidcAuthorizationCodeAuthService(
			ctx,
			&config.OidcAuthorizationCodeAuthServiceParams{
				ClientId:                    cfg.Oauth.GetClientId(),
				ClientSecret:                cfg.Oauth.GetClientSecret(),
				IssuerUrl:                   issuerUrl,
				AppUri:                      cfg.Oauth.GetAppUrl(),
				CallbackPath:                cb,
				AuthEndpointQueryParams:     cfg.Oauth.GetAuthEndpointQueryParams(),
				Scopes:                      cfg.Oauth.GetScopes(),
				SessionParams:               oidc.SessionParameters{},
				Headers:                     &oidc.HeaderConfig{},
				DiscoveryDataOverride:       &oidc.DiscoveryData{},
				DiscoveryPollInterval:       DefaultOIDCDiscoveryPollInterval,
				InvalidJwksOnDemandStrategy: jwks.NewNilKeySourceFactory(),
			})

		if err != nil {
			return nil, eaconfig.GetName().GetValue(), err
		}
		return authService, eaconfig.GetName().GetValue(), nil

	case *extauthv1.ExtAuthConfig_Config_Oauth2:

		switch oauthCfg := cfg.Oauth2.OauthType.(type) {
		case *extauthv1.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode:
			oidcCfg := oauthCfg.OidcAuthorizationCode

			cb := oidcCfg.CallbackPath
			if cb == "" {
				cb = DefaultCallback
			}

			oidcCfg.IssuerUrl = addTrailingSlash(oidcCfg.IssuerUrl)
			var sessionParameters oidc.SessionParameters
			if oidcCfg.GetUserSession() != nil {
				sessionParameters, err = UserSessionConfigToSessionParameters(oidcCfg.GetUserSession())
			} else {
				// GetSession() is deprecated, but we still need to support it in case APIs are not in sync
				// GetSession() for backwards compatibility
				sessionParameters, err = ToSessionParameters(oidcCfg.GetSession())
			}
			if err != nil {
				return nil, eaconfig.GetName().GetValue(), err
			}

			headersConfig := ToHeaderConfig(oidcCfg.GetHeaders())
			if headersConfig == nil {
				headersConfig = &oidc.HeaderConfig{}
			}

			discoveryDataOverride := ToDiscoveryDataOverride(oidcCfg.GetDiscoveryOverride())
			if discoveryDataOverride == nil {
				discoveryDataOverride = &oidc.DiscoveryData{}
			}

			discoveryPollInterval := oidcCfg.GetDiscoveryPollInterval()
			if discoveryPollInterval == nil {
				discoveryPollInterval = ptypes.DurationProto(DefaultOIDCDiscoveryPollInterval)
			}

			autoMapFromMetadata := ToAutoMapFromMetadata(oidcCfg.GetAutoMapFromMetadata())
			if autoMapFromMetadata == nil {
				autoMapFromMetadata = &oidc.AutoMapFromMetadata{}
			}

			endSessionProperties := ToEndSessionEndpointProperties(oidcCfg.GetEndSessionProperties())

			jwksOnDemandCacheRefreshPolicy := ToOnDemandCacheRefreshPolicy(oidcCfg.GetJwksCacheRefreshPolicy())

			authService, err := t.serviceFactory.NewOidcAuthorizationCodeAuthService(
				ctx,
				&config.OidcAuthorizationCodeAuthServiceParams{
					ClientId:                    oidcCfg.GetClientId(),
					ClientSecret:                oidcCfg.GetClientSecret(),
					IssuerUrl:                   oidcCfg.GetIssuerUrl(),
					AppUri:                      oidcCfg.GetAppUrl(),
					CallbackPath:                cb,
					LogoutPath:                  oidcCfg.GetLogoutPath(),
					AfterLogoutUri:              oidcCfg.GetAfterLogoutUrl(),
					SessionIdHeaderName:         oidcCfg.GetSessionIdHeaderName(),
					AuthEndpointQueryParams:     oidcCfg.GetAuthEndpointQueryParams(),
					TokenEndpointQueryParams:    oidcCfg.GetTokenEndpointQueryParams(),
					Scopes:                      oidcCfg.GetScopes(),
					SessionParams:               sessionParameters,
					Headers:                     headersConfig,
					DiscoveryDataOverride:       discoveryDataOverride,
					DiscoveryPollInterval:       discoveryPollInterval.AsDuration(),
					InvalidJwksOnDemandStrategy: jwksOnDemandCacheRefreshPolicy,
					ParseCallbackPathAsRegex:    oidcCfg.GetParseCallbackPathAsRegex(),
					AutoMapFromMetadata:         autoMapFromMetadata,
					EndSessionProperties:        endSessionProperties,
				})

			if err != nil {
				return nil, eaconfig.GetName().GetValue(), err
			}
			return authService, eaconfig.GetName().GetValue(), nil

		case *extauthv1.ExtAuthConfig_OAuth2Config_AccessTokenValidationConfig:
			userInfoUrl := oauthCfg.AccessTokenValidationConfig.GetUserinfoUrl()
			scopeValidator := utils.NewMatchAllValidator(oauthCfg.AccessTokenValidationConfig.GetRequiredScopes().GetScope())

			cacheTtl := oauthCfg.AccessTokenValidationConfig.CacheTimeout
			if cacheTtl == nil {
				cacheTtl = ptypes.DurationProto(DefaultOAuthCacheTtl)
			}

			switch validationType := oauthCfg.AccessTokenValidationConfig.GetValidationType().(type) {
			case *extauthv1.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionUrl:
				authService := t.serviceFactory.NewOAuth2TokenIntrospectionAuthService(
					&config.OAuth2TokenIntrospectionAuthServiceParams{
						IntrospectionUrl: validationType.IntrospectionUrl,
						ScopeValidator:   scopeValidator,
						UserInfoUrl:      userInfoUrl,
						CacheTtl:         cacheTtl.AsDuration(),
					})
				return authService, eaconfig.GetName().GetValue(), nil
			case *extauthv1.ExtAuthConfig_AccessTokenValidationConfig_Introspection:
				authService := t.serviceFactory.NewOAuth2TokenIntrospectionAuthService(
					&config.OAuth2TokenIntrospectionAuthServiceParams{
						ClientId:         validationType.Introspection.GetClientId(),
						ClientSecret:     validationType.Introspection.GetClientSecret(),
						IntrospectionUrl: validationType.Introspection.GetIntrospectionUrl(),
						ScopeValidator:   scopeValidator,
						UserInfoUrl:      userInfoUrl,
						CacheTtl:         cacheTtl.AsDuration(),
						UserIdAttribute:  validationType.Introspection.GetUserIdAttributeName(),
					})
				return authService, eaconfig.GetName().GetValue(), nil

			case *extauthv1.ExtAuthConfig_AccessTokenValidationConfig_Jwt:
				authService, err := t.serviceFactory.NewOAuth2JwtAccessTokenAuthService(
					ctx,
					&config.OAuth2JwtAccessTokenAuthServiceParams{
						JwksStr:                   validationType.Jwt.GetLocalJwks().GetInlineString(),
						RemoteJwksUri:             validationType.Jwt.GetRemoteJwks().GetUrl(),
						RemoteJwksRefreshInterval: validationType.Jwt.GetRemoteJwks().GetRefreshInterval().AsDuration(),
						Issuer:                    validationType.Jwt.GetIssuer(),
						ScopeValidator:            scopeValidator,
						UserInfoUrl:               userInfoUrl,
						CacheTtl:                  cacheTtl.AsDuration(),
					})
				if err != nil {
					return nil, "", err
				}
				return authService, eaconfig.GetName().GetValue(), nil

			default:
				return nil, eaconfig.GetName().GetValue(), errors.Errorf("Unhandled access token validation type: %+v", oauthCfg.AccessTokenValidationConfig.ValidationType)
			}
		case *extauthv1.ExtAuthConfig_OAuth2Config_Oauth2Config:
			plainOAuth2Cfg := oauthCfg.Oauth2Config

			cb := plainOAuth2Cfg.GetCallbackPath()
			if cb == "" {
				cb = DefaultCallback
			}

			sessionIdHeader := ""
			var sessionParameters oauth2.SessionParameters
			if plainOAuth2Cfg.GetUserSession() != nil {
				sessionParameters, err = UserSessionConfigToSessionParametersOAuth2(plainOAuth2Cfg.GetUserSession())
				if redisSession := plainOAuth2Cfg.GetUserSession().GetRedis(); redisSession != nil {
					sessionIdHeader = redisSession.GetHeaderName()
				}
			} else {
				// GetSession() is deprecated, but we still need to support it in case APIs are not in sync
				// GetSession() for backwards compatibility
				sessionParameters, err = ToSessionParametersOAuth2(plainOAuth2Cfg.GetSession())
				if redisSession := plainOAuth2Cfg.GetSession().GetRedis(); redisSession != nil {
					sessionIdHeader = redisSession.GetHeaderName()
				}
			}
			if err != nil {
				return nil, eaconfig.GetName().GetValue(), err
			}

			authService, err := t.serviceFactory.NewPlainOAuth2AuthService(
				ctx,
				&config.PlainOAuth2AuthServiceParams{
					ClientId:                 plainOAuth2Cfg.GetClientId(),
					ClientSecret:             plainOAuth2Cfg.GetClientSecret(),
					AppUri:                   plainOAuth2Cfg.GetAppUrl(),
					CallbackPath:             cb,
					LogoutPath:               plainOAuth2Cfg.GetLogoutPath(),
					AfterLogoutUri:           plainOAuth2Cfg.GetAfterLogoutUrl(),
					SessionIdHeaderName:      sessionIdHeader,
					AuthEndpointQueryParams:  plainOAuth2Cfg.GetAuthEndpointQueryParams(),
					TokenEndpointQueryParams: plainOAuth2Cfg.GetTokenEndpointQueryParams(),
					Scopes:                   plainOAuth2Cfg.GetScopes(),
					SessionParams:            sessionParameters,
					AuthEndpoint:             plainOAuth2Cfg.GetAuthEndpoint(),
					TokenEndpoint:            plainOAuth2Cfg.GetTokenEndpoint(),
					RevocationEndpoint:       plainOAuth2Cfg.GetRevocationEndpoint(),
				})

			return authService, eaconfig.GetName().GetValue(), err
		}

	case *extauthv1.ExtAuthConfig_Config_ApiKeyAuth:
		switch cfg.ApiKeyAuth.GetStorageBackend().(type) {
		case *extauthv1.ExtAuthConfig_ApiKeyAuthConfig_K8SSecretApikeyStorage:
			{
				validApiKeys := map[string]apikeys.KeyMetadata{}
				for apiKey, metadata := range cfg.ApiKeyAuth.ValidApiKeys {
					validApiKeys[apiKey] = apikeys.KeyMetadata{
						UserName: metadata.Username,
						Metadata: metadata.Metadata,
					}
				}
				secretsConf := &secrets.Config{
					ApiKeyHeaderName:       cfg.ApiKeyAuth.GetHeaderName(),
					HeadersFromKeyMetadata: cfg.ApiKeyAuth.GetHeadersFromKeyMetadata(),
					ValidApiKeys:           validApiKeys,
					LabelSelector:          cfg.ApiKeyAuth.GetK8SSecretApikeyStorage().GetLabelSelector(),
					ApiKeySecretRefs:       cfg.ApiKeyAuth.GetK8SSecretApikeyStorage().GetApiKeySecretRefs(),
				}
				apiKeyAuthService := secrets.NewAPIKeyService(secretsConf)
				return apiKeyAuthService, eaconfig.GetName().GetValue(), nil
			}
		case *extauthv1.ExtAuthConfig_ApiKeyAuthConfig_AerospikeApikeyStorage:
			{
				aerospikeConf, err := TranslateAerospikeConfig(cfg)
				apiKeyAuthService, err := externalstorage.NewAPIKeyService(ctx, aerospikeConf)
				if err != nil {
					return nil, eaconfig.GetName().GetValue(), err
				}
				return apiKeyAuthService, eaconfig.GetName().GetValue(), nil
			}
		default:
			{
				validApiKeys := map[string]apikeys.KeyMetadata{}
				for apiKey, metadata := range cfg.ApiKeyAuth.ValidApiKeys {
					validApiKeys[apiKey] = apikeys.KeyMetadata{
						UserName: metadata.Username,
						Metadata: metadata.Metadata,
					}
				}
				secretsConf := &secrets.Config{
					ApiKeyHeaderName:       cfg.ApiKeyAuth.GetHeaderName(),
					HeadersFromKeyMetadata: cfg.ApiKeyAuth.GetHeadersFromKeyMetadata(),
					ValidApiKeys:           validApiKeys,
				}
				apiKeyAuthService := secrets.NewAPIKeyService(secretsConf)
				return apiKeyAuthService, eaconfig.GetName().GetValue(), nil
			}
		}

	case *extauthv1.ExtAuthConfig_Config_PluginAuth:
		p, err := t.serviceFactory.LoadAuthPlugin(ctx, cfg.PluginAuth)
		return p, cfg.PluginAuth.Name, err // plugin name takes precedent over auth config name
	case *extauthv1.ExtAuthConfig_Config_OpaAuth:
		options := opa.Options{
			FastInputConversion: cfg.OpaAuth.GetOptions().GetFastInputConversion(),
		}
		opaCfg, err := opa.NewWithOptions(ctx, cfg.OpaAuth.Query, cfg.OpaAuth.Modules, options)
		if err != nil {
			return nil, "", err
		}
		return opaCfg, eaconfig.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_LdapInternal:
		ldapSvc, err := getLdapAuthServiceWithSecret(ctx, cfg.LdapInternal)
		if err != nil {
			return nil, "", err
		}
		return ldapSvc, eaconfig.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_Ldap:
		ldapSvc, err := getLdapAuthService(ctx, cfg.Ldap)
		if err != nil {
			return nil, "", err
		}
		return ldapSvc, eaconfig.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_HmacAuth:
		passwords := cfg.HmacAuth.GetSecretList().GetSecretList()
		// When we add multiple implementations, there will need to be a check here for the implmentation type, but for now it's always HeadersRequestToHmac
		hmacConversionFunc := hmac.HeadersRequestToHmac
		hmacSvc := t.serviceFactory.NewHmacAuthService(
			&config.HmacAuthServiceParams{
				HmacPasswords: passwords,
				Unwrapper:     hmacConversionFunc,
			})
		return hmacSvc, eaconfig.GetName().GetValue(), nil
	case *extauthv1.ExtAuthConfig_Config_PassThroughAuth:
		switch protocolConfig := cfg.PassThroughAuth.GetProtocol().(type) {
		case *extauthv1.PassThroughAuth_Grpc:
			grpcSvc, err := getPassThroughGrpcAuthService(ctx, cfg.PassThroughAuth.GetConfig(), protocolConfig.Grpc, cfg.PassThroughAuth.GetFailureModeAllow())
			if err != nil {
				return nil, "", err
			}
			return grpcSvc, eaconfig.GetName().GetValue(), nil
		case *extauthv1.PassThroughAuth_Http:
			svc, err := getPassThroughHttpService(ctx, cfg.PassThroughAuth.GetConfig(), protocolConfig.Http, cfg.PassThroughAuth.GetFailureModeAllow())
			if err != nil {
				return nil, "", err
			}
			return svc, eaconfig.GetName().GetValue(), nil
		default:
			return nil, eaconfig.GetName().GetValue(), errors.Errorf("Unhandled pass through auth protocol: %+v", cfg.PassThroughAuth.Protocol)
		}

	}
	return nil, "", errors.New("unknown auth configuration")
}

func addTrailingSlash(url string) string {
	if len(url) != 0 && url[len(url)-1:] == "/" {
		return url
	}
	return url + "/"
}

// TranslateAerospikeConfig receives an ExtAuthConfig and produces a ServiceStorageConfig configured with Aerospike connection details
// This function is extracted to unit test Aerospike config translation since the call to NewAPIKeyService now attempts to
// establish the database connection, where previously this was the responsibility of the Start method.
func TranslateAerospikeConfig(cfg *extauthv1.ExtAuthConfig_Config_ApiKeyAuth) (*externalstorage.ServiceStorageConfig, error) {
	inConf := cfg.ApiKeyAuth.GetAerospikeApikeyStorage()
	soloApisAerospikeConf := &extauthSoloApis.AerospikeApiKeyStorage{
		Hostname:      cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetHostname(),
		Namespace:     cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetNamespace(),
		Set:           cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetSet(),
		Port:          cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetPort(),
		BatchSize:     cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetBatchSize(),
		NodeTlsName:   cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetNodeTlsName(),
		CertPath:      cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetCertPath(),
		KeyPath:       cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetKeyPath(),
		AllowInsecure: cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetAllowInsecure(),
		RootCaPath:    cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetRootCaPath(),
		TlsVersion:    cfg.ApiKeyAuth.GetAerospikeApikeyStorage().GetTlsVersion(),
	}

	if _, ok := inConf.GetCommitLevel().(*extauthv1.AerospikeApiKeyStorage_CommitMaster); ok {
		soloApisAerospikeConf.CommitLevel = &extauthSoloApis.AerospikeApiKeyStorage_CommitMaster{}
	} else {
		soloApisAerospikeConf.CommitLevel = &extauthSoloApis.AerospikeApiKeyStorage_CommitAll{}
	}

	switch inConf.GetReadModeSc().GetReadModeSc().(type) {
	case *extauthv1.AerospikeApiKeyStorageReadModeSc_ReadModeScAllowUnavailable:
		soloApisAerospikeConf.ReadModeSc = &extauthSoloApis.AerospikeApiKeyStorageReadModeSc{
			ReadModeSc: &extauthSoloApis.AerospikeApiKeyStorageReadModeSc_ReadModeScAllowUnavailable{},
		}
	case *extauthv1.AerospikeApiKeyStorageReadModeSc_ReadModeScLinearize:
		soloApisAerospikeConf.ReadModeSc = &extauthSoloApis.AerospikeApiKeyStorageReadModeSc{
			ReadModeSc: &extauthSoloApis.AerospikeApiKeyStorageReadModeSc_ReadModeScLinearize{},
		}
	case *extauthv1.AerospikeApiKeyStorageReadModeSc_ReadModeScReplica:
		soloApisAerospikeConf.ReadModeSc = &extauthSoloApis.AerospikeApiKeyStorageReadModeSc{
			ReadModeSc: &extauthSoloApis.AerospikeApiKeyStorageReadModeSc_ReadModeScReplica{},
		}
	default:
		soloApisAerospikeConf.ReadModeSc = &extauthSoloApis.AerospikeApiKeyStorageReadModeSc{
			ReadModeSc: &extauthSoloApis.AerospikeApiKeyStorageReadModeSc_ReadModeScSession{},
		}
	}

	switch inConf.GetReadModeAp().GetReadModeAp().(type) {
	case *extauthv1.AerospikeApiKeyStorageReadModeAp_ReadModeApAll:
		soloApisAerospikeConf.ReadModeAp = &extauthSoloApis.AerospikeApiKeyStorageReadModeAp{
			ReadModeAp: &extauthSoloApis.AerospikeApiKeyStorageReadModeAp_ReadModeApAll{},
		}
	default:
		soloApisAerospikeConf.ReadModeAp = &extauthSoloApis.AerospikeApiKeyStorageReadModeAp{
			ReadModeAp: &extauthSoloApis.AerospikeApiKeyStorageReadModeAp_ReadModeApOne{},
		}
	}

	for _, tlsCurveGroup := range inConf.GetTlsCurveGroups() {
		switch tlsCurveGroup.GetCurveId().(type) {
		case *extauthv1.AerospikeApiKeyStorageTlsCurveID_CurveP256:
			soloApisAerospikeConf.TlsCurveGroups =
				append(soloApisAerospikeConf.TlsCurveGroups,
					&extauthSoloApis.AerospikeApiKeyStorageTlsCurveID{
						CurveId: &extauthSoloApis.AerospikeApiKeyStorageTlsCurveID_CurveP256{},
					})
		case *extauthv1.AerospikeApiKeyStorageTlsCurveID_CurveP384:
			soloApisAerospikeConf.TlsCurveGroups =
				append(soloApisAerospikeConf.TlsCurveGroups,
					&extauthSoloApis.AerospikeApiKeyStorageTlsCurveID{
						CurveId: &extauthSoloApis.AerospikeApiKeyStorageTlsCurveID_CurveP384{},
					})
		case *extauthv1.AerospikeApiKeyStorageTlsCurveID_CurveP521:
			soloApisAerospikeConf.TlsCurveGroups =
				append(soloApisAerospikeConf.TlsCurveGroups,
					&extauthSoloApis.AerospikeApiKeyStorageTlsCurveID{
						CurveId: &extauthSoloApis.AerospikeApiKeyStorageTlsCurveID_CurveP521{},
					})
		case *extauthv1.AerospikeApiKeyStorageTlsCurveID_X_25519:
			soloApisAerospikeConf.TlsCurveGroups =
				append(soloApisAerospikeConf.TlsCurveGroups,
					&extauthSoloApis.AerospikeApiKeyStorageTlsCurveID{
						CurveId: &extauthSoloApis.AerospikeApiKeyStorageTlsCurveID_X_25519{},
					})
		default:
			return nil, eris.New("invalid tls curve id")
		}
	}

	return &externalstorage.ServiceStorageConfig{
		StorageConfig: &externalstorage.StorageConfig{
			AerospikeStorageConfig: soloApisAerospikeConf,
			BackendType:            apikeys.AerospikeBackend,
		},
		ServiceConfig: &externalstorage.ServiceConfig{
			ApiKeyHeaderName:       cfg.ApiKeyAuth.GetHeaderName(),
			HeadersFromKeyMetadata: cfg.ApiKeyAuth.GetHeadersFromKeyMetadata(),
			// When Aerospike is used as the persistence layer for Secrets,
			// the AuthService will be responsible for filtering valid keys
			// based on the labelSelector defined on the AuthConfig.
			// This is a _key difference_ from the k8s Secret implementation.
			// When K8s Secrets are used, the labelSelector is processed at translation time
			LabelSelector: sanitizeAerospikeLabelSelector(inConf.GetLabelSelector()),
		},
	}, nil
}

// developerPortalLabelMapping is responsible for ensuring that the labelSelector which is passed to the AuthService,
// is formatted properly. There are 2 ways that AuthConfigs are generated: by users, by the Developer Portal
// This function is intended to sanitize configuration generated by the Developer Portal.
// THIS IS KNOWN TECHNICAL DEBT
// We understand long-term that the Developer Portal should not configure this, but for now we sanitize it here
// because to our understanding, the Aerospike backend is only used by customers who rely on the Developer Portal.
//
// DevPortal-generated labelSelector:
//
//	apiproducts.portal.gloo.solo.io: petstore-product.default
//	environments.portal.gloo.solo.io: dev.default
//	usageplans.portal.gloo.solo.io: basic
//
// Desired labelSelector:
//
//	product: petstore-product.default
//	environment: dev.default
//	plan: basic
var developerPortalLabelMapping = map[string]string{
	"apiproducts.portal.gloo.solo.io":  "product",
	"environments.portal.gloo.solo.io": "environment",
	"usageplans.portal.gloo.solo.io":   "plan",
}

func sanitizeAerospikeLabelSelector(labelSelector map[string]string) map[string]string {
	sanitizedLabelSelector := make(map[string]string, len(labelSelector))

	for key, value := range labelSelector {
		if newKey, ok := developerPortalLabelMapping[key]; ok {
			sanitizedLabelSelector[newKey] = value
		} else {
			sanitizedLabelSelector[key] = value
		}

	}
	return sanitizedLabelSelector
}

func getLdapAuthService(ctx context.Context, ldapCfg *extauthv1.Ldap) (api.AuthService, error) {
	poolInitCap, poolMaxCap := getLdapConnectionPoolParams(ldapCfg.GetPool())

	// Connection pool will be cleaned up when the context is cancelled
	ldapClientBuilder, err := ldap.NewPooledClientBuilder(ctx, &ldap.ClientPoolConfig{
		ServerAddress:   ldapCfg.Address,
		InitialCapacity: poolInitCap,
		MaximumCapacity: poolMaxCap,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start LDAP connection pool")
	}
	ldapAuthService, err := ldap.NewLdapAuthService(ldapClientBuilder, &ldap.Config{
		UserDnTemplate:                ldapCfg.UserDnTemplate,
		MembershipAttributeName:       ldapCfg.MembershipAttributeName,
		AllowedGroups:                 ldapCfg.AllowedGroups,
		SearchFilter:                  ldapCfg.SearchFilter,
		DisableGroupChecking:          ldapCfg.DisableGroupChecking,
		CheckGroupsWithServiceAccount: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create LDAP auth service")
	}
	return ldapAuthService, nil
}
func getLdapAuthServiceWithSecret(ctx context.Context, ldapCfg *extauthv1.ExtAuthConfig_LdapConfig) (api.AuthService, error) {
	poolInitCap, poolMaxCap := getLdapConnectionPoolParams(ldapCfg.GetPool())

	// Connection pool will be cleaned up when the context is cancelled
	ldapClientBuilder, err := ldap.NewPooledClientBuilder(ctx, &ldap.ClientPoolConfig{
		ServerAddress:   ldapCfg.Address,
		InitialCapacity: poolInitCap,
		MaximumCapacity: poolMaxCap,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to start LDAP connection pool")
	}
	ldapAuthService, err := ldap.NewLdapAuthService(ldapClientBuilder, &ldap.Config{
		UserDnTemplate:                ldapCfg.UserDnTemplate,
		MembershipAttributeName:       ldapCfg.MembershipAttributeName,
		AllowedGroups:                 ldapCfg.AllowedGroups,
		SearchFilter:                  ldapCfg.SearchFilter,
		DisableGroupChecking:          ldapCfg.DisableGroupChecking,
		CheckGroupsWithServiceAccount: ldapCfg.GetGroupLookupSettings().GetCheckGroupsWithServiceAccount(),
		ServiceAccountUsername:        ldapCfg.GetGroupLookupSettings().GetUsername(),
		ServiceAccountPassword:        ldapCfg.GetGroupLookupSettings().GetPassword(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create LDAP auth service")
	}
	return ldapAuthService, nil
}

func getLdapConnectionPoolParams(pool *extauthv1.Ldap_ConnectionPool) (int, int) {
	initCap := 2
	maxCap := 5

	if initSize := pool.GetInitialSize(); initSize != nil {
		initCap = int(initSize.Value)
	}

	if maxSize := pool.GetMaxSize(); maxSize != nil {
		maxCap = int(maxSize.Value)
	}
	return initCap, maxCap
}

func getPassThroughGrpcAuthService(ctx context.Context, passthroughAuthCfg *structpb.Struct, grpcConfig *extauthv1.PassThroughGrpc, failureModeAllow bool) (api.AuthService, error) {
	connectionTimeout := 5 * time.Second

	if timeout := grpcConfig.GetConnectionTimeout(); timeout != nil {
		timeout, err := ptypes.Duration(timeout)
		if err != nil {
			return nil, err
		}
		connectionTimeout = timeout
	}

	clientManagerConfig := &grpcPassthrough.ClientManagerConfig{
		Address:           grpcConfig.GetAddress(),
		ConnectionTimeout: connectionTimeout,
		FailureModeAllow:  failureModeAllow,
		UseSecure:         grpcConfig.GetTlsConfig() != nil,
	}

	grpcClientManager, err := grpcPassthrough.NewGrpcClientManager(ctx, clientManagerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client manager")
	}

	return grpcPassthrough.NewGrpcService(grpcClientManager, passthroughAuthCfg), nil
}

func getPassThroughHttpService(ctx context.Context, authCfgCfg *structpb.Struct, httpPassthroughConfig *extauthv1.PassThroughHttp, failureModeAllow bool) (api.AuthService, error) {
	connectionTimeout := 5 * time.Second
	if timeout := httpPassthroughConfig.GetConnectionTimeout(); timeout != nil {
		timeout, err := ptypes.Duration(timeout)
		if err != nil {
			return nil, err
		}
		connectionTimeout = timeout
	}

	allowedHeadersMap := map[string]bool{}
	for _, header := range httpPassthroughConfig.GetRequest().GetAllowedHeaders() {
		allowedHeadersMap[header] = true
	}

	var tlsConfig *tls.Config
	if rootCa := os.Getenv(translation.HttpsPassthroughCaCert); rootCa != "" {
		rootCaBytes, err := base64.StdEncoding.DecodeString(rootCa)
		if err != nil {
			return nil, errors.Wrapf(err, "error base64 decoding root ca %s", rootCa)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(rootCaBytes)
		if !ok {
			return nil, errors.Errorf("ca cert base64 encoded - (%s) is not OK", rootCa)
		}

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	cfg := &httpPassthrough.PassthroughConfig{
		PassThroughFilterMetadata:         httpPassthroughConfig.GetRequest().GetPassThroughFilterMetadata(),
		PassThroughState:                  httpPassthroughConfig.GetRequest().GetPassThroughState(),
		PassThroughBody:                   httpPassthroughConfig.GetRequest().GetPassThroughBody(),
		AllowedHeaders:                    allowedHeadersMap,
		HeadersToAdd:                      httpPassthroughConfig.GetRequest().GetHeadersToAdd(),
		Url:                               httpPassthroughConfig.Url,
		ConnectionTimeout:                 connectionTimeout,
		AllowedUpstreamHeaders:            httpPassthroughConfig.GetResponse().GetAllowedUpstreamHeaders(),
		AllowedUpstreamHeadersToOverwrite: httpPassthroughConfig.GetResponse().GetAllowedUpstreamHeadersToOverwrite(),
		AllowedClientHeaders:              httpPassthroughConfig.GetResponse().GetAllowedClientHeadersOnDenied(),
		ReadStateFromResponse:             httpPassthroughConfig.GetResponse().GetReadStateFromResponse(),
		TLSClientConfig:                   tlsConfig,
		AllowFailure:                      failureModeAllow,
	}

	return httpPassthrough.NewHttpService(cfg, authCfgCfg), nil
}

func convertAprUsers(users map[string]*extauthv1.BasicAuth_Apr_SaltedHashedPassword) map[string]apr.SaltAndHashedPassword {
	ret := map[string]apr.SaltAndHashedPassword{}
	for k, v := range users {
		ret[k] = apr.SaltAndHashedPassword{
			HashedPassword: v.HashedPassword,
			Salt:           v.Salt,
		}
	}
	return ret
}

// sessionToStore will create a session store based off the user session configuration.
// If the user session is nil, it will create no store.
// If it is redis, it will create a redis client store.
// If it is a cookie, it will create a cookie store.
// throws an error if the user session is unknown.
func sessionToStore(us *extauthv1.ExtAuthConfig_UserSessionConfig) (*translation.StoreParameters, error) {
	usersession := us.GetSession()
	if usersession == nil {
		return nil, nil
	}

	encryptionKey := us.GetCipherConfig().GetKey()
	var err error
	var encryptionCipher cipher.Cipher
	if encryptionKey != "" {
		encryptionCipher, err = cipher.NewGCMEncryption([]byte(encryptionKey))
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create encryption cipher for user session")
		}
	}

	switch s := usersession.(type) {
	case *extauthv1.ExtAuthConfig_UserSessionConfig_Cookie:
		allowRefreshing := false
		if allowRefreshSetting := s.Cookie.GetAllowRefreshing(); allowRefreshSetting != nil {
			allowRefreshing = allowRefreshSetting.Value
		}
		store := translation.StoreParameters{
			Session:          oidc.NewCookieSessionStore(s.Cookie.GetKeyPrefix(), encryptionCipher),
			RefreshIfExpired: allowRefreshing,
			TargetDomain:     s.Cookie.GetTargetDomain(),
			PreExpiryBuffer:  0,
		}
		return &store, nil
	case *extauthv1.ExtAuthConfig_UserSessionConfig_Redis:
		options := s.Redis.GetOptions()
		client, err := extRedis.NewRedisUniversalClient(getSoloApisRedisOptions(options))
		// there is an error creating the TLS Config
		if err != nil {
			return nil, err
		}
		rs := redissession.NewRedisSession(client, s.Redis.CookieName, s.Redis.KeyPrefix)

		allowRefreshing := true
		if allowRefreshSetting := s.Redis.AllowRefreshing; allowRefreshSetting != nil {
			allowRefreshing = allowRefreshSetting.Value
		}

		preExpiryBuffer := &duration.Duration{Seconds: 2, Nanos: 0}
		if preExpiryBufferSetting := s.Redis.PreExpiryBuffer; preExpiryBufferSetting != nil {
			preExpiryBuffer = preExpiryBufferSetting
		}
		store := translation.StoreParameters{
			Session:          rs,
			RefreshIfExpired: allowRefreshing,
			PreExpiryBuffer:  preExpiryBuffer.AsDuration(),
			TargetDomain:     s.Redis.GetTargetDomain(),
		}
		return &store, nil
	}
	return nil, fmt.Errorf("no matching session config")
}

// userSessionToStore is deprecated, use sessionToStore() instead.
func userSessionToStore(us *extauthv1.UserSession) (session.SessionStore, bool, time.Duration, error) {
	if us == nil {
		return nil, false, 0, nil
	}
	usersession := us.Session
	if usersession == nil {
		return nil, false, 0, nil
	}

	switch s := usersession.(type) {
	case *extauthv1.UserSession_Cookie:
		allowRefreshing := false
		if allowRefreshSetting := s.Cookie.GetAllowRefreshing(); allowRefreshSetting != nil {
			allowRefreshing = allowRefreshSetting.Value
		}
		return oidc.NewCookieSessionStore(s.Cookie.GetKeyPrefix(), nil), allowRefreshing, 0, nil
	case *extauthv1.UserSession_Redis:
		options := s.Redis.GetOptions()
		client, err := extRedis.NewRedisUniversalClient(getSoloApisRedisOptions(options))
		// there is an error creating the TLS Config
		if err != nil {
			return nil, false, 0, err
		}
		rs := redissession.NewRedisSession(client, s.Redis.CookieName, s.Redis.KeyPrefix)

		allowRefreshing := true
		if allowRefreshSetting := s.Redis.AllowRefreshing; allowRefreshSetting != nil {
			allowRefreshing = allowRefreshSetting.Value
		}

		preExpiryBuffer := &duration.Duration{Seconds: 2, Nanos: 0}
		if preExpiryBufferSetting := s.Redis.PreExpiryBuffer; preExpiryBufferSetting != nil {
			preExpiryBuffer = preExpiryBufferSetting
		}

		return rs, allowRefreshing, preExpiryBuffer.AsDuration(), nil
	}
	return nil, false, 0, fmt.Errorf("no matching session config")
}

func cookieConfigToSessionOptions(cookieOptions *extauthv1.UserSession_CookieOptions) *session.Options {
	var sessionOptions *session.Options
	if cookieOptions != nil {
		var path *string
		if pathFromOpt := cookieOptions.GetPath(); pathFromOpt != nil {
			tmp := pathFromOpt.Value
			path = &tmp
		}
		maxAge := defaultMaxAge
		if maxAgeConfig := cookieOptions.MaxAge; maxAgeConfig != nil {
			maxAge = int(maxAgeConfig.Value)
		}
		httpOnly := true
		if cookieOptions.GetHttpOnly() != nil {
			httpOnly = cookieOptions.GetHttpOnly().Value
		}

		sessionOptions = &session.Options{
			Path:     path,
			Domain:   cookieOptions.GetDomain(),
			HttpOnly: httpOnly,
			Secure:   !cookieOptions.GetNotSecure(),
			MaxAge:   maxAge,
			SameSite: session.SameSite(cookieOptions.SameSite),
		}
	}
	return sessionOptions
}

func ToHeaderConfig(hc *extauthv1.HeaderConfiguration) *oidc.HeaderConfig {
	var headersConfig *oidc.HeaderConfig
	if hc != nil {
		useBearerSchemaForAuthorization := defaultUseBearerSchemaForAuthorization
		if bearerSchemaWrapper := hc.GetUseBearerSchemaForAuthorization(); bearerSchemaWrapper != nil {
			useBearerSchemaForAuthorization = bearerSchemaWrapper.GetValue()
		}

		headersConfig = &oidc.HeaderConfig{
			IdTokenHeader:                   hc.GetIdTokenHeader(),
			AccessTokenHeader:               hc.GetAccessTokenHeader(),
			UseBearerSchemaForAuthorization: useBearerSchemaForAuthorization,
		}
	}
	return headersConfig
}

func ToDiscoveryDataOverride(discoveryOverride *extauthv1.DiscoveryOverride) *oidc.DiscoveryData {
	var discoveryDataOverride *oidc.DiscoveryData
	if discoveryOverride != nil {
		discoveryDataOverride = &oidc.DiscoveryData{
			// IssuerUrl is intentionally excluded as it cannot be overridden
			AuthEndpoint:       discoveryOverride.GetAuthEndpoint(),
			RevocationEndpoint: discoveryOverride.GetRevocationEndpoint(),
			EndSessionEndpoint: discoveryOverride.GetEndSessionEndpoint(),
			TokenEndpoint:      discoveryOverride.GetTokenEndpoint(),
			KeysUri:            discoveryOverride.GetJwksUri(),
			ResponseTypes:      discoveryOverride.GetResponseTypes(),
			Subjects:           discoveryOverride.GetSubjects(),
			IDTokenAlgs:        discoveryOverride.GetIdTokenAlgs(),
			Scopes:             discoveryOverride.GetScopes(),
			AuthMethods:        discoveryOverride.GetAuthMethods(),
			Claims:             discoveryOverride.GetClaims(),
		}
	}
	return discoveryDataOverride
}

// UserSessionConfigToSessionParameters sets the Session Parameters and the store
func UserSessionConfigToSessionParameters(userSession *extauthv1.ExtAuthConfig_UserSessionConfig) (oidc.SessionParameters, error) {
	if userSession == nil {
		return oidc.SessionParameters{}, nil
	}
	sessionOptions := cookieConfigToSessionOptions(userSession.GetCookieOptions())
	store, err := sessionToStore(userSession)
	if err != nil {
		return oidc.SessionParameters{}, err
	}
	var sessionParams oidc.SessionParameters
	if store != nil {
		sessionParams = oidc.SessionParameters{
			Store:            store.Session,
			RefreshIfExpired: store.RefreshIfExpired,
			PreExpiryBuffer:  store.PreExpiryBuffer,
			TargetDomain:     store.TargetDomain,
		}
	}
	sessionParams.Options = sessionOptions
	sessionParams.ErrOnSessionFetch = userSession.GetFailOnFetchFailure()
	return sessionParams, nil
}

// toSessionParameters is deprecated: use UserSessionConfigToSessionParameters instead.
func ToSessionParameters(userSession *extauthv1.UserSession) (oidc.SessionParameters, error) {
	sessionOptions := cookieConfigToSessionOptions(userSession.GetCookieOptions())
	sessionStore, refreshIfExpired, preExpiryBuffer, err := userSessionToStore(userSession)
	if err != nil {
		return oidc.SessionParameters{}, err
	}
	return oidc.SessionParameters{
		ErrOnSessionFetch: userSession.GetFailOnFetchFailure(),
		Store:             sessionStore,
		Options:           sessionOptions,
		RefreshIfExpired:  refreshIfExpired,
		PreExpiryBuffer:   preExpiryBuffer,
	}, nil
}

// UserSessionConfigToSessionParametersOAuth2 will convert the user session config to a SessionParameters struct.
// configuration for the session store, cipher, and session options.
func UserSessionConfigToSessionParametersOAuth2(userSession *extauthv1.ExtAuthConfig_UserSessionConfig) (oauth2.SessionParameters, error) {
	if userSession == nil {
		return oauth2.SessionParameters{}, nil
	}
	sessionOptions := cookieConfigToSessionOptions(userSession.GetCookieOptions())
	store, err := sessionToStore(userSession)
	if err != nil {
		return oauth2.SessionParameters{}, err
	}

	var sessionParams oauth2.SessionParameters
	if store != nil {
		sessionParams = oauth2.SessionParameters{
			Store:            store.Session,
			RefreshIfExpired: store.RefreshIfExpired,
			PreExpiryBuffer:  store.PreExpiryBuffer,
			TargetDomain:     store.TargetDomain,
		}
	}
	sessionParams.ErrOnSessionFetch = userSession.GetFailOnFetchFailure()
	sessionParams.Options = sessionOptions
	return sessionParams, nil
}

// ToSessionParametersOAuth2 is deprecated use UserSessionConfigToSessionParametersOAuth2 instead.
func ToSessionParametersOAuth2(userSession *extauthv1.UserSession) (oauth2.SessionParameters, error) {
	sessionOptions := cookieConfigToSessionOptions(userSession.GetCookieOptions())
	sessionStore, refreshIfExpired, preExpiryBuffer, err := userSessionToStore(userSession)
	if err != nil {
		return oauth2.SessionParameters{}, err
	}
	return oauth2.SessionParameters{
		ErrOnSessionFetch: userSession.GetFailOnFetchFailure(),
		Store:             sessionStore,
		Options:           sessionOptions,
		RefreshIfExpired:  refreshIfExpired,
		PreExpiryBuffer:   preExpiryBuffer,
	}, nil
}

func ToOnDemandCacheRefreshPolicy(policy *extauthv1.JwksOnDemandCacheRefreshPolicy) jwks.KeySourceFactory {
	// The onDemandCacheRefreshPolicy determines how the JWKS cache should be refreshed when a request is made
	// that contains a key not contained in the JWKS cache
	switch cacheRefreshPolicy := policy.GetPolicy().(type) {
	case *extauthv1.JwksOnDemandCacheRefreshPolicy_Never:
		// Never refresh the cache on missing key
		return jwks.NewNilKeySourceFactory()

	case *extauthv1.JwksOnDemandCacheRefreshPolicy_Always:
		// Always refresh the cache on missing key
		return jwks.NewHttpKeySourceFactory(nil)

	case *extauthv1.JwksOnDemandCacheRefreshPolicy_MaxIdpReqPerPollingInterval:
		// Refresh the cache on missing key `MaxIdpReqPerPollingInterval` times per interval
		return jwks.NewMaxRequestHttpKeySourceFactory(nil, cacheRefreshPolicy.MaxIdpReqPerPollingInterval)
	}

	// The default case is Never refresh
	return jwks.NewNilKeySourceFactory()

}

func ToAutoMapFromMetadata(autoMapFromMetadata *extauthv1.AutoMapFromMetadata) *oidc.AutoMapFromMetadata {
	return oidc.NewAutoMapFromMetadata(autoMapFromMetadata.GetNamespace())
}

func getSoloApisRedisOptions(options *extauthv1.RedisOptions) *extauthSoloApis.RedisOptions {
	if options == nil {
		return nil
	}
	return &extauthSoloApis.RedisOptions{
		Host:             options.GetHost(),
		Db:               options.GetDb(),
		PoolSize:         options.GetPoolSize(),
		TlsCertMountPath: options.GetTlsCertMountPath(),
		SocketType:       extauthSoloApis.RedisOptions_SocketType(options.GetSocketType()),
	}
}

// ToEndSessionEndpointProperties translates from gloo to ext-auth-service
func ToEndSessionEndpointProperties(endSessionEndpointProperties *extauthv1.EndSessionProperties) *oidc.EndSessionProperties {
	return oidc.NewEndSessionProperties(oidc.EndSessionMethodType(endSessionEndpointProperties.GetMethodType()))
}
