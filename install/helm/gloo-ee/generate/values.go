package generate

import glooGen "github.com/solo-io/gloo/install/helm/gloo/generate"

type HelmConfig struct {
	Config
	Global *glooGen.Global `json:"global,omitempty"`
}
type Config struct {
	Settings            *glooGen.Settings `json:"settings,omitempty"`
	LicenseKey          string            `json:"license_key,omitempty"`
	CreateLicenseSecret bool              `json:"create_license_secret"`
	Gloo                *glooGen.Config   `json:"gloo,omitempty"`
	Redis               *Redis            `json:"redis,omitempty"`
	RateLimit           *RateLimit        `json:"rateLimit,omitempty"`
	Observability       *Observability    `json:"observability,omitempty"`
	Rbac                *Rbac             `json:"rbac"`
	Grafana             interface{}       `json:"grafana,omitempty"`
	Prometheus          interface{}       `json:"prometheus,omitempty"`
	Tags                map[string]string `json:"tags,omitempty"`
	ApiServer           *ApiServer        `json:"apiServer,omitempty"`
	DevPortal           *DevPortal        `json:"devPortal,omitempty"`
}

// Common

type Rbac struct {
	Create bool `json:"create"`
}

// Gloo-ee

type GlooEeExtensions struct {
	ExtAuth *ExtAuth `json:"extAuth,omitempty"`
}

type RateLimit struct {
	Enabled    bool                 `json:"enabled" desc:"if true, deploy rate limit service (default true)"`
	Deployment *RateLimitDeployment `json:"deployment,omitempty"`
	Service    *RateLimitService    `json:"service,omitempty"`
}

type DynamoDb struct {
	Region             string `json:"region" desc:"aws region to run DynamoDB requests in"`
	SecretName         string `json:"secretName,omitempty" desc:"name of the aws secret in gloo's installation namespace that has aws creds (if provided, uses DynamoDB to back rate-limiting service instead of Redis)"`
	RateLimitTableName string `json:"tableName" desc:"DynamoDB table name used to back rate limit service (default rate-limits)"`
	ConsistentReads    bool   `json:"consistentReads" desc:"if true, reads from DynamoDB will be strongly consistent (default false)"`
	BatchSize          uint8  `json:"batchSize" desc:"batch size for get requests to DynamoDB (max 100, default 100)"`
}

type RateLimitDeployment struct {
	RedisUrl    string         `json:"redisUrl"`
	GlooAddress string         `json:"glooAddress"`
	DynamoDb    DynamoDb       `json:"dynamodb"`
	Image       *glooGen.Image `json:"image,omitempty"`
	Stats       *glooGen.Stats `json:"stats"`
	*glooGen.DeploymentSpec
}

type RateLimitService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type Redis struct {
	Deployment *RedisDeployment `json:"deployment,omitempty"`
	Service    RedisService     `json:"service,omitempty"`
}

type RedisDeployment struct {
	Image      *glooGen.Image `json:"image,omitempty"`
	StaticPort uint           `json:"staticPort"`
	*glooGen.DeploymentSpec
}

type RedisService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type Observability struct {
	Enabled                   bool                     `json:"enabled,omitempty" desc:"if true, deploy observability service (default true)"`
	Deployment                *ObservabilityDeployment `json:"deployment,omitempty"`
	CustomGrafana             *CustomGrafana           `json:"customGrafana" desc:"Configure a custom grafana deployment to work with Gloo observability, rather than the default Gloo grafana"`
	UpstreamDashboardTemplate string                   `json:"upstreamDashboardTemplate" desc:"Provide a custom dashboard template to use when generating per-upstream dashboards. The only variables available for use in this template are: {{.Uid}} and {{.EnvoyClusterName}}. Recommended to use Helm's --set-file to provide this value."`
}

type ObservabilityDeployment struct {
	Image *glooGen.Image `json:"image,omitempty"`
	Stats *glooGen.Stats `json:"stats"`
	*glooGen.DeploymentSpec
}

