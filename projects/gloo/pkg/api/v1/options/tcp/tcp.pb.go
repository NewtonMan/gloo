// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.6.1
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/tcp/tcp.proto

package tcp

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "github.com/solo-io/protoc-gen-ext/extproto"
	_ "github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
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

// Contains various settings for Envoy's tcp proxy filter.
// See here for more information: https://www.envoyproxy.io/docs/envoy/v1.10.0/api-v2/config/filter/network/tcp_proxy/v2/tcp_proxy.proto#envoy-api-msg-config-filter-network-tcp-proxy-v2-tcpproxy
type TcpProxySettings struct {
	state              protoimpl.MessageState  `protogen:"open.v1"`
	MaxConnectAttempts *wrapperspb.UInt32Value `protobuf:"bytes,1,opt,name=max_connect_attempts,json=maxConnectAttempts,proto3" json:"max_connect_attempts,omitempty"`
	IdleTimeout        *durationpb.Duration    `protobuf:"bytes,2,opt,name=idle_timeout,json=idleTimeout,proto3" json:"idle_timeout,omitempty"`
	// If set, this configures tunneling, e.g. configuration options to tunnel multiple TCP
	// payloads over a shared HTTP tunnel. If this message is absent, the payload
	// will be proxied upstream as per usual.
	TunnelingConfig *TcpProxySettings_TunnelingConfig `protobuf:"bytes,12,opt,name=tunneling_config,json=tunnelingConfig,proto3" json:"tunneling_config,omitempty"`
	// If set, Envoy will flush the access log on this time interval. Must be a
	// minimum of 1 ms. By default, will only write to the access log when a
	// connection is closed.
	AccessLogFlushInterval *durationpb.Duration `protobuf:"bytes,15,opt,name=access_log_flush_interval,json=accessLogFlushInterval,proto3" json:"access_log_flush_interval,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *TcpProxySettings) Reset() {
	*x = TcpProxySettings{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TcpProxySettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpProxySettings) ProtoMessage() {}

func (x *TcpProxySettings) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpProxySettings.ProtoReflect.Descriptor instead.
func (*TcpProxySettings) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescGZIP(), []int{0}
}

func (x *TcpProxySettings) GetMaxConnectAttempts() *wrapperspb.UInt32Value {
	if x != nil {
		return x.MaxConnectAttempts
	}
	return nil
}

func (x *TcpProxySettings) GetIdleTimeout() *durationpb.Duration {
	if x != nil {
		return x.IdleTimeout
	}
	return nil
}

func (x *TcpProxySettings) GetTunnelingConfig() *TcpProxySettings_TunnelingConfig {
	if x != nil {
		return x.TunnelingConfig
	}
	return nil
}

func (x *TcpProxySettings) GetAccessLogFlushInterval() *durationpb.Duration {
	if x != nil {
		return x.AccessLogFlushInterval
	}
	return nil
}

// Header name/value pair plus option to control append behavior.
type HeaderValueOption struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Header name/value pair that this option applies to.
	Header *HeaderValue `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// If true (default), the value is appended to existing values.
	Append        *wrapperspb.BoolValue `protobuf:"bytes,2,opt,name=append,proto3" json:"append,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeaderValueOption) Reset() {
	*x = HeaderValueOption{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeaderValueOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderValueOption) ProtoMessage() {}

func (x *HeaderValueOption) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderValueOption.ProtoReflect.Descriptor instead.
func (*HeaderValueOption) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescGZIP(), []int{1}
}

func (x *HeaderValueOption) GetHeader() *HeaderValue {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *HeaderValueOption) GetAppend() *wrapperspb.BoolValue {
	if x != nil {
		return x.Append
	}
	return nil
}

// Header name/value pair.
type HeaderValue struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Header name.
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Header value.
	Value         string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HeaderValue) Reset() {
	*x = HeaderValue{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HeaderValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderValue) ProtoMessage() {}

func (x *HeaderValue) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderValue.ProtoReflect.Descriptor instead.
func (*HeaderValue) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescGZIP(), []int{2}
}

func (x *HeaderValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *HeaderValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Configuration for tunneling TCP over other transports or application layers.
type TcpProxySettings_TunnelingConfig struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The hostname to send in the synthesized CONNECT headers to the upstream proxy.
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	// Additional request headers to be sent to upstream proxy. Mainly used to
	// trigger upstream to convert POST request back to CONNECT requests.
	HeadersToAdd  []*HeaderValueOption `protobuf:"bytes,13,rep,name=headers_to_add,json=headersToAdd,proto3" json:"headers_to_add,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TcpProxySettings_TunnelingConfig) Reset() {
	*x = TcpProxySettings_TunnelingConfig{}
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TcpProxySettings_TunnelingConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpProxySettings_TunnelingConfig) ProtoMessage() {}

