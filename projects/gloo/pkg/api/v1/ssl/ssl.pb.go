// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/v1/ssl/ssl.proto

package ssl

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/solo-io/protoc-gen-ext/extproto"
	core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SslConfig_OcspStaplePolicy int32

const (
	// OCSP responses are optional. If none is provided, or the provided response is expired, the associated certificate will be used without the OCSP response.
	SslConfig_LENIENT_STAPLING SslConfig_OcspStaplePolicy = 0
	// OCSP responses are optional. If none is provided, the associated certificate will be used without the OCSP response.
	// If a response is present, but expired, the certificate will not be used for connections.
	// If no suitable certificate is found, the connection is rejected.
	SslConfig_STRICT_STAPLING SslConfig_OcspStaplePolicy = 1
	// OCSP responses are required. If no `ocsp_staple` is set on a certificate, configuration will fail.
	// If a response is expired, the associated certificate will not be used.
	// If no suitable certificate is found, the connection is rejected.
	SslConfig_MUST_STAPLE SslConfig_OcspStaplePolicy = 2
)

// Enum value maps for SslConfig_OcspStaplePolicy.
var (
	SslConfig_OcspStaplePolicy_name = map[int32]string{
		0: "LENIENT_STAPLING",
		1: "STRICT_STAPLING",
		2: "MUST_STAPLE",
	}
	SslConfig_OcspStaplePolicy_value = map[string]int32{
		"LENIENT_STAPLING": 0,
		"STRICT_STAPLING":  1,
		"MUST_STAPLE":      2,
	}
)

func (x SslConfig_OcspStaplePolicy) Enum() *SslConfig_OcspStaplePolicy {
	p := new(SslConfig_OcspStaplePolicy)
	*p = x
	return p
}

func (x SslConfig_OcspStaplePolicy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SslConfig_OcspStaplePolicy) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_enumTypes[0].Descriptor()
}

func (SslConfig_OcspStaplePolicy) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_enumTypes[0]
}

func (x SslConfig_OcspStaplePolicy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SslConfig_OcspStaplePolicy.Descriptor instead.
func (SslConfig_OcspStaplePolicy) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{0, 0}
}

type SslParameters_ProtocolVersion int32

const (
	// Envoy will choose the optimal TLS version.
	SslParameters_TLS_AUTO SslParameters_ProtocolVersion = 0
	// TLS 1.0
	SslParameters_TLSv1_0 SslParameters_ProtocolVersion = 1
	// TLS 1.1
	SslParameters_TLSv1_1 SslParameters_ProtocolVersion = 2
	// TLS 1.2
	SslParameters_TLSv1_2 SslParameters_ProtocolVersion = 3
	// TLS 1.3
	SslParameters_TLSv1_3 SslParameters_ProtocolVersion = 4
)

// Enum value maps for SslParameters_ProtocolVersion.
var (
	SslParameters_ProtocolVersion_name = map[int32]string{
		0: "TLS_AUTO",
		1: "TLSv1_0",
		2: "TLSv1_1",
		3: "TLSv1_2",
		4: "TLSv1_3",
	}
	SslParameters_ProtocolVersion_value = map[string]int32{
		"TLS_AUTO": 0,
		"TLSv1_0":  1,
		"TLSv1_1":  2,
		"TLSv1_2":  3,
		"TLSv1_3":  4,
	}
)

func (x SslParameters_ProtocolVersion) Enum() *SslParameters_ProtocolVersion {
	p := new(SslParameters_ProtocolVersion)
	*p = x
	return p
}

func (x SslParameters_ProtocolVersion) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SslParameters_ProtocolVersion) Descriptor() protoreflect.EnumDescriptor {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_enumTypes[1].Descriptor()
}

func (SslParameters_ProtocolVersion) Type() protoreflect.EnumType {
	return &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_enumTypes[1]
}

func (x SslParameters_ProtocolVersion) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SslParameters_ProtocolVersion.Descriptor instead.
func (SslParameters_ProtocolVersion) EnumDescriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{5, 0}
}