type CustomGrafana struct {
	Enabled  bool   `json:"enabled",omitempty,desc:"Set to true to indicate that the observability pod should talk to a custom grafana instance"`
	Username string `json:"username",omitempty,desc:"Set this and the 'password' field to authenticate to the custom grafana instance using basic auth"`
	Password string `json:"password",omitempty,desc:"Set this and the 'username' field to authenticate to the custom grafana instance using basic auth"`
	ApiKey   string `json:"apiKey",omitempty,desc:"Authenticate to the custom grafana instance using this api key"`
	Url      string `json:"url",omitempty,desc:"The URL for the custom grafana instance"`
	CaBundle string `json:"caBundle",omitempty,desc:"The Certificate Authority used to verify the server certificates.'"`
}

type ExtAuth struct {
	Enabled              bool                      `json:"enabled,omitempty" desc:"if true, deploy ExtAuth service (default true)"`
	UserIdHeader         string                    `json:"userIdHeader,omitempty"`
	Deployment           *ExtAuthDeployment        `json:"deployment,omitempty"`
	Service              *ExtAuthService           `json:"service,omitempty"`
	SigningKey           *ExtAuthSigningKey        `json:"signingKey,omitempty"`
	Plugins              map[string]*ExtAuthPlugin `json:"plugins,omitempty"`
	EnvoySidecar         bool                      `json:"envoySidecar" desc:"if true, deploy ExtAuth as a sidecar with envoy (defaults to false)"`
	StandaloneDeployment bool                      `json:"standaloneDeployment" desc:"if true, create a standalone ExtAuth deployment (defaults to true)"`
	ServiceName          string                    `json:"serviceName,omitempty"`
	RequestTimeout       string                    `json:"requestTimeout,omitempty" desc:"Timeout for the ext auth service to respond (defaults to 200ms)"`
}

type ExtAuthDeployment struct {
	Name        string         `json:"name"`
	GlooAddress string         `json:"glooAddress,omitempty"`
	Port        uint           `json:"port"`
	Image       *glooGen.Image `json:"image,omitempty"`
	Stats       *glooGen.Stats `json:"stats"`
	*glooGen.DeploymentSpec
}

type ExtAuthService struct {
	Port uint   `json:"port"`
	Name string `json:"name"`
}

type ExtAuthSigningKey struct {
	Name       string `json:"name"`
	SigningKey string `json:"signing-key"`
}

type ExtAuthPlugin struct {
	Image *glooGen.Image `json:"image,omitempty"`
}

type ApiServer struct {
	Enable bool `json:"enable,omitempty" desc:"If set, will deploy a read-only UI for Gloo"`
	// used for gating config (like license secret) that are only relevant to the enterprise UI
	Enterprise bool                 `json:"enterprise,omitempty"`
	Deployment *ApiServerDeployment `json:"deployment,omitempty"`
	Service    *ApiServerService    `json:"service,omitempty"`
	ConfigMap  *ApiServerConfigMap  `json:"configMap,omitempty"`
	EnableBeta bool                 `json:"enableBeta,omitempty"`
}

type ApiServerDeployment struct {
	Server *ApiServerServerDeployment `json:"server,omitempty"`
	Ui     *ApiServerUiDeployment     `json:"ui,omitempty"`
	Envoy  *ApiServerEnvoyDeployment  `json:"envoy,omitempty"`
	Stats  *glooGen.Stats             `json:"stats"`
	*glooGen.DeploymentSpec
}

type ApiServerServerDeployment struct {
	GrpcPort uint           `json:"grpcPort"`
	OAuth    *OAuth         `json:"oauth,omitempty"`
	Image    *glooGen.Image `json:"image"`
	*glooGen.DeploymentSpec
}

type ApiServerEnvoyDeployment struct {
	Image *glooGen.Image `json:"image"`
	*glooGen.DeploymentSpec
}

type ApiServerUiDeployment struct {
	StaticPort uint           `json:"staticPort"`
	Image      *glooGen.Image `json:"image,omitempty"`
	*glooGen.DeploymentSpec
}

type ApiServerService struct {
	Name string `json:"name"`
}

type ApiServerConfigMap struct {
	Name string `json:"name"`
}

type OAuth struct {
	Server string `json:"server"`
	Client string `json:"client"`
}

type DevPortal struct {
	Enabled bool `json:"enabled" desc:"if true, deploy the developer portal. Defaults to false"`
}
