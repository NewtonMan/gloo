package extauth

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/hashicorp/go-multierror"

	"github.com/rotisserie/eris"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/solo-io/ext-auth-service/pkg/config/opa"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

var (
	unknownConfigTypeError                 = errors.New("unknown ext auth configuration")
	emptyQueryError                        = errors.New("no query provided")
	noValidUsersError                      = errors.New("No valid users found")
	duplicateOidcClientAuthenticationError = errors.New("can not define both clientAuthentication and deprecated fields clientSecretRef/disableClientSecret")
	NonApiKeySecretError                   = func(secret *v1.Secret) error {
		return errors.Errorf("secret [%s] is not an API key secret", secret.Metadata.Ref().Key())
	}
	EmptyApiKeyError = func(secret *v1.Secret) error {
		return errors.Errorf("no API key found on API key secret [%s]", secret.Metadata.Ref().Key())
	}
	MissingRequiredMetadataError = func(requiredKey string, secret *v1.Secret) error {
		return errors.Errorf("API key secret [%s] does not contain the required [%s] metadata entry", secret.Metadata.Ref().Key(), requiredKey)
	}
	NonAccountCredentialsSecretError = func(secret *v1.Secret) error {
		return errors.Errorf("secret [%s] is not an Account Credentials secret", secret.Metadata.Ref().Key())
	}
	MissingSecretError = func() error {
		return errors.Errorf("Authenticating with service account configured without account credentials")
	}
	duplicateModuleError           = func(s string) error { return fmt.Errorf("%s is a duplicate module", s) }
	unknownPassThroughProtocolType = func(protocol interface{}) error {
		return errors.Errorf("unknown passthrough protocol type [%v]", protocol)
	}
	noMatchesForGroupError = func(labelSelector map[string]string) error {
		return errors.Errorf("no matching apikey secrets for the provided label selector %v", labelSelector)
	}
)

const (
	OidcPkJwtClientAuthValidForDefaultSeconds = 5
)

// TranslateExtAuthConfig Returns {nil, nil} if the input config is empty or if it contains only custom auth entries
// Deprecated: Prefer ConvertExternalAuthConfigToXdsAuthConfig
func TranslateExtAuthConfig(ctx context.Context, snapshot *v1snap.ApiSnapshot, authConfigRef *core.ResourceRef) (*extauth.ExtAuthConfig, error) {
	configResource, err := snapshot.AuthConfigs.Find(authConfigRef.Strings())
	if err != nil {
		return nil, errors.Errorf("could not find auth config [%s] in snapshot", authConfigRef.Key())
	}

	return ConvertExternalAuthConfigToXdsAuthConfig(ctx, snapshot, configResource)
}

// ConvertExternalAuthConfigToXdsAuthConfig converts a user facing AuthConfig object
// into an AuthConfig object that will be sent over xDS to the ext-auth-service.
// Returns {nil, nil} if the input config is empty
func ConvertExternalAuthConfigToXdsAuthConfig(ctx context.Context, snapshot *v1snap.ApiSnapshot, externalAuthConfig *extauth.AuthConfig) (*extauth.ExtAuthConfig, error) {
	var translatedConfigs []*extauth.ExtAuthConfig_Config
	for _, config := range externalAuthConfig.Configs {
		translated, err := translateConfig(ctx, snapshot, config)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to translate ext auth config")
		}

		translatedConfigs = append(translatedConfigs, translated)
	}

	if len(translatedConfigs) == 0 {
		return nil, nil
	}

	// We sort translatedConfigs to ensure that on each translation run the same configs
	// hash to the same value. This is one solution to ensure this.
	// However, we may choose a more robust way of implementing this
	sort.SliceStable(translatedConfigs, func(i, j int) bool {
		return translatedConfigs[i].GetName().GetValue() < translatedConfigs[j].GetName().GetValue()
	})

	return &extauth.ExtAuthConfig{
		BooleanExpr:       externalAuthConfig.GetBooleanExpr(),
		AuthConfigRefName: externalAuthConfig.GetMetadata().Ref().Key(),
		Configs:           translatedConfigs,
	}, nil
}