// SslConfig contains the options necessary to configure a virtual host or listener to use TLS termination
type SslConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to SslSecrets:
	//
	//	*SslConfig_SecretRef
	//	*SslConfig_SslFiles
	//	*SslConfig_Sds
	SslSecrets isSslConfig_SslSecrets `protobuf_oneof:"ssl_secrets"`
	// optional. the SNI domains that should be considered for TLS connections
	SniDomains []string `protobuf:"bytes,3,rep,name=sni_domains,json=sniDomains,proto3" json:"sni_domains,omitempty"`
	// Verify that the Subject Alternative Name in the peer certificate is one of the specified values.
	// note that a root_ca must be provided if this option is used.
	VerifySubjectAltName []string       `protobuf:"bytes,5,rep,name=verify_subject_alt_name,json=verifySubjectAltName,proto3" json:"verify_subject_alt_name,omitempty"`
	Parameters           *SslParameters `protobuf:"bytes,6,opt,name=parameters,proto3" json:"parameters,omitempty"`
	// Set Application Level Protocol Negotiation
	// If empty, defaults to ["h2", "http/1.1"].
	// As an advanced option you may use ["allow_empty"] to avoid defaults and set alpn to have no alpn set (ie pass empty slice).
	AlpnProtocols []string `protobuf:"bytes,7,rep,name=alpn_protocols,json=alpnProtocols,proto3" json:"alpn_protocols,omitempty"`
	// If the SSL config has the ca.crt (root CA) provided, Gloo uses it to perform mTLS by default.
	// Set oneWayTls to true to disable mTLS in favor of server-only TLS (one-way TLS), even if Gloo has the root CA.
	// If unset, defaults to false.
	OneWayTls *wrapperspb.BoolValue `protobuf:"bytes,8,opt,name=one_way_tls,json=oneWayTls,proto3" json:"one_way_tls,omitempty"`
	// If set to true, the TLS session resumption will be deactivated, note that it deactivates only the tickets based tls session resumption (not the cache).
	DisableTlsSessionResumption *wrapperspb.BoolValue `protobuf:"bytes,9,opt,name=disable_tls_session_resumption,json=disableTlsSessionResumption,proto3" json:"disable_tls_session_resumption,omitempty"`
	// If present and nonzero, the amount of time to allow incoming connections to complete any
	// transport socket negotiations. If this expires before the transport reports connection
	// establishment, the connection is summarily closed.
	TransportSocketConnectTimeout *durationpb.Duration `protobuf:"bytes,10,opt,name=transport_socket_connect_timeout,json=transportSocketConnectTimeout,proto3" json:"transport_socket_connect_timeout,omitempty"`
	// The OCSP staple policy to use for this listener.
	// Defaults to `LENIENT_STAPLING`.
	// https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/transport_sockets/tls/v3/tls.proto#enum-extensions-transport-sockets-tls-v3-downstreamtlscontext-ocspstaplepolicy
	OcspStaplePolicy SslConfig_OcspStaplePolicy `protobuf:"varint,11,opt,name=ocsp_staple_policy,json=ocspStaplePolicy,proto3,enum=gloo.solo.io.SslConfig_OcspStaplePolicy" json:"ocsp_staple_policy,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *SslConfig) Reset() {
	*x = SslConfig{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SslConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SslConfig) ProtoMessage() {}

func (x *SslConfig) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SslConfig.ProtoReflect.Descriptor instead.
func (*SslConfig) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{0}
}

func (x *SslConfig) GetSslSecrets() isSslConfig_SslSecrets {
	if x != nil {
		return x.SslSecrets
	}
	return nil
}

func (x *SslConfig) GetSecretRef() *core.ResourceRef {
	if x != nil {
		if x, ok := x.SslSecrets.(*SslConfig_SecretRef); ok {
			return x.SecretRef
		}
	}
	return nil
}

func (x *SslConfig) GetSslFiles() *SSLFiles {
	if x != nil {
		if x, ok := x.SslSecrets.(*SslConfig_SslFiles); ok {
			return x.SslFiles
		}
	}
	return nil
}

func (x *SslConfig) GetSds() *SDSConfig {
	if x != nil {
		if x, ok := x.SslSecrets.(*SslConfig_Sds); ok {
			return x.Sds
		}
	}
	return nil
}

func (x *SslConfig) GetSniDomains() []string {
	if x != nil {
		return x.SniDomains
	}
	return nil
}

func (x *SslConfig) GetVerifySubjectAltName() []string {
	if x != nil {
		return x.VerifySubjectAltName
	}
	return nil
}

func (x *SslConfig) GetParameters() *SslParameters {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *SslConfig) GetAlpnProtocols() []string {
	if x != nil {
		return x.AlpnProtocols
	}
	return nil
}

func (x *SslConfig) GetOneWayTls() *wrapperspb.BoolValue {
	if x != nil {
		return x.OneWayTls
	}
	return nil
}

func (x *SslConfig) GetDisableTlsSessionResumption() *wrapperspb.BoolValue {
	if x != nil {
		return x.DisableTlsSessionResumption
	}
	return nil
}

func (x *SslConfig) GetTransportSocketConnectTimeout() *durationpb.Duration {
	if x != nil {
		return x.TransportSocketConnectTimeout
	}
	return nil
}

func (x *SslConfig) GetOcspStaplePolicy() SslConfig_OcspStaplePolicy {
	if x != nil {
		return x.OcspStaplePolicy
	}
	return SslConfig_LENIENT_STAPLING
}

type isSslConfig_SslSecrets interface {
	isSslConfig_SslSecrets()
}

type SslConfig_SecretRef struct {
	// SecretRef contains the secret ref to a gloo tls secret or a kubernetes tls secret.
	// gloo tls secret can contain a root ca as well if verification is needed.
	SecretRef *core.ResourceRef `protobuf:"bytes,1,opt,name=secret_ref,json=secretRef,proto3,oneof"`
}

type SslConfig_SslFiles struct {
	// SSLFiles reference paths to certificates which are local to the proxy
	SslFiles *SSLFiles `protobuf:"bytes,2,opt,name=ssl_files,json=sslFiles,proto3,oneof"`
}

type SslConfig_Sds struct {
	// Use secret discovery service.
	Sds *SDSConfig `protobuf:"bytes,4,opt,name=sds,proto3,oneof"`
}

func (*SslConfig_SecretRef) isSslConfig_SslSecrets() {}

func (*SslConfig_SslFiles) isSslConfig_SslSecrets() {}

func (*SslConfig_Sds) isSslConfig_SslSecrets() {}

// SSLFiles reference paths to certificates which can be read by the proxy off of its local filesystem
type SSLFiles struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	TlsCert string                 `protobuf:"bytes,1,opt,name=tls_cert,json=tlsCert,proto3" json:"tls_cert,omitempty"`
	TlsKey  string                 `protobuf:"bytes,2,opt,name=tls_key,json=tlsKey,proto3" json:"tls_key,omitempty"`
	// for client cert validation. optional
	RootCa string `protobuf:"bytes,3,opt,name=root_ca,json=rootCa,proto3" json:"root_ca,omitempty"`
	// stapled ocsp response. optional
	// should be der-encoded
	OcspStaple    string `protobuf:"bytes,4,opt,name=ocsp_staple,json=ocspStaple,proto3" json:"ocsp_staple,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SSLFiles) Reset() {
	*x = SSLFiles{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SSLFiles) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SSLFiles) ProtoMessage() {}

