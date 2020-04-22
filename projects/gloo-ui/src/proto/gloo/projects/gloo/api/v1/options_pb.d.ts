/* eslint-disable */
// package: gloo.solo.io
// file: gloo/projects/gloo/api/v1/options.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../protoc-gen-ext/extproto/ext_pb";
import * as gloo_projects_gloo_api_v1_extensions_pb from "../../../../../gloo/projects/gloo/api/v1/extensions_pb";
import * as gloo_projects_gloo_api_v1_options_cors_cors_pb from "../../../../../gloo/projects/gloo/api/v1/options/cors/cors_pb";
import * as gloo_projects_gloo_api_v1_options_rest_rest_pb from "../../../../../gloo/projects/gloo/api/v1/options/rest/rest_pb";
import * as gloo_projects_gloo_api_v1_options_grpc_grpc_pb from "../../../../../gloo/projects/gloo/api/v1/options/grpc/grpc_pb";
import * as gloo_projects_gloo_api_v1_options_als_als_pb from "../../../../../gloo/projects/gloo/api/v1/options/als/als_pb";
import * as gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb from "../../../../../gloo/projects/gloo/api/v1/options/grpc_web/grpc_web_pb";
import * as gloo_projects_gloo_api_v1_options_hcm_hcm_pb from "../../../../../gloo/projects/gloo/api/v1/options/hcm/hcm_pb";
import * as gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb from "../../../../../gloo/projects/gloo/api/v1/options/lbhash/lbhash_pb";
import * as gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb from "../../../../../gloo/projects/gloo/api/v1/options/shadowing/shadowing_pb";
import * as gloo_projects_gloo_api_v1_options_tcp_tcp_pb from "../../../../../gloo/projects/gloo/api/v1/options/tcp/tcp_pb";
import * as gloo_projects_gloo_api_v1_options_tracing_tracing_pb from "../../../../../gloo/projects/gloo/api/v1/options/tracing/tracing_pb";
import * as gloo_projects_gloo_api_v1_options_retries_retries_pb from "../../../../../gloo/projects/gloo/api/v1/options/retries/retries_pb";
import * as gloo_projects_gloo_api_v1_options_stats_stats_pb from "../../../../../gloo/projects/gloo/api/v1/options/stats/stats_pb";
import * as gloo_projects_gloo_api_v1_options_faultinjection_fault_pb from "../../../../../gloo/projects/gloo/api/v1/options/faultinjection/fault_pb";
import * as gloo_projects_gloo_api_v1_options_headers_headers_pb from "../../../../../gloo/projects/gloo/api/v1/options/headers/headers_pb";
import * as gloo_projects_gloo_api_v1_options_aws_aws_pb from "../../../../../gloo/projects/gloo/api/v1/options/aws/aws_pb";
import * as gloo_projects_gloo_api_v1_options_wasm_wasm_pb from "../../../../../gloo/projects/gloo/api/v1/options/wasm/wasm_pb";
import * as gloo_projects_gloo_api_v1_options_azure_azure_pb from "../../../../../gloo/projects/gloo/api/v1/options/azure/azure_pb";
import * as gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb from "../../../../../gloo/projects/gloo/api/v1/options/healthcheck/healthcheck_pb";
import * as gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb from "../../../../../gloo/projects/gloo/api/v1/options/protocol_upgrade/protocol_upgrade_pb";
import * as gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb from "../../../../../gloo/projects/gloo/api/external/envoy/extensions/transformation/transformation_pb";
import * as gloo_projects_gloo_api_external_envoy_extensions_proxylatency_proxylatency_pb from "../../../../../gloo/projects/gloo/api/external/envoy/extensions/proxylatency/proxylatency_pb";
import * as gloo_projects_gloo_api_external_envoy_config_filter_http_gzip_v2_gzip_pb from "../../../../../gloo/projects/gloo/api/external/envoy/config/filter/http/gzip/v2/gzip_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/jwt/jwt_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/rbac/rbac_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/waf/waf_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/dlp/dlp_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class ListenerOptions extends jspb.Message {
  hasAccessLoggingService(): boolean;
  clearAccessLoggingService(): void;
  getAccessLoggingService(): gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService | undefined;
  setAccessLoggingService(value?: gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasPerConnectionBufferLimitBytes(): boolean;
  clearPerConnectionBufferLimitBytes(): void;
  getPerConnectionBufferLimitBytes(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setPerConnectionBufferLimitBytes(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListenerOptions.AsObject;
  static toObject(includeInstance: boolean, msg: ListenerOptions): ListenerOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListenerOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListenerOptions;
  static deserializeBinaryFromReader(message: ListenerOptions, reader: jspb.BinaryReader): ListenerOptions;
}

export namespace ListenerOptions {
  export type AsObject = {
    accessLoggingService?: gloo_projects_gloo_api_v1_options_als_als_pb.AccessLoggingService.AsObject,
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    perConnectionBufferLimitBytes?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }
}

export class HttpListenerOptions extends jspb.Message {
  hasGrpcWeb(): boolean;
  clearGrpcWeb(): void;
  getGrpcWeb(): gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb | undefined;
  setGrpcWeb(value?: gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb): void;

  hasHttpConnectionManagerSettings(): boolean;
  clearHttpConnectionManagerSettings(): void;
  getHttpConnectionManagerSettings(): gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings | undefined;
  setHttpConnectionManagerSettings(value?: gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings): void;

  hasHealthCheck(): boolean;
  clearHealthCheck(): void;
  getHealthCheck(): gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck | undefined;
  setHealthCheck(value?: gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasWaf(): boolean;
  clearWaf(): void;
  getWaf(): gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings | undefined;
  setWaf(value?: gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings): void;

  hasDlp(): boolean;
  clearDlp(): void;
  getDlp(): gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig | undefined;
  setDlp(value?: gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig): void;

  hasWasm(): boolean;
  clearWasm(): void;
  getWasm(): gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource | undefined;
  setWasm(value?: gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource): void;

  hasGzip(): boolean;
  clearGzip(): void;
  getGzip(): gloo_projects_gloo_api_external_envoy_config_filter_http_gzip_v2_gzip_pb.Gzip | undefined;
  setGzip(value?: gloo_projects_gloo_api_external_envoy_config_filter_http_gzip_v2_gzip_pb.Gzip): void;

  hasProxyLatency(): boolean;
  clearProxyLatency(): void;
  getProxyLatency(): gloo_projects_gloo_api_external_envoy_extensions_proxylatency_proxylatency_pb.ProxyLatency | undefined;
  setProxyLatency(value?: gloo_projects_gloo_api_external_envoy_extensions_proxylatency_proxylatency_pb.ProxyLatency): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpListenerOptions.AsObject;
  static toObject(includeInstance: boolean, msg: HttpListenerOptions): HttpListenerOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpListenerOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpListenerOptions;
  static deserializeBinaryFromReader(message: HttpListenerOptions, reader: jspb.BinaryReader): HttpListenerOptions;
}

export namespace HttpListenerOptions {
  export type AsObject = {
    grpcWeb?: gloo_projects_gloo_api_v1_options_grpc_web_grpc_web_pb.GrpcWeb.AsObject,
    httpConnectionManagerSettings?: gloo_projects_gloo_api_v1_options_hcm_hcm_pb.HttpConnectionManagerSettings.AsObject,
    healthCheck?: gloo_projects_gloo_api_v1_options_healthcheck_healthcheck_pb.HealthCheck.AsObject,
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    waf?: gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.AsObject,
    dlp?: gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.FilterConfig.AsObject,
    wasm?: gloo_projects_gloo_api_v1_options_wasm_wasm_pb.PluginSource.AsObject,
    gzip?: gloo_projects_gloo_api_external_envoy_config_filter_http_gzip_v2_gzip_pb.Gzip.AsObject,
    proxyLatency?: gloo_projects_gloo_api_external_envoy_extensions_proxylatency_proxylatency_pb.ProxyLatency.AsObject,
  }
}

export class TcpListenerOptions extends jspb.Message {
  hasTcpProxySettings(): boolean;
  clearTcpProxySettings(): void;
  getTcpProxySettings(): gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings | undefined;
  setTcpProxySettings(value?: gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TcpListenerOptions.AsObject;
  static toObject(includeInstance: boolean, msg: TcpListenerOptions): TcpListenerOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TcpListenerOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TcpListenerOptions;
  static deserializeBinaryFromReader(message: TcpListenerOptions, reader: jspb.BinaryReader): TcpListenerOptions;
}

export namespace TcpListenerOptions {
  export type AsObject = {
    tcpProxySettings?: gloo_projects_gloo_api_v1_options_tcp_tcp_pb.TcpProxySettings.AsObject,
  }
}

export class VirtualHostOptions extends jspb.Message {
  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasRetries(): boolean;
  clearRetries(): void;
  getRetries(): gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy | undefined;
  setRetries(value?: gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy): void;

  hasStats(): boolean;
  clearStats(): void;
  getStats(): gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats | undefined;
  setStats(value?: gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats): void;

  hasHeaderManipulation(): boolean;
  clearHeaderManipulation(): void;
  getHeaderManipulation(): gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation | undefined;
  setHeaderManipulation(value?: gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation): void;

  hasCors(): boolean;
  clearCors(): void;
  getCors(): gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy | undefined;
  setCors(value?: gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy): void;

  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations | undefined;
  setTransformations(value?: gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations): void;

  hasRatelimitBasic(): boolean;
  clearRatelimitBasic(): void;
  getRatelimitBasic(): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit | undefined;
  setRatelimitBasic(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit): void;

  hasRatelimit(): boolean;
  clearRatelimit(): void;
  getRatelimit(): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension | undefined;
  setRatelimit(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension): void;

  hasWaf(): boolean;
  clearWaf(): void;
  getWaf(): gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings | undefined;
  setWaf(value?: gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings): void;

  hasJwt(): boolean;
  clearJwt(): void;
  getJwt(): gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension | undefined;
  setJwt(value?: gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension): void;

  hasRbac(): boolean;
  clearRbac(): void;
  getRbac(): gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings | undefined;
  setRbac(value?: gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension | undefined;
  setExtauth(value?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension): void;

  hasDlp(): boolean;
  clearDlp(): void;
  getDlp(): gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config | undefined;
  setDlp(value?: gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualHostOptions.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualHostOptions): VirtualHostOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualHostOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualHostOptions;
  static deserializeBinaryFromReader(message: VirtualHostOptions, reader: jspb.BinaryReader): VirtualHostOptions;
}

export namespace VirtualHostOptions {
  export type AsObject = {
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    retries?: gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.AsObject,
    stats?: gloo_projects_gloo_api_v1_options_stats_stats_pb.Stats.AsObject,
    headerManipulation?: gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.AsObject,
    cors?: gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.AsObject,
    transformations?: gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.AsObject,
    ratelimitBasic?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.AsObject,
    ratelimit?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitVhostExtension.AsObject,
    waf?: gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.AsObject,
    jwt?: gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.VhostExtension.AsObject,
    rbac?: gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.AsObject,
    extauth?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.AsObject,
    dlp?: gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.AsObject,
  }
}

export class RouteOptions extends jspb.Message {
  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations | undefined;
  setTransformations(value?: gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations): void;

  hasFaults(): boolean;
  clearFaults(): void;
  getFaults(): gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults | undefined;
  setFaults(value?: gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults): void;

  hasPrefixRewrite(): boolean;
  clearPrefixRewrite(): void;
  getPrefixRewrite(): google_protobuf_wrappers_pb.StringValue | undefined;
  setPrefixRewrite(value?: google_protobuf_wrappers_pb.StringValue): void;

  hasTimeout(): boolean;
  clearTimeout(): void;
  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasRetries(): boolean;
  clearRetries(): void;
  getRetries(): gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy | undefined;
  setRetries(value?: gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasTracing(): boolean;
  clearTracing(): void;
  getTracing(): gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings | undefined;
  setTracing(value?: gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings): void;

  hasShadowing(): boolean;
  clearShadowing(): void;
  getShadowing(): gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing | undefined;
  setShadowing(value?: gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing): void;

  hasHeaderManipulation(): boolean;
  clearHeaderManipulation(): void;
  getHeaderManipulation(): gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation | undefined;
  setHeaderManipulation(value?: gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation): void;

  hasHostRewrite(): boolean;
  clearHostRewrite(): void;
  getHostRewrite(): string;
  setHostRewrite(value: string): void;

  hasAutoHostRewrite(): boolean;
  clearAutoHostRewrite(): void;
  getAutoHostRewrite(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setAutoHostRewrite(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasCors(): boolean;
  clearCors(): void;
  getCors(): gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy | undefined;
  setCors(value?: gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy): void;

  hasLbHash(): boolean;
  clearLbHash(): void;
  getLbHash(): gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig | undefined;
  setLbHash(value?: gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig): void;

  clearUpgradesList(): void;
  getUpgradesList(): Array<gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>;
  setUpgradesList(value: Array<gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>): void;
  addUpgrades(value?: gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig, index?: number): gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig;

  hasRatelimitBasic(): boolean;
  clearRatelimitBasic(): void;
  getRatelimitBasic(): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit | undefined;
  setRatelimitBasic(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit): void;

  hasRatelimit(): boolean;
  clearRatelimit(): void;
  getRatelimit(): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension | undefined;
  setRatelimit(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension): void;

  hasWaf(): boolean;
  clearWaf(): void;
  getWaf(): gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings | undefined;
  setWaf(value?: gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings): void;

  hasJwt(): boolean;
  clearJwt(): void;
  getJwt(): gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension | undefined;
  setJwt(value?: gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension): void;

  hasRbac(): boolean;
  clearRbac(): void;
  getRbac(): gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings | undefined;
  setRbac(value?: gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension | undefined;
  setExtauth(value?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension): void;

  hasDlp(): boolean;
  clearDlp(): void;
  getDlp(): gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config | undefined;
  setDlp(value?: gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config): void;

  getHostRewriteTypeCase(): RouteOptions.HostRewriteTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RouteOptions.AsObject;
  static toObject(includeInstance: boolean, msg: RouteOptions): RouteOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: RouteOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RouteOptions;
  static deserializeBinaryFromReader(message: RouteOptions, reader: jspb.BinaryReader): RouteOptions;
}

export namespace RouteOptions {
  export type AsObject = {
    transformations?: gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.AsObject,
    faults?: gloo_projects_gloo_api_v1_options_faultinjection_fault_pb.RouteFaults.AsObject,
    prefixRewrite?: google_protobuf_wrappers_pb.StringValue.AsObject,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    retries?: gloo_projects_gloo_api_v1_options_retries_retries_pb.RetryPolicy.AsObject,
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    tracing?: gloo_projects_gloo_api_v1_options_tracing_tracing_pb.RouteTracingSettings.AsObject,
    shadowing?: gloo_projects_gloo_api_v1_options_shadowing_shadowing_pb.RouteShadowing.AsObject,
    headerManipulation?: gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.AsObject,
    hostRewrite: string,
    autoHostRewrite?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    cors?: gloo_projects_gloo_api_v1_options_cors_cors_pb.CorsPolicy.AsObject,
    lbHash?: gloo_projects_gloo_api_v1_options_lbhash_lbhash_pb.RouteActionHashConfig.AsObject,
    upgradesList: Array<gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.AsObject>,
    ratelimitBasic?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.IngressRateLimit.AsObject,
    ratelimit?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.RateLimitRouteExtension.AsObject,
    waf?: gloo_projects_gloo_api_v1_enterprise_options_waf_waf_pb.Settings.AsObject,
    jwt?: gloo_projects_gloo_api_v1_enterprise_options_jwt_jwt_pb.RouteExtension.AsObject,
    rbac?: gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.ExtensionSettings.AsObject,
    extauth?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.AsObject,
    dlp?: gloo_projects_gloo_api_v1_enterprise_options_dlp_dlp_pb.Config.AsObject,
  }

  export enum HostRewriteTypeCase {
    HOST_REWRITE_TYPE_NOT_SET = 0,
    HOST_REWRITE = 10,
    AUTO_HOST_REWRITE = 19,
  }
}

export class DestinationSpec extends jspb.Message {
  hasAws(): boolean;
  clearAws(): void;
  getAws(): gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec | undefined;
  setAws(value?: gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec): void;

  hasAzure(): boolean;
  clearAzure(): void;
  getAzure(): gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec | undefined;
  setAzure(value?: gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec): void;

  hasRest(): boolean;
  clearRest(): void;
  getRest(): gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec | undefined;
  setRest(value?: gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec): void;

  hasGrpc(): boolean;
  clearGrpc(): void;
  getGrpc(): gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec | undefined;
  setGrpc(value?: gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec): void;

  getDestinationTypeCase(): DestinationSpec.DestinationTypeCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DestinationSpec.AsObject;
  static toObject(includeInstance: boolean, msg: DestinationSpec): DestinationSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DestinationSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DestinationSpec;
  static deserializeBinaryFromReader(message: DestinationSpec, reader: jspb.BinaryReader): DestinationSpec;
}

export namespace DestinationSpec {
  export type AsObject = {
    aws?: gloo_projects_gloo_api_v1_options_aws_aws_pb.DestinationSpec.AsObject,
    azure?: gloo_projects_gloo_api_v1_options_azure_azure_pb.DestinationSpec.AsObject,
    rest?: gloo_projects_gloo_api_v1_options_rest_rest_pb.DestinationSpec.AsObject,
    grpc?: gloo_projects_gloo_api_v1_options_grpc_grpc_pb.DestinationSpec.AsObject,
  }

  export enum DestinationTypeCase {
    DESTINATION_TYPE_NOT_SET = 0,
    AWS = 1,
    AZURE = 2,
    REST = 3,
    GRPC = 4,
  }
}

export class WeightedDestinationOptions extends jspb.Message {
  hasHeaderManipulation(): boolean;
  clearHeaderManipulation(): void;
  getHeaderManipulation(): gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation | undefined;
  setHeaderManipulation(value?: gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation): void;

  hasTransformations(): boolean;
  clearTransformations(): void;
  getTransformations(): gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations | undefined;
  setTransformations(value?: gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension | undefined;
  setExtauth(value?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WeightedDestinationOptions.AsObject;
  static toObject(includeInstance: boolean, msg: WeightedDestinationOptions): WeightedDestinationOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: WeightedDestinationOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WeightedDestinationOptions;
  static deserializeBinaryFromReader(message: WeightedDestinationOptions, reader: jspb.BinaryReader): WeightedDestinationOptions;
}

export namespace WeightedDestinationOptions {
  export type AsObject = {
    headerManipulation?: gloo_projects_gloo_api_v1_options_headers_headers_pb.HeaderManipulation.AsObject,
    transformations?: gloo_projects_gloo_api_external_envoy_extensions_transformation_transformation_pb.RouteTransformations.AsObject,
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    extauth?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.ExtAuthExtension.AsObject,
  }
}