func (x *TcpProxySettings_TunnelingConfig) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpProxySettings_TunnelingConfig.ProtoReflect.Descriptor instead.
func (*TcpProxySettings_TunnelingConfig) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescGZIP(), []int{0, 0}
}

func (x *TcpProxySettings_TunnelingConfig) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *TcpProxySettings_TunnelingConfig) GetHeadersToAdd() []*HeaderValueOption {
	if x != nil {
		return x.HeadersToAdd
	}
	return nil
}

var File_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto protoreflect.FileDescriptor

const file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDesc = "" +
	"\n" +
	"Bgithub.com/solo-io/gloo/projects/gloo/api/v1/options/tcp/tcp.proto\x12\x18tcp.options.gloo.solo.io\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x1egoogle/protobuf/duration.proto\x1aEgithub.com/solo-io/solo-kit/api/external/envoy/api/v2/core/base.proto\x1a\x12extproto/ext.proto\"\xe0\x03\n" +
	"\x10TcpProxySettings\x12N\n" +
	"\x14max_connect_attempts\x18\x01 \x01(\v2\x1c.google.protobuf.UInt32ValueR\x12maxConnectAttempts\x12<\n" +
	"\fidle_timeout\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\vidleTimeout\x12e\n" +
	"\x10tunneling_config\x18\f \x01(\v2:.tcp.options.gloo.solo.io.TcpProxySettings.TunnelingConfigR\x0ftunnelingConfig\x12T\n" +
	"\x19access_log_flush_interval\x18\x0f \x01(\v2\x19.google.protobuf.DurationR\x16accessLogFlushInterval\x1a\x80\x01\n" +
	"\x0fTunnelingConfig\x12\x1a\n" +
	"\bhostname\x18\x01 \x01(\tR\bhostname\x12Q\n" +
	"\x0eheaders_to_add\x18\r \x03(\v2+.tcp.options.gloo.solo.io.HeaderValueOptionR\fheadersToAdd\"\x86\x01\n" +
	"\x11HeaderValueOption\x12=\n" +
	"\x06header\x18\x01 \x01(\v2%.tcp.options.gloo.solo.io.HeaderValueR\x06header\x122\n" +
	"\x06append\x18\x02 \x01(\v2\x1a.google.protobuf.BoolValueR\x06append\"5\n" +
	"\vHeaderValue\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05valueBJ\xb8\xf5\x04\x01\xc0\xf5\x04\x01\xd0\xf5\x04\x01Z<github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tcpb\x06proto3"

var (
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescOnce sync.Once
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescData []byte
)

func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDesc)))
	})
	return file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDescData
}

var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_goTypes = []any{
	(*TcpProxySettings)(nil),                 // 0: tcp.options.gloo.solo.io.TcpProxySettings
	(*HeaderValueOption)(nil),                // 1: tcp.options.gloo.solo.io.HeaderValueOption
	(*HeaderValue)(nil),                      // 2: tcp.options.gloo.solo.io.HeaderValue
	(*TcpProxySettings_TunnelingConfig)(nil), // 3: tcp.options.gloo.solo.io.TcpProxySettings.TunnelingConfig
	(*wrapperspb.UInt32Value)(nil),           // 4: google.protobuf.UInt32Value
	(*durationpb.Duration)(nil),              // 5: google.protobuf.Duration
	(*wrapperspb.BoolValue)(nil),             // 6: google.protobuf.BoolValue
}
var file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_depIdxs = []int32{
	4, // 0: tcp.options.gloo.solo.io.TcpProxySettings.max_connect_attempts:type_name -> google.protobuf.UInt32Value
	5, // 1: tcp.options.gloo.solo.io.TcpProxySettings.idle_timeout:type_name -> google.protobuf.Duration
	3, // 2: tcp.options.gloo.solo.io.TcpProxySettings.tunneling_config:type_name -> tcp.options.gloo.solo.io.TcpProxySettings.TunnelingConfig
	5, // 3: tcp.options.gloo.solo.io.TcpProxySettings.access_log_flush_interval:type_name -> google.protobuf.Duration
	2, // 4: tcp.options.gloo.solo.io.HeaderValueOption.header:type_name -> tcp.options.gloo.solo.io.HeaderValue
	6, // 5: tcp.options.gloo.solo.io.HeaderValueOption.append:type_name -> google.protobuf.BoolValue
	1, // 6: tcp.options.gloo.solo.io.TcpProxySettings.TunnelingConfig.headers_to_add:type_name -> tcp.options.gloo.solo.io.HeaderValueOption
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_init() }
func file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_init() {
	if File_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDesc), len(file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto = out.File
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_goTypes = nil
	file_github_com_solo_io_gloo_projects_gloo_api_v1_options_tcp_tcp_proto_depIdxs = nil
}