func (x *SSLFiles) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SSLFiles.ProtoReflect.Descriptor instead.
func (*SSLFiles) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{1}
}

func (x *SSLFiles) GetTlsCert() string {
	if x != nil {
		return x.TlsCert
	}
	return ""
}

func (x *SSLFiles) GetTlsKey() string {
	if x != nil {
		return x.TlsKey
	}
	return ""
}

func (x *SSLFiles) GetRootCa() string {
	if x != nil {
		return x.RootCa
	}
	return ""
}

func (x *SSLFiles) GetOcspStaple() string {
	if x != nil {
		return x.OcspStaple
	}
	return ""
}

// SslConfig contains the options necessary to configure an upstream to use TLS origination
type UpstreamSslConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to SslSecrets:
	//
	//	*UpstreamSslConfig_SecretRef
	//	*UpstreamSslConfig_SslFiles
	//	*UpstreamSslConfig_Sds
	SslSecrets isUpstreamSslConfig_SslSecrets `protobuf_oneof:"ssl_secrets"`
	// optional. the SNI domains that should be considered for TLS connections
	Sni string `protobuf:"bytes,3,opt,name=sni,proto3" json:"sni,omitempty"`
	// Verify that the Subject Alternative Name in the peer certificate is one of the specified values.
	// note that a root_ca must be provided if this option is used.
	VerifySubjectAltName []string       `protobuf:"bytes,5,rep,name=verify_subject_alt_name,json=verifySubjectAltName,proto3" json:"verify_subject_alt_name,omitempty"`
	Parameters           *SslParameters `protobuf:"bytes,7,opt,name=parameters,proto3" json:"parameters,omitempty"`
	// Set Application Level Protocol Negotiation.
	// If empty, it is not set.
	AlpnProtocols []string `protobuf:"bytes,8,rep,name=alpn_protocols,json=alpnProtocols,proto3" json:"alpn_protocols,omitempty"`
	// Allow Tls renegotiation, the default value is false.
	// TLS renegotiation is considered insecure and shouldn’t be used unless absolutely necessary.
	AllowRenegotiation *wrapperspb.BoolValue `protobuf:"bytes,10,opt,name=allow_renegotiation,json=allowRenegotiation,proto3" json:"allow_renegotiation,omitempty"`
	// If the SSL config has the ca.crt (root CA) provided, Gloo uses it to perform mTLS by default.
	// Set oneWayTls to true to disable mTLS in favor of server-only TLS (one-way TLS), even if Gloo has the root CA.
	// This flag does nothing if SDS is configured.
	// If unset, defaults to false.
	OneWayTls     *wrapperspb.BoolValue `protobuf:"bytes,11,opt,name=one_way_tls,json=oneWayTls,proto3" json:"one_way_tls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpstreamSslConfig) Reset() {
	*x = UpstreamSslConfig{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpstreamSslConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpstreamSslConfig) ProtoMessage() {}

func (x *UpstreamSslConfig) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpstreamSslConfig.ProtoReflect.Descriptor instead.
func (*UpstreamSslConfig) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{2}
}

func (x *UpstreamSslConfig) GetSslSecrets() isUpstreamSslConfig_SslSecrets {
	if x != nil {
		return x.SslSecrets
	}
	return nil
}

func (x *UpstreamSslConfig) GetSecretRef() *core.ResourceRef {
	if x != nil {
		if x, ok := x.SslSecrets.(*UpstreamSslConfig_SecretRef); ok {
			return x.SecretRef
		}
	}
	return nil
}

func (x *UpstreamSslConfig) GetSslFiles() *SSLFiles {
	if x != nil {
		if x, ok := x.SslSecrets.(*UpstreamSslConfig_SslFiles); ok {
			return x.SslFiles
		}
	}
	return nil
}

func (x *UpstreamSslConfig) GetSds() *SDSConfig {
	if x != nil {
		if x, ok := x.SslSecrets.(*UpstreamSslConfig_Sds); ok {
			return x.Sds
		}
	}
	return nil
}

func (x *UpstreamSslConfig) GetSni() string {
	if x != nil {
		return x.Sni
	}
	return ""
}

func (x *UpstreamSslConfig) GetVerifySubjectAltName() []string {
	if x != nil {
		return x.VerifySubjectAltName
	}
	return nil
}

func (x *UpstreamSslConfig) GetParameters() *SslParameters {
	if x != nil {
		return x.Parameters
	}
	return nil
}

func (x *UpstreamSslConfig) GetAlpnProtocols() []string {
	if x != nil {
		return x.AlpnProtocols
	}
	return nil
}

func (x *UpstreamSslConfig) GetAllowRenegotiation() *wrapperspb.BoolValue {
	if x != nil {
		return x.AllowRenegotiation
	}
	return nil
}

func (x *UpstreamSslConfig) GetOneWayTls() *wrapperspb.BoolValue {
	if x != nil {
		return x.OneWayTls
	}
	return nil
}

type isUpstreamSslConfig_SslSecrets interface {
	isUpstreamSslConfig_SslSecrets()
}

type UpstreamSslConfig_SecretRef struct {
	// SecretRef contains the secret ref to a gloo tls secret or a kubernetes tls secret.
	// gloo tls secret can contain a root ca as well if verification is needed.
	SecretRef *core.ResourceRef `protobuf:"bytes,1,opt,name=secret_ref,json=secretRef,proto3,oneof"`
}

type UpstreamSslConfig_SslFiles struct {
	// SSLFiles reference paths to certificates which are local to the proxy
	SslFiles *SSLFiles `protobuf:"bytes,2,opt,name=ssl_files,json=sslFiles,proto3,oneof"`
}

type UpstreamSslConfig_Sds struct {
	// Use secret discovery service.
	Sds *SDSConfig `protobuf:"bytes,4,opt,name=sds,proto3,oneof"`
}

func (*UpstreamSslConfig_SecretRef) isUpstreamSslConfig_SslSecrets() {}

func (*UpstreamSslConfig_SslFiles) isUpstreamSslConfig_SslSecrets() {}

func (*UpstreamSslConfig_Sds) isUpstreamSslConfig_SslSecrets() {}

type SDSConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Target uri for the sds channel. currently only a unix domain socket is supported.
	TargetUri string `protobuf:"bytes,1,opt,name=target_uri,json=targetUri,proto3" json:"target_uri,omitempty"`
	// Types that are valid to be assigned to SdsBuilder:
	//
	//	*SDSConfig_CallCredentials
	//	*SDSConfig_ClusterName
	SdsBuilder isSDSConfig_SdsBuilder `protobuf_oneof:"sds_builder"`
	// The name of the secret containing the certificate
	CertificatesSecretName string `protobuf:"bytes,3,opt,name=certificates_secret_name,json=certificatesSecretName,proto3" json:"certificates_secret_name,omitempty"`
	// The name of secret containing the validation context (i.e. root ca)
	ValidationContextName string `protobuf:"bytes,4,opt,name=validation_context_name,json=validationContextName,proto3" json:"validation_context_name,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *SDSConfig) Reset() {
	*x = SDSConfig{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SDSConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SDSConfig) ProtoMessage() {}

