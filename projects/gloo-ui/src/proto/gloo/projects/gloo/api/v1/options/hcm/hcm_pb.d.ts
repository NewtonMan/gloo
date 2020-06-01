/* eslint-disable */
// package: hcm.options.gloo.solo.io
// file: gloo/projects/gloo/api/v1/options/hcm/hcm.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../../../gogoproto/gogo_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as gloo_projects_gloo_api_v1_options_tracing_tracing_pb from "../../../../../../../gloo/projects/gloo/api/v1/options/tracing/tracing_pb";
import * as gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb from "../../../../../../../gloo/projects/gloo/api/v1/options/protocol_upgrade/protocol_upgrade_pb";
import * as extproto_ext_pb from "../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class HttpConnectionManagerSettings extends jspb.Message {
  getSkipXffAppend(): boolean;
  setSkipXffAppend(value: boolean): void;

  getVia(): string;
  setVia(value: string): void;

  getXffNumTrustedHops(): number;
  setXffNumTrustedHops(value: number): void;

  hasUseRemoteAddress(): boolean;
  clearUseRemoteAddress(): void;
  getUseRemoteAddress(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setUseRemoteAddress(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasGenerateRequestId(): boolean;
  clearGenerateRequestId(): void;
  getGenerateRequestId(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setGenerateRequestId(value?: google_protobuf_wrappers_pb.BoolValue): void;

  getProxy100Continue(): boolean;
  setProxy100Continue(value: boolean): void;

  hasStreamIdleTimeout(): boolean;
  clearStreamIdleTimeout(): void;
  getStreamIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setStreamIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasIdleTimeout(): boolean;
  clearIdleTimeout(): void;
  getIdleTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setIdleTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasMaxRequestHeadersKb(): boolean;
  clearMaxRequestHeadersKb(): void;
  getMaxRequestHeadersKb(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setMaxRequestHeadersKb(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  hasRequestTimeout(): boolean;
  clearRequestTimeout(): void;
  getRequestTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setRequestTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasDrainTimeout(): boolean;
  clearDrainTimeout(): void;
  getDrainTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setDrainTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasDelayedCloseTimeout(): boolean;
  clearDelayedCloseTimeout(): void;
  getDelayedCloseTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setDelayedCloseTimeout(value?: google_protobuf_duration_pb.Duration): void;

  getServerName(): string;
  setServerName(value: string): void;

  getAcceptHttp10(): boolean;
  setAcceptHttp10(value: boolean): void;

  getDefaultHostForHttp10(): string;
  setDefaultHostForHttp10(value: string): void;

  getProperCaseHeaderKeyFormat(): boolean;
  setProperCaseHeaderKeyFormat(value: boolean): void;

  hasTracing(): boolean;
  clearTracing(): void;
  getTracing(): gloo_projects_gloo_api_v1_options_tracing_tracing_pb.ListenerTracingSettings | undefined;
  setTracing(value?: gloo_projects_gloo_api_v1_options_tracing_tracing_pb.ListenerTracingSettings): void;

  getForwardClientCertDetails(): HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap];
  setForwardClientCertDetails(value: HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap]): void;

  hasSetCurrentClientCertDetails(): boolean;
  clearSetCurrentClientCertDetails(): void;
  getSetCurrentClientCertDetails(): HttpConnectionManagerSettings.SetCurrentClientCertDetails | undefined;
  setSetCurrentClientCertDetails(value?: HttpConnectionManagerSettings.SetCurrentClientCertDetails): void;

  getPreserveExternalRequestId(): boolean;
  setPreserveExternalRequestId(value: boolean): void;

  clearUpgradesList(): void;
  getUpgradesList(): Array<gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>;
  setUpgradesList(value: Array<gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig>): void;
  addUpgrades(value?: gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig, index?: number): gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HttpConnectionManagerSettings.AsObject;
  static toObject(includeInstance: boolean, msg: HttpConnectionManagerSettings): HttpConnectionManagerSettings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: HttpConnectionManagerSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HttpConnectionManagerSettings;
  static deserializeBinaryFromReader(message: HttpConnectionManagerSettings, reader: jspb.BinaryReader): HttpConnectionManagerSettings;
}

export namespace HttpConnectionManagerSettings {
  export type AsObject = {
    skipXffAppend: boolean,
    via: string,
    xffNumTrustedHops: number,
    useRemoteAddress?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    generateRequestId?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    proxy100Continue: boolean,
    streamIdleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    idleTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    maxRequestHeadersKb?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
    requestTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    drainTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    delayedCloseTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    serverName: string,
    acceptHttp10: boolean,
    defaultHostForHttp10: string,
    properCaseHeaderKeyFormat: boolean,
    tracing?: gloo_projects_gloo_api_v1_options_tracing_tracing_pb.ListenerTracingSettings.AsObject,
    forwardClientCertDetails: HttpConnectionManagerSettings.ForwardClientCertDetailsMap[keyof HttpConnectionManagerSettings.ForwardClientCertDetailsMap],
    setCurrentClientCertDetails?: HttpConnectionManagerSettings.SetCurrentClientCertDetails.AsObject,
    preserveExternalRequestId: boolean,
    upgradesList: Array<gloo_projects_gloo_api_v1_options_protocol_upgrade_protocol_upgrade_pb.ProtocolUpgradeConfig.AsObject>,
  }

  export class SetCurrentClientCertDetails extends jspb.Message {
    hasSubject(): boolean;
    clearSubject(): void;
    getSubject(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setSubject(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getCert(): boolean;
    setCert(value: boolean): void;

    getChain(): boolean;
    setChain(value: boolean): void;

    getDns(): boolean;
    setDns(value: boolean): void;

    getUri(): boolean;
    setUri(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SetCurrentClientCertDetails.AsObject;
    static toObject(includeInstance: boolean, msg: SetCurrentClientCertDetails): SetCurrentClientCertDetails.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SetCurrentClientCertDetails, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SetCurrentClientCertDetails;
    static deserializeBinaryFromReader(message: SetCurrentClientCertDetails, reader: jspb.BinaryReader): SetCurrentClientCertDetails;
  }

  export namespace SetCurrentClientCertDetails {
    export type AsObject = {
      subject?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      cert: boolean,
      chain: boolean,
      dns: boolean,
      uri: boolean,
    }
  }

  export interface ForwardClientCertDetailsMap {
    SANITIZE: 0;
    FORWARD_ONLY: 1;
    APPEND_FORWARD: 2;
    SANITIZE_SET: 3;
    ALWAYS_FORWARD_ONLY: 4;
  }

  export const ForwardClientCertDetails: ForwardClientCertDetailsMap;
}