func translateConfig(ctx context.Context, snap *v1snap.ApiSnapshot, cfg *extauth.AuthConfig_Config) (*extauth.ExtAuthConfig_Config, error) {
	extAuthConfig := &extauth.ExtAuthConfig_Config{
		Name: cfg.Name,
	}

	switch config := cfg.AuthConfig.(type) {
	case *extauth.AuthConfig_Config_BasicAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_BasicAuth{
			BasicAuth: config.BasicAuth,
		}
	// handle deprecated case
	case *extauth.AuthConfig_Config_Oauth:
		cfg, err := translateOauth(snap, config.Oauth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth{
			Oauth: cfg,
		}
	case *extauth.AuthConfig_Config_Oauth2:

		switch oauthCfg := config.Oauth2.OauthType.(type) {
		case *extauth.OAuth2_OidcAuthorizationCode:
			cfg, err := translateOidcAuthorizationCode(snap, oauthCfg.OidcAuthorizationCode)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth2{
				Oauth2: &extauth.ExtAuthConfig_OAuth2Config{
					OauthType: &extauth.ExtAuthConfig_OAuth2Config_OidcAuthorizationCode{OidcAuthorizationCode: cfg},
				},
			}
		case *extauth.OAuth2_AccessTokenValidation:
			accessTokenValidationConfig, err := translateAccessTokenValidationConfig(snap, oauthCfg.AccessTokenValidation)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth2{
				Oauth2: &extauth.ExtAuthConfig_OAuth2Config{
					OauthType: &extauth.ExtAuthConfig_OAuth2Config_AccessTokenValidationConfig{
						AccessTokenValidationConfig: accessTokenValidationConfig,
					},
				},
			}
		case *extauth.OAuth2_Oauth2:
			plainOAuth2Config, err := translatePlainOAuth2(snap, oauthCfg.Oauth2)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Oauth2{
				Oauth2: &extauth.ExtAuthConfig_OAuth2Config{
					OauthType: &extauth.ExtAuthConfig_OAuth2Config_Oauth2Config{
						Oauth2Config: plainOAuth2Config,
					},
				},
			}
		}
	case *extauth.AuthConfig_Config_ApiKeyAuth:
		apiKeyConfig, err := translateApiKey(ctx, snap, config.ApiKeyAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_ApiKeyAuth{
			ApiKeyAuth: apiKeyConfig,
		}
	case *extauth.AuthConfig_Config_PluginAuth:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PluginAuth{
			PluginAuth: config.PluginAuth,
		}
	case *extauth.AuthConfig_Config_OpaAuth:
		cfg, err := translateOpaConfig(ctx, snap, config.OpaAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_OpaAuth{OpaAuth: cfg}
	case *extauth.AuthConfig_Config_Ldap:
		if config.Ldap.GroupLookupSettings == nil {
			//use old API if we do not have service account settings
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Ldap{Ldap: config.Ldap}
		} else {
			// translate the config to the new API that includes the service account user and password taken from a secret
			cfg, err := translateLdapConfig(snap, config.Ldap)
			if err != nil {
				return nil, err
			}
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_LdapInternal{
				LdapInternal: cfg,
			}
		}
	case *extauth.AuthConfig_Config_Jwt:
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_Jwt{}
	case *extauth.AuthConfig_Config_HmacAuth:

		cfg, err := translateHmacConfig(ctx, snap, config.HmacAuth)
		if err != nil {
			return nil, err
		}
		extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_HmacAuth{
			HmacAuth: cfg,
		}
	case *extauth.AuthConfig_Config_PassThroughAuth:
		switch protocolConfig := config.PassThroughAuth.GetProtocol().(type) {
		case *extauth.PassThroughAuth_Grpc:
			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PassThroughAuth{
				PassThroughAuth: &extauth.PassThroughAuth{
					Protocol: &extauth.PassThroughAuth_Grpc{
						Grpc: protocolConfig.Grpc,
					},
					Config:           config.PassThroughAuth.Config,
					FailureModeAllow: config.PassThroughAuth.GetFailureModeAllow(),
				},
			}
		case *extauth.PassThroughAuth_Http:
			cfg, err := translateHttpPassthroughConfig(ctx, snap, config.PassThroughAuth)
			if err != nil {
				return nil, err
			}

			extAuthConfig.AuthConfig = &extauth.ExtAuthConfig_Config_PassThroughAuth{
				PassThroughAuth: cfg,
			}
		default:
			return nil, unknownPassThroughProtocolType(config.PassThroughAuth.Protocol)
		}
	default:
		return nil, unknownConfigTypeError
	}
	return extAuthConfig, nil
}

func translateHttpPassthroughConfig(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.PassThroughAuth) (*extauth.PassThroughAuth, error) {
	duplicatedHeaders := []string{}
	for _, header := range config.GetHttp().GetResponse().GetAllowedUpstreamHeaders() {
		if slices.Contains(config.GetHttp().GetResponse().GetAllowedUpstreamHeadersToOverwrite(), header) {
			duplicatedHeaders = append(duplicatedHeaders, header)
		}
	}

	if len(duplicatedHeaders) > 0 {
		return nil, eris.Errorf("The following headers are configured for both append and overwrite in the upstream: %s", strings.Join(duplicatedHeaders, ", "))
	}

	return &extauth.PassThroughAuth{
		Protocol: &extauth.PassThroughAuth_Http{
			Http: config.GetHttp(),
		},
		Config:           config.Config,
		FailureModeAllow: config.GetFailureModeAllow(),
	}, nil
}

func translateOpaConfig(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.OpaAuth) (*extauth.ExtAuthConfig_OpaAuthConfig, error) {

	modules := map[string]string{}
	for _, refs := range config.Modules {
		artifact, err := snap.Artifacts.Find(refs.Namespace, refs.Name)
		if err != nil {
			return nil, err
		}

		for k, v := range artifact.Data {
			if _, ok := modules[k]; !ok {
				modules[k] = v
			} else {
				return nil, duplicateModuleError(k)
			}
		}
	}

	options := opa.Options{
		FastInputConversion: config.GetOptions().GetFastInputConversion(),
	}

	if strings.TrimSpace(config.Query) == "" {
		return nil, emptyQueryError
	}

	// validate that it is a valid opa config
	_, err := opa.NewWithOptions(ctx, config.Query, modules, options)

	return &extauth.ExtAuthConfig_OpaAuthConfig{
		Modules: modules,
		Query:   config.Query,
		Options: config.Options,
	}, err
}

func translateApiKey(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	switch config.GetStorageBackend().(type) {
	case *extauth.ApiKeyAuth_K8SSecretApikeyStorage:
		return translateSecretsApiKey(ctx, snap, config)
	case *extauth.ApiKeyAuth_AerospikeApikeyStorage:
		return translateAerospikeApiKey(ctx, snap, config)
	default:
		return translateSecretsApiKey(ctx, snap, config)
	}
}

func translateAerospikeApiKey(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	if config == nil {
		return nil, errors.New("nil settings")
	}
	storageConfig := config.GetAerospikeApikeyStorage()
	if storageConfig == nil {
		return nil, errors.New("nil storage config")
	}
	// Add metadata if present
	var headersFromKeyMetadata map[string]string
	if len(config.HeadersFromMetadata) > 0 {
		headersFromKeyMetadata = make(map[string]string)
		for k, v := range config.HeadersFromMetadata {
			headersFromKeyMetadata[k] = v.GetName()
		}
		contextutils.LoggerFrom(ctx).Debugw("found headersFromKeyMetadata config",
			zap.Any("headersFromKeyMetadata", headersFromKeyMetadata))
	}
	retConf := &extauth.ExtAuthConfig_ApiKeyAuthConfig{
		StorageBackend: &extauth.ExtAuthConfig_ApiKeyAuthConfig_AerospikeApikeyStorage{
			AerospikeApikeyStorage: storageConfig,
		},
		HeadersFromKeyMetadata: headersFromKeyMetadata,
		HeaderName:             config.HeaderName,
	}
	return retConf, nil
}
func translateSecretsApiKey(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.ApiKeyAuth) (*extauth.ExtAuthConfig_ApiKeyAuthConfig, error) {
	var (
		matchingSecrets []*v1.Secret
		searchErrs      = &multierror.Error{}
		secretErrs      = &multierror.Error{}
	)

	// Find directly referenced secrets
	for _, secretRef := range config.ApiKeySecretRefs {
		secret, err := snap.Secrets.Find(secretRef.Namespace, secretRef.Name)
		if err != nil {
			searchErrs = multierror.Append(searchErrs, err)
			continue
		}
		matchingSecrets = append(matchingSecrets, secret)
	}

	// Find secrets matching provided label selector
	if config.LabelSelector != nil && len(config.LabelSelector) > 0 {
		foundAny := false
		for _, secret := range snap.Secrets {
			selector := labels.Set(config.LabelSelector).AsSelectorPreValidated()
			if selector.Matches(labels.Set(secret.Metadata.Labels)) {
				matchingSecrets = append(matchingSecrets, secret)
				foundAny = true
			}
		}
		if !foundAny {
			// A label may be applied before the underlying secret has been persisted.
			// In this case, we should accept the configuration and just warn the user.
			// Otherwise, this situation blocks configuration from being processed.
			//
			// We do not yet support warnings on AuthConfig CRs, so we log a warning instead
			// Technical Debt: https://github.com/solo-io/solo-projects/issues/2950
			err := noMatchesForGroupError(config.LabelSelector)
			contextutils.LoggerFrom(ctx).Warnf("%v, continuing processing", err)
		}
	}

	if err := searchErrs.ErrorOrNil(); err != nil {
		return nil, err
	}

	var allSecretKeys map[string]string
	if len(config.HeadersFromMetadata) > 0 {
		allSecretKeys = make(map[string]string)
		for k, v := range config.HeadersFromMetadata {
			if v.Required {
				allSecretKeys[k] = v.GetName()
			}
		}
	}
	if len(config.HeadersFromMetadataEntry) > 0 {
		if allSecretKeys == nil {
			allSecretKeys = make(map[string]string)
		}
		for k, v := range config.HeadersFromMetadataEntry {
			if v.Required {
				allSecretKeys[k] = v.GetName()
			}
		}
	}

	var requiredSecretKeys []string
	for _, secretKey := range allSecretKeys {
		requiredSecretKeys = append(requiredSecretKeys, secretKey)
	}

	validApiKeys := make(map[string]*extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata)
	for _, secret := range matchingSecrets {
		apiKeySecret := secret.GetApiKey()
		if apiKeySecret == nil {
			secretErrs = multierror.Append(secretErrs, NonApiKeySecretError(secret))
			continue
		}

		apiKey := apiKeySecret.GetApiKey()
		if apiKey == "" {
			secretErrs = multierror.Append(secretErrs, EmptyApiKeyError(secret))
			continue
		}

		// If there is required metadata, make sure the secret contains it
		secretMetadata := apiKeySecret.GetMetadata()
		for _, requiredKey := range requiredSecretKeys {
			if _, ok := secretMetadata[requiredKey]; !ok {
				secretErrs = multierror.Append(secretErrs, MissingRequiredMetadataError(requiredKey, secret))
				continue
			}
		}

		apiKeyMetadata := &extauth.ExtAuthConfig_ApiKeyAuthConfig_KeyMetadata{
			Username: secret.Metadata.Name,
		}

		if len(secretMetadata) > 0 {
			apiKeyMetadata.Metadata = make(map[string]string)
			for k, v := range secretMetadata {
				apiKeyMetadata.Metadata[k] = v
			}
		}

		validApiKeys[apiKey] = apiKeyMetadata
	}

	apiKeyAuthConfig := &extauth.ExtAuthConfig_ApiKeyAuthConfig{
		HeaderName:   config.HeaderName,
		ValidApiKeys: validApiKeys,
	}

	// Add metadata if present
	if len(config.HeadersFromMetadata) > 0 {
		apiKeyAuthConfig.HeadersFromKeyMetadata = make(map[string]string)
		for k, v := range config.HeadersFromMetadata {
			apiKeyAuthConfig.HeadersFromKeyMetadata[k] = v.GetName()
		}
	}
	if len(config.HeadersFromMetadataEntry) > 0 {
		if apiKeyAuthConfig.HeadersFromKeyMetadata == nil {
			apiKeyAuthConfig.HeadersFromKeyMetadata = make(map[string]string)
		}
		for k, v := range config.HeadersFromMetadataEntry {
			apiKeyAuthConfig.HeadersFromKeyMetadata[k] = v.GetName()
		}
	}

	return apiKeyAuthConfig, secretErrs.ErrorOrNil()
}

// translate deprecated config
func translateOauth(snap *v1snap.ApiSnapshot, config *extauth.OAuth) (*extauth.ExtAuthConfig_OAuthConfig, error) {

	secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
	if err != nil {
		return nil, err
	}

	return &extauth.ExtAuthConfig_OAuthConfig{
		AppUrl:                  config.AppUrl,
		ClientId:                config.ClientId,
		ClientSecret:            secret.GetOauth().GetClientSecret(),
		IssuerUrl:               config.IssuerUrl,
		AuthEndpointQueryParams: config.AuthEndpointQueryParams,
		CallbackPath:            config.CallbackPath,
		Scopes:                  config.Scopes,
	}, nil
}

func translatePlainOAuth2(snap *v1snap.ApiSnapshot, config *extauth.PlainOAuth2) (*extauth.ExtAuthConfig_PlainOAuth2Config, error) {
	secretDisabled := config.GetDisableClientSecret().GetValue()
	clientSecret := ""
	if !secretDisabled {
		secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
		if err != nil {
			return nil, err
		}
		clientSecret = secret.GetOauth().GetClientSecret()
	}
	var session *extauth.UserSession
	// userSession will be set to nil if the cipher config is set. This needs to be applied to any new features as well.
	userSession, err := translateUserSession(snap, config.Session)
	if err != nil {
		return nil, err
	}
	// if the cipher config is set to nil, we want to pass the session as is for backwards compatibility.
	// NOTE: changes to the UserSession must be relected in this conditional as well.
	if config.Session.GetCipherConfig() == nil {
		// session is deprecated, use userSession.
		session, err = userSessionToSession(userSession)
		if err != nil {
			return nil, err
		}
		// if the client sets the cipherConfig, we do not want to set uncrypted cookies for the client. Client
		// should receive 400 errors because they can not authenticate with identity providers.
		userSession = nil
	}

	return &extauth.ExtAuthConfig_PlainOAuth2Config{
		AppUrl:                   config.AppUrl,
		ClientId:                 config.ClientId,
		ClientSecret:             clientSecret,
		AuthEndpointQueryParams:  config.AuthEndpointQueryParams,
		TokenEndpointQueryParams: config.TokenEndpointQueryParams,
		CallbackPath:             config.CallbackPath,
		AfterLogoutUrl:           config.AfterLogoutUrl,
		LogoutPath:               config.LogoutPath,
		Scopes:                   config.Scopes,
		// Session is deprecated, use UserSession. setting session here because upgrades could neglect UserSession from race condition
		Session:            session,
		UserSession:        userSession,
		TokenEndpoint:      config.TokenEndpoint,
		AuthEndpoint:       config.AuthEndpoint,
		RevocationEndpoint: config.RevocationEndpoint,
	}, nil
}

// userSessionToSession will construct a deprecated user session to a session store.
// The userSession is used to construct the redis or cookie session store.
func userSessionToSession(userSession *extauth.ExtAuthConfig_UserSessionConfig) (*extauth.UserSession, error) {
	if userSession == nil {
		return nil, nil
	}
	session := &extauth.UserSession{
		FailOnFetchFailure: userSession.FailOnFetchFailure,
		CookieOptions:      userSession.CookieOptions,
	}
	switch u := userSession.Session.(type) {
	case *extauth.ExtAuthConfig_UserSessionConfig_Cookie:
		session.Session = &extauth.UserSession_Cookie{
			Cookie: u.Cookie,
		}
	case *extauth.ExtAuthConfig_UserSessionConfig_Redis:
		session.Session = &extauth.UserSession_Redis{
			Redis: u.Redis,
		}
	case nil: // no option set which is ok
		break
	default:
		return nil, errors.Errorf("unknown user session type [%T] cannot convert to session", u)
	}
	return session, nil
}

func translateOidcAuthorizationCode(snap *v1snap.ApiSnapshot, config *extauth.OidcAuthorizationCode) (*extauth.ExtAuthConfig_OidcAuthorizationCodeConfig, error) {
	// Configuration used for OIDC client authentication
	clientSecret := ""
	var pkJwtClientAuthenticationConfig *extauth.ExtAuthConfig_OidcAuthorizationCodeConfig_PkJwtClientAuthenticationConfig

	// Translate the Specific OIDC Client Authentication Type
	switch config.GetClientAuthentication().GetClientAuthenticationConfig().(type) {
	// Client Secret Client Authentication
	case *extauth.OidcAuthorizationCode_ClientAuthentication_ClientSecret_:
		if config.GetClientSecretRef() != nil || config.GetDisableClientSecret() != nil {
			return nil, duplicateOidcClientAuthenticationError
		}

		clientSecretConfig := config.GetClientAuthentication().GetClientSecret()
		clientSecretDisabled := clientSecretConfig.GetDisableClientSecret().GetValue()

		// Require the client secret to be defined unless it is specifically disabled
		if !clientSecretDisabled {
			secret, err := snap.Secrets.Find(clientSecretConfig.GetClientSecretRef().GetNamespace(), clientSecretConfig.GetClientSecretRef().GetName())
			if err != nil {
				return nil, errors.New("client secret expected and not found")
			}
			clientSecret = secret.GetOauth().GetClientSecret()
		}

	// Private Key JWT Client Authentication
	case *extauth.OidcAuthorizationCode_ClientAuthentication_PrivateKeyJwt_:
		if config.GetClientSecretRef() != nil || config.GetDisableClientSecret() != nil {
			return nil, duplicateOidcClientAuthenticationError
		}
		signingKey := ""
		// Todo: handle depreaction of client secret
		authenticationConfig := config.GetClientAuthentication().GetPrivateKeyJwt()

		signingKeySecret, err := snap.Secrets.Find(authenticationConfig.GetSigningKeyRef().GetNamespace(), authenticationConfig.GetSigningKeyRef().GetName())
		if err != nil {
			return nil, errors.New("client secret expected and not found")
		}
		signingKey = signingKeySecret.GetOauth().GetClientSecret()

		// Set the default validFor if not set
		validFor := authenticationConfig.GetValidFor()
		if err := validFor.CheckValid(); err != nil {
			validFor = &durationpb.Duration{Seconds: OidcPkJwtClientAuthValidForDefaultSeconds}
		}

		pkJwtClientAuthenticationConfig = &extauth.ExtAuthConfig_OidcAuthorizationCodeConfig_PkJwtClientAuthenticationConfig{
			SigningKey: signingKey,
			ValidFor:   validFor,
		}

	// default is the deprecated Client Secret Client Authentication format or a mistake
	default:
		// Validate Deprecated clientSecret/disableClientSecret fields
		// We're here beacuse we didn't catch any of the expected config types, but make sure there wasn't an empty/malformed codeExchangeType
		if config.GetClientAuthentication() != nil {
			return nil, errors.New("Invalid codeExchangeType, expected clientSecret or privateKeyJwt")
		}

		secretDisabled := config.GetDisableClientSecret().GetValue()
		// Require the client secret to be defined unless it is specifically disabled
		if !secretDisabled {
			secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
			if err != nil {
				return nil, errors.New("client secret expected and not found")
			}
			clientSecret = secret.GetOauth().GetClientSecret()
		}
	}

	sessionIdHeaderName := config.GetSessionIdHeaderName()
	// prefer the session id header name set in redis config, if set.
	switch session := config.GetSession().GetSession().(type) {
	case *extauth.UserSession_Redis:
		if headerName := session.Redis.GetHeaderName(); headerName != "" {
			sessionIdHeaderName = headerName
		}
	}
	var session *extauth.UserSession
	userSessionConfig, err := translateUserSession(snap, config.GetSession())
	if err != nil {
		return nil, err
	}
	// if the cipher config is set to nil, we want to pass the session as is for backwards compatibility.
	// NOTE: changes to the UserSession must be relected in this conditional as well.
	if config.GetSession().GetCipherConfig() == nil {
		// session is deprecated, use userSession.
		// If userSession is nil, there might be an old API that is being used, in this case default to session.
		session, err = userSessionToSession(userSessionConfig)
		if err != nil {
			return nil, err
		}
		// if the client sets the cipherConfig, we do not want to set uncrypted cookies for the client. Client
		// should receive 400 errors because they can not authenticate with identity providers.
		userSessionConfig = nil
	}

	return &extauth.ExtAuthConfig_OidcAuthorizationCodeConfig{
		AppUrl:                   config.AppUrl,
		ClientId:                 config.ClientId,
		ClientSecret:             clientSecret,
		IssuerUrl:                config.IssuerUrl,
		AuthEndpointQueryParams:  config.AuthEndpointQueryParams,
		TokenEndpointQueryParams: config.TokenEndpointQueryParams,
		CallbackPath:             config.CallbackPath,
		AfterLogoutUrl:           config.AfterLogoutUrl,
		SessionIdHeaderName:      sessionIdHeaderName,
		LogoutPath:               config.LogoutPath,
		Scopes:                   config.Scopes,
		// Session is deprecated, use UserSession. setting session here because upgrades could neglect UserSession from race condition
		Session:                         session,
		UserSession:                     userSessionConfig,
		Headers:                         config.Headers,
		DiscoveryOverride:               config.DiscoveryOverride,
		DiscoveryPollInterval:           config.GetDiscoveryPollInterval(),
		JwksCacheRefreshPolicy:          config.GetJwksCacheRefreshPolicy(),
		ParseCallbackPathAsRegex:        config.ParseCallbackPathAsRegex,
		AutoMapFromMetadata:             config.AutoMapFromMetadata,
		EndSessionProperties:            config.EndSessionProperties,
		PkJwtClientAuthenticationConfig: pkJwtClientAuthenticationConfig,
	}, nil

}

// translateUserSession will create the UserSessionConfig used in the deprecated UserSession.
// NOTE: changes to the UserSession must be relected in this conditional as well when building the UserSession.
func translateUserSession(snap *v1snap.ApiSnapshot, config *extauth.UserSession) (*extauth.ExtAuthConfig_UserSessionConfig, error) {
	if config == nil {
		return nil, nil
	}
	userSessionConfig := &extauth.ExtAuthConfig_UserSessionConfig{
		FailOnFetchFailure: config.FailOnFetchFailure,
		CookieOptions:      config.CookieOptions,
	}
	if err := translateUserSessionCipherConfig(snap, config, userSessionConfig); err != nil {
		return nil, err
	}
	if err := translateUserSessionSession(snap, config, userSessionConfig); err != nil {
		return nil, err
	}
	return userSessionConfig, nil
}

func translateUserSessionCipherConfig(snap *v1snap.ApiSnapshot, config *extauth.UserSession, userSessionConfig *extauth.ExtAuthConfig_UserSessionConfig) error {
	const encryptionKeyLength = 32
	cipherKeyRef := config.GetCipherConfig().GetKeyRef()
	if cipherKeyRef != nil {
		cipherSecret, err := snap.Secrets.Find(cipherKeyRef.GetNamespace(), cipherKeyRef.GetName())
		if err != nil {
			return err
		}
		encryptionKey := cipherSecret.GetEncryption().GetKey()
		if len(encryptionKey) != encryptionKeyLength {
			return fmt.Errorf("the encryption key needs to be %d characters in length", encryptionKeyLength)
		}
		userSessionConfig.CipherConfig = &extauth.ExtAuthConfig_UserSessionConfig_CipherConfig{
			Key: encryptionKey,
		}
	}
	return nil
}

func translateUserSessionSession(snap *v1snap.ApiSnapshot, config *extauth.UserSession, userSessionConfig *extauth.ExtAuthConfig_UserSessionConfig) error {
	if config.GetSession() == nil {
		return nil
	}
	switch t := config.GetSession().(type) {
	case *extauth.UserSession_Cookie:
		userSessionConfig.Session = &extauth.ExtAuthConfig_UserSessionConfig_Cookie{
			Cookie: t.Cookie,
		}
	case *extauth.UserSession_Redis:
		userSessionConfig.Session = &extauth.ExtAuthConfig_UserSessionConfig_Redis{
			Redis: t.Redis,
		}
	default:
		return fmt.Errorf("this is not a valid type to for the userSession session: %t", t)
	}
	return nil
}

func translateLdapConfig(snap *v1snap.ApiSnapshot, config *extauth.Ldap) (*extauth.ExtAuthConfig_LdapConfig, error) {
	var translatedGroupLookupSettings *extauth.ExtAuthConfig_LdapServiceAccountConfig
	if config.GroupLookupSettings != nil {
		translatedGroupLookupSettings = &extauth.ExtAuthConfig_LdapServiceAccountConfig{}
		translatedGroupLookupSettings.CheckGroupsWithServiceAccount = config.GetGroupLookupSettings().GetCheckGroupsWithServiceAccount()
		if translatedGroupLookupSettings.CheckGroupsWithServiceAccount && config.GetGroupLookupSettings().CredentialsSecretRef == nil {
			return nil, MissingSecretError()
		}
		if config.GetGroupLookupSettings().CredentialsSecretRef != nil {
			secret, err := snap.Secrets.Find(config.GetGroupLookupSettings().GetCredentialsSecretRef().GetNamespace(),
				config.GetGroupLookupSettings().GetCredentialsSecretRef().GetName())
			if err != nil {
				return nil, err
			}
			if _, ok := secret.GetKind().(*v1.Secret_Credentials); !ok {
				return nil, NonAccountCredentialsSecretError(secret)
			}
			translatedGroupLookupSettings.Username = secret.GetCredentials().GetUsername()
			translatedGroupLookupSettings.Password = secret.GetCredentials().GetPassword()
		}
	}
	translatedConfig := &extauth.ExtAuthConfig_LdapConfig{
		Address:                 config.GetAddress(),
		UserDnTemplate:          config.GetUserDnTemplate(),
		MembershipAttributeName: config.GetMembershipAttributeName(),
		AllowedGroups:           config.GetAllowedGroups(),
		Pool:                    config.GetPool(),
		SearchFilter:            config.GetSearchFilter(),
		DisableGroupChecking:    config.GetDisableGroupChecking(),
		GroupLookupSettings:     translatedGroupLookupSettings,
	}
	return translatedConfig, nil
}
func translateHmacConfig(ctx context.Context, snap *v1snap.ApiSnapshot, config *extauth.HmacAuth) (*extauth.ExtAuthConfig_HmacAuthConfig, error) {
	passwords := make(map[string]string, len(config.GetSecretRefs().GetSecretRefs()))
	secretErrors := &multierror.Error{}
	for _, secretRef := range config.GetSecretRefs().GetSecretRefs() {
		secret, err := snap.Secrets.Find(secretRef.GetNamespace(), secretRef.GetName())
		if err != nil {
			secretErrors = multierror.Append(secretErrors, err)
		} else if _, ok := secret.GetKind().(*v1.Secret_Credentials); !ok {
			secretErrors = multierror.Append(secretErrors, NonAccountCredentialsSecretError(secret))
		} else {
			passwords[secret.GetCredentials().GetUsername()] = secret.GetCredentials().GetPassword()
		}
	}
	if compiledError := secretErrors.ErrorOrNil(); compiledError != nil {
		contextutils.LoggerFrom(ctx).Warnf("Some secrets could not be read. Any valid secrets will be avaiable for authentication. Errors %v", secretErrors)
	}
	if len(passwords) == 0 {
		return nil, noValidUsersError
	}
	translatedConfig := &extauth.ExtAuthConfig_HmacAuthConfig{
		SecretStorage: &extauth.ExtAuthConfig_HmacAuthConfig_SecretList{
			SecretList: &extauth.ExtAuthConfig_InMemorySecretList{SecretList: passwords},
		},
	}
	// When there is more than one implementation, there might be config from some of them to pass through
	switch hmacImpl := config.GetImplementationType().(type) {
	case *extauth.HmacAuth_ParametersInHeaders:
		translatedConfig.ImplementationType = &extauth.ExtAuthConfig_HmacAuthConfig_ParametersInHeaders{
			// this is always empty so it's not technically necessary
			ParametersInHeaders: hmacImpl.ParametersInHeaders,
		}
	default:
		translatedConfig.ImplementationType = &extauth.ExtAuthConfig_HmacAuthConfig_ParametersInHeaders{
			ParametersInHeaders: &extauth.HmacParametersInHeaders{},
		}
	}
	return translatedConfig, nil
}
func translateAccessTokenValidationConfig(snap *v1snap.ApiSnapshot, config *extauth.AccessTokenValidation) (*extauth.ExtAuthConfig_AccessTokenValidationConfig, error) {
	accessTokenValidationConfig := &extauth.ExtAuthConfig_AccessTokenValidationConfig{
		UserinfoUrl:  config.GetUserinfoUrl(),
		CacheTimeout: config.GetCacheTimeout(),
	}

	// ValidationType
	switch validationTypeConfig := config.ValidationType.(type) {
	case *extauth.AccessTokenValidation_IntrospectionUrl:
		accessTokenValidationConfig.ValidationType = &extauth.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionUrl{
			IntrospectionUrl: config.GetIntrospectionUrl(),
		}
	case *extauth.AccessTokenValidation_Introspection:
		introspectionCfg, err := translateAccessTokenValidationIntrospection(snap, validationTypeConfig.Introspection)
		if err != nil {
			return nil, err
		}
		accessTokenValidationConfig.ValidationType = &extauth.ExtAuthConfig_AccessTokenValidationConfig_Introspection{
			Introspection: introspectionCfg,
		}
	case *extauth.AccessTokenValidation_Jwt:
		jwtCfg, err := translateAccessTokenValidationJwt(validationTypeConfig.Jwt)
		if err != nil {
			return nil, err
		}
		accessTokenValidationConfig.ValidationType = &extauth.ExtAuthConfig_AccessTokenValidationConfig_Jwt{
			Jwt: jwtCfg,
		}
	}

	// ScopeValidation
	switch scopeValidationConfig := config.ScopeValidation.(type) {
	case *extauth.AccessTokenValidation_RequiredScopes:
		accessTokenValidationConfig.ScopeValidation = &extauth.ExtAuthConfig_AccessTokenValidationConfig_RequiredScopes{
			RequiredScopes: &extauth.ExtAuthConfig_AccessTokenValidationConfig_ScopeList{
				Scope: scopeValidationConfig.RequiredScopes.GetScope(),
			},
		}
	}

	return accessTokenValidationConfig, nil
}

func translateAccessTokenValidationIntrospection(snap *v1snap.ApiSnapshot, config *extauth.IntrospectionValidation) (*extauth.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionValidation, error) {
	var clientSecret string
	if config.GetClientSecretRef() != nil {
		secret, err := snap.Secrets.Find(config.GetClientSecretRef().GetNamespace(), config.GetClientSecretRef().GetName())
		if err != nil {
			return nil, err
		}
		clientSecret = secret.GetOauth().GetClientSecret()
	}

	return &extauth.ExtAuthConfig_AccessTokenValidationConfig_IntrospectionValidation{
		IntrospectionUrl:    config.GetIntrospectionUrl(),
		ClientId:            config.GetClientId(),
		ClientSecret:        clientSecret,
		UserIdAttributeName: config.GetUserIdAttributeName(),
	}, nil
}

func translateAccessTokenValidationJwt(config *extauth.JwtValidation) (*extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation, error) {
	jwtValidation := &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation{
		Issuer: config.GetIssuer(),
	}

	switch jwksSourceSpecifierConfig := config.JwksSourceSpecifier.(type) {
	case *extauth.JwtValidation_LocalJwks_:
		jwtValidation.JwksSourceSpecifier = &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_LocalJwks_{
			LocalJwks: &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_LocalJwks{
				InlineString: jwksSourceSpecifierConfig.LocalJwks.GetInlineString(),
			},
		}

	case *extauth.JwtValidation_RemoteJwks_:
		jwtValidation.JwksSourceSpecifier = &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_RemoteJwks_{
			RemoteJwks: &extauth.ExtAuthConfig_AccessTokenValidationConfig_JwtValidation_RemoteJwks{
				Url:             jwksSourceSpecifierConfig.RemoteJwks.GetUrl(),
				RefreshInterval: jwksSourceSpecifierConfig.RemoteJwks.GetRefreshInterval(),
			},
		}
	}

	return jwtValidation, nil
}