func (x *SDSConfig) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SDSConfig.ProtoReflect.Descriptor instead.
func (*SDSConfig) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{3}
}

func (x *SDSConfig) GetTargetUri() string {
	if x != nil {
		return x.TargetUri
	}
	return ""
}

func (x *SDSConfig) GetSdsBuilder() isSDSConfig_SdsBuilder {
	if x != nil {
		return x.SdsBuilder
	}
	return nil
}

func (x *SDSConfig) GetCallCredentials() *CallCredentials {
	if x != nil {
		if x, ok := x.SdsBuilder.(*SDSConfig_CallCredentials); ok {
			return x.CallCredentials
		}
	}
	return nil
}

func (x *SDSConfig) GetClusterName() string {
	if x != nil {
		if x, ok := x.SdsBuilder.(*SDSConfig_ClusterName); ok {
			return x.ClusterName
		}
	}
	return ""
}

func (x *SDSConfig) GetCertificatesSecretName() string {
	if x != nil {
		return x.CertificatesSecretName
	}
	return ""
}

func (x *SDSConfig) GetValidationContextName() string {
	if x != nil {
		return x.ValidationContextName
	}
	return ""
}

type isSDSConfig_SdsBuilder interface {
	isSDSConfig_SdsBuilder()
}

type SDSConfig_CallCredentials struct {
	// Call credentials.
	CallCredentials *CallCredentials `protobuf:"bytes,2,opt,name=call_credentials,json=callCredentials,proto3,oneof"`
}

type SDSConfig_ClusterName struct {
	// The name of the sds cluster in envoy
	ClusterName string `protobuf:"bytes,5,opt,name=cluster_name,json=clusterName,proto3,oneof"`
}

func (*SDSConfig_CallCredentials) isSDSConfig_SdsBuilder() {}

func (*SDSConfig_ClusterName) isSDSConfig_SdsBuilder() {}

type CallCredentials struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Call credentials are coming from a file,
	FileCredentialSource *CallCredentials_FileCredentialSource `protobuf:"bytes,1,opt,name=file_credential_source,json=fileCredentialSource,proto3" json:"file_credential_source,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *CallCredentials) Reset() {
	*x = CallCredentials{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CallCredentials) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallCredentials) ProtoMessage() {}

func (x *CallCredentials) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallCredentials.ProtoReflect.Descriptor instead.
func (*CallCredentials) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{4}
}

func (x *CallCredentials) GetFileCredentialSource() *CallCredentials_FileCredentialSource {
	if x != nil {
		return x.FileCredentialSource
	}
	return nil
}

// General TLS parameters. See the [envoy docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/transport_sockets/tls/v3/common.proto#extensions-transport-sockets-tls-v3-tlsparameters)
// for more information on the meaning of these values.
type SslParameters struct {
	state                  protoimpl.MessageState        `protogen:"open.v1"`
	MinimumProtocolVersion SslParameters_ProtocolVersion `protobuf:"varint,1,opt,name=minimum_protocol_version,json=minimumProtocolVersion,proto3,enum=gloo.solo.io.SslParameters_ProtocolVersion" json:"minimum_protocol_version,omitempty"`
	MaximumProtocolVersion SslParameters_ProtocolVersion `protobuf:"varint,2,opt,name=maximum_protocol_version,json=maximumProtocolVersion,proto3,enum=gloo.solo.io.SslParameters_ProtocolVersion" json:"maximum_protocol_version,omitempty"`
	CipherSuites           []string                      `protobuf:"bytes,3,rep,name=cipher_suites,json=cipherSuites,proto3" json:"cipher_suites,omitempty"`
	EcdhCurves             []string                      `protobuf:"bytes,4,rep,name=ecdh_curves,json=ecdhCurves,proto3" json:"ecdh_curves,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *SslParameters) Reset() {
	*x = SslParameters{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SslParameters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SslParameters) ProtoMessage() {}

func (x *SslParameters) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SslParameters.ProtoReflect.Descriptor instead.
func (*SslParameters) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{5}
}

func (x *SslParameters) GetMinimumProtocolVersion() SslParameters_ProtocolVersion {
	if x != nil {
		return x.MinimumProtocolVersion
	}
	return SslParameters_TLS_AUTO
}

func (x *SslParameters) GetMaximumProtocolVersion() SslParameters_ProtocolVersion {
	if x != nil {
		return x.MaximumProtocolVersion
	}
	return SslParameters_TLS_AUTO
}

func (x *SslParameters) GetCipherSuites() []string {
	if x != nil {
		return x.CipherSuites
	}
	return nil
}

func (x *SslParameters) GetEcdhCurves() []string {
	if x != nil {
		return x.EcdhCurves
	}
	return nil
}

type CallCredentials_FileCredentialSource struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// File containing auth token.
	TokenFileName string `protobuf:"bytes,1,opt,name=token_file_name,json=tokenFileName,proto3" json:"token_file_name,omitempty"`
	// Header to carry the token.
	Header        string `protobuf:"bytes,2,opt,name=header,proto3" json:"header,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CallCredentials_FileCredentialSource) Reset() {
	*x = CallCredentials_FileCredentialSource{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CallCredentials_FileCredentialSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CallCredentials_FileCredentialSource) ProtoMessage() {}

func (x *CallCredentials_FileCredentialSource) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CallCredentials_FileCredentialSource.ProtoReflect.Descriptor instead.
func (*CallCredentials_FileCredentialSource) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP(), []int{4, 0}
}

func (x *CallCredentials_FileCredentialSource) GetTokenFileName() string {
	if x != nil {
		return x.TokenFileName
	}
	return ""
}

func (x *CallCredentials_FileCredentialSource) GetHeader() string {
	if x != nil {
		return x.Header
	}
	return ""
}

var File_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDesc = "" +
	"\n" +
	":github.com/solo-io/gloo/projects/gloo/api/v1/ssl/ssl.proto\x12\fgloo.solo.io\x1a\x1egoogle/protobuf/duration.proto\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x12extproto/ext.proto\x1a,github.com/solo-io/solo-kit/api/v1/ref.proto\"\x9f\x06\n" +
	"\tSslConfig\x12:\n" +
	"\n" +
	"secret_ref\x18\x01 \x01(\v2\x19.core.solo.io.ResourceRefH\x00R\tsecretRef\x125\n" +
	"\tssl_files\x18\x02 \x01(\v2\x16.gloo.solo.io.SSLFilesH\x00R\bsslFiles\x12+\n" +
	"\x03sds\x18\x04 \x01(\v2\x17.gloo.solo.io.SDSConfigH\x00R\x03sds\x12\x1f\n" +
	"\vsni_domains\x18\x03 \x03(\tR\n" +
	"sniDomains\x125\n" +
	"\x17verify_subject_alt_name\x18\x05 \x03(\tR\x14verifySubjectAltName\x12;\n" +
	"\n" +
	"parameters\x18\x06 \x01(\v2\x1b.gloo.solo.io.SslParametersR\n" +
	"parameters\x12%\n" +
	"\x0ealpn_protocols\x18\a \x03(\tR\ralpnProtocols\x12:\n" +
	"\vone_way_tls\x18\b \x01(\v2\x1a.google.protobuf.BoolValueR\toneWayTls\x12_\n" +
	"\x1edisable_tls_session_resumption\x18\t \x01(\v2\x1a.google.protobuf.BoolValueR\x1bdisableTlsSessionResumption\x12b\n" +
	" transport_socket_connect_timeout\x18\n" +
	" \x01(\v2\x19.google.protobuf.DurationR\x1dtransportSocketConnectTimeout\x12V\n" +
	"\x12ocsp_staple_policy\x18\v \x01(\x0e2(.gloo.solo.io.SslConfig.OcspStaplePolicyR\x10ocspStaplePolicy\"N\n" +
	"\x10OcspStaplePolicy\x12\x14\n" +
	"\x10LENIENT_STAPLING\x10\x00\x12\x13\n" +
	"\x0fSTRICT_STAPLING\x10\x01\x12\x0f\n" +
	"\vMUST_STAPLE\x10\x02B\r\n" +
	"\vssl_secrets\"x\n" +
	"\bSSLFiles\x12\x19\n" +
	"\btls_cert\x18\x01 \x01(\tR\atlsCert\x12\x17\n" +
	"\atls_key\x18\x02 \x01(\tR\x06tlsKey\x12\x17\n" +
	"\aroot_ca\x18\x03 \x01(\tR\x06rootCa\x12\x1f\n" +
	"\vocsp_staple\x18\x04 \x01(\tR\n" +
	"ocspStaple\"\xf8\x03\n" +
	"\x11UpstreamSslConfig\x12:\n" +
	"\n" +
	"secret_ref\x18\x01 \x01(\v2\x19.core.solo.io.ResourceRefH\x00R\tsecretRef\x125\n" +
	"\tssl_files\x18\x02 \x01(\v2\x16.gloo.solo.io.SSLFilesH\x00R\bsslFiles\x12+\n" +
	"\x03sds\x18\x04 \x01(\v2\x17.gloo.solo.io.SDSConfigH\x00R\x03sds\x12\x10\n" +
	"\x03sni\x18\x03 \x01(\tR\x03sni\x125\n" +
	"\x17verify_subject_alt_name\x18\x05 \x03(\tR\x14verifySubjectAltName\x12;\n" +
	"\n" +
	"parameters\x18\a \x01(\v2\x1b.gloo.solo.io.SslParametersR\n" +
	"parameters\x12%\n" +
	"\x0ealpn_protocols\x18\b \x03(\tR\ralpnProtocols\x12K\n" +
	"\x13allow_renegotiation\x18\n" +
	" \x01(\v2\x1a.google.protobuf.BoolValueR\x12allowRenegotiation\x12:\n" +
	"\vone_way_tls\x18\v \x01(\v2\x1a.google.protobuf.BoolValueR\toneWayTlsB\r\n" +
	"\vssl_secrets\"\x9c\x02\n" +
	"\tSDSConfig\x12\x1d\n" +
	"\n" +
	"target_uri\x18\x01 \x01(\tR\ttargetUri\x12J\n" +
	"\x10call_credentials\x18\x02 \x01(\v2\x1d.gloo.solo.io.CallCredentialsH\x00R\x0fcallCredentials\x12#\n" +
	"\fcluster_name\x18\x05 \x01(\tH\x00R\vclusterName\x128\n" +
	"\x18certificates_secret_name\x18\x03 \x01(\tR\x16certificatesSecretName\x126\n" +
	"\x17validation_context_name\x18\x04 \x01(\tR\x15validationContextNameB\r\n" +
	"\vsds_builder\"\xd3\x01\n" +
	"\x0fCallCredentials\x12h\n" +
	"\x16file_credential_source\x18\x01 \x01(\v22.gloo.solo.io.CallCredentials.FileCredentialSourceR\x14fileCredentialSource\x1aV\n" +
	"\x14FileCredentialSource\x12&\n" +
	"\x0ftoken_file_name\x18\x01 \x01(\tR\rtokenFileName\x12\x16\n" +
	"\x06header\x18\x02 \x01(\tR\x06header\"\xf8\x02\n" +
	"\rSslParameters\x12e\n" +
	"\x18minimum_protocol_version\x18\x01 \x01(\x0e2+.gloo.solo.io.SslParameters.ProtocolVersionR\x16minimumProtocolVersion\x12e\n" +
	"\x18maximum_protocol_version\x18\x02 \x01(\x0e2+.gloo.solo.io.SslParameters.ProtocolVersionR\x16maximumProtocolVersion\x12#\n" +
	"\rcipher_suites\x18\x03 \x03(\tR\fcipherSuites\x12\x1f\n" +
	"\vecdh_curves\x18\x04 \x03(\tR\n" +
	"ecdhCurves\"S\n" +
	"\x0fProtocolVersion\x12\f\n" +
	"\bTLS_AUTO\x10\x00\x12\v\n" +
	"\aTLSv1_0\x10\x01\x12\v\n" +
	"\aTLSv1_1\x10\x02\x12\v\n" +
	"\aTLSv1_2\x10\x03\x12\v\n" +
	"\aTLSv1_3\x10\x04BB\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01Z4github.com/solo-io/gloo/projects/gloo/pkg/api/v1/sslb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_goTypes = []any{
	(SslConfig_OcspStaplePolicy)(0),              // 0: gloo.solo.io.SslConfig.OcspStaplePolicy
	(SslParameters_ProtocolVersion)(0),           // 1: gloo.solo.io.SslParameters.ProtocolVersion
	(*SslConfig)(nil),                            // 2: gloo.solo.io.SslConfig
	(*SSLFiles)(nil),                             // 3: gloo.solo.io.SSLFiles
	(*UpstreamSslConfig)(nil),                    // 4: gloo.solo.io.UpstreamSslConfig
	(*SDSConfig)(nil),                            // 5: gloo.solo.io.SDSConfig
	(*CallCredentials)(nil),                      // 6: gloo.solo.io.CallCredentials
	(*SslParameters)(nil),                        // 7: gloo.solo.io.SslParameters
	(*CallCredentials_FileCredentialSource)(nil), // 8: gloo.solo.io.CallCredentials.FileCredentialSource
	(*core.ResourceRef)(nil),                     // 9: core.solo.io.ResourceRef
	(*wrapperspb.BoolValue)(nil),                 // 10: google.protobuf.BoolValue
	(*durationpb.Duration)(nil),                  // 11: google.protobuf.Duration
}
var file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_depIdxs = []int32{
	9,  // 0: gloo.solo.io.SslConfig.secret_ref:type_name -> core.solo.io.ResourceRef
	3,  // 1: gloo.solo.io.SslConfig.ssl_files:type_name -> gloo.solo.io.SSLFiles
	5,  // 2: gloo.solo.io.SslConfig.sds:type_name -> gloo.solo.io.SDSConfig
	7,  // 3: gloo.solo.io.SslConfig.parameters:type_name -> gloo.solo.io.SslParameters
	10, // 4: gloo.solo.io.SslConfig.one_way_tls:type_name -> google.protobuf.BoolValue
	10, // 5: gloo.solo.io.SslConfig.disable_tls_session_resumption:type_name -> google.protobuf.BoolValue
	11, // 6: gloo.solo.io.SslConfig.transport_socket_connect_timeout:type_name -> google.protobuf.Duration
	0,  // 7: gloo.solo.io.SslConfig.ocsp_staple_policy:type_name -> gloo.solo.io.SslConfig.OcspStaplePolicy
	9,  // 8: gloo.solo.io.UpstreamSslConfig.secret_ref:type_name -> core.solo.io.ResourceRef
	3,  // 9: gloo.solo.io.UpstreamSslConfig.ssl_files:type_name -> gloo.solo.io.SSLFiles
	5,  // 10: gloo.solo.io.UpstreamSslConfig.sds:type_name -> gloo.solo.io.SDSConfig
	7,  // 11: gloo.solo.io.UpstreamSslConfig.parameters:type_name -> gloo.solo.io.SslParameters
	10, // 12: gloo.solo.io.UpstreamSslConfig.allow_renegotiation:type_name -> google.protobuf.BoolValue
	10, // 13: gloo.solo.io.UpstreamSslConfig.one_way_tls:type_name -> google.protobuf.BoolValue
	6,  // 14: gloo.solo.io.SDSConfig.call_credentials:type_name -> gloo.solo.io.CallCredentials
	8,  // 15: gloo.solo.io.CallCredentials.file_credential_source:type_name -> gloo.solo.io.CallCredentials.FileCredentialSource
	1,  // 16: gloo.solo.io.SslParameters.minimum_protocol_version:type_name -> gloo.solo.io.SslParameters.ProtocolVersion
	1,  // 17: gloo.solo.io.SslParameters.maximum_protocol_version:type_name -> gloo.solo.io.SslParameters.ProtocolVersion
	18, // [18:18] is the sub-list for method output_type
	18, // [18:18] is the sub-list for method input_type
	18, // [18:18] is the sub-list for extension type_name
	18, // [18:18] is the sub-list for extension extendee
	0,  // [0:18] is the sub-list for field type_name
}

func init() { file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_init() }
func file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto != nil {
		return
	}
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[0].OneofWrappers = []any{
		(*SslConfig_SecretRef)(nil),
		(*SslConfig_SslFiles)(nil),
		(*SslConfig_Sds)(nil),
	}
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[2].OneofWrappers = []any{
		(*UpstreamSslConfig_SecretRef)(nil),
		(*UpstreamSslConfig_SslFiles)(nil),
		(*UpstreamSslConfig_Sds)(nil),
	}
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes[3].OneofWrappers = []any{
		(*SDSConfig_CallCredentials)(nil),
		(*SDSConfig_ClusterName)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_depIdxs,
		EnumInfos:         file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_enumTypes,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_ssl_ssl_proto_depIdxs = nil
}
