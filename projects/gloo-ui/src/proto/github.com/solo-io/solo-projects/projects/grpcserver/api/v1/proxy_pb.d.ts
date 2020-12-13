/* eslint-disable */
// package: glooeeapi.solo.io
// file: github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../protoc-gen-ext/extproto/ext_pb";
import * as github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb from "../../../../../../../github.com/solo-io/gloo/projects/gloo/api/v1/proxy_pb";
import * as github_com_solo_io_solo_kit_api_v1_ref_pb from "../../../../../../../github.com/solo-io/solo-kit/api/v1/ref_pb";
import * as github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb from "../../../../../../../github.com/solo-io/solo-projects/projects/grpcserver/api/v1/types_pb";

export class ProxyDetails extends jspb.Message {
  hasProxy(): boolean;
  clearProxy(): void;
  getProxy(): github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Proxy | undefined;
  setProxy(value?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Proxy): void;

  hasRaw(): boolean;
  clearRaw(): void;
  getRaw(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw | undefined;
  setRaw(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Status | undefined;
  setStatus(value?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Status): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProxyDetails.AsObject;
  static toObject(includeInstance: boolean, msg: ProxyDetails): ProxyDetails.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ProxyDetails, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProxyDetails;
  static deserializeBinaryFromReader(message: ProxyDetails, reader: jspb.BinaryReader): ProxyDetails;
}

export namespace ProxyDetails {
  export type AsObject = {
    proxy?: github_com_solo_io_gloo_projects_gloo_api_v1_proxy_pb.Proxy.AsObject,
    raw?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Raw.AsObject,
    status?: github_com_solo_io_solo_projects_projects_grpcserver_api_v1_types_pb.Status.AsObject,
  }
}

export class GetProxyRequest extends jspb.Message {
  hasRef(): boolean;
  clearRef(): void;
  getRef(): github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef | undefined;
  setRef(value?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProxyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProxyRequest): GetProxyRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProxyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProxyRequest;
  static deserializeBinaryFromReader(message: GetProxyRequest, reader: jspb.BinaryReader): GetProxyRequest;
}

export namespace GetProxyRequest {
  export type AsObject = {
    ref?: github_com_solo_io_solo_kit_api_v1_ref_pb.ResourceRef.AsObject,
  }
}

export class GetProxyResponse extends jspb.Message {
  hasProxyDetails(): boolean;
  clearProxyDetails(): void;
  getProxyDetails(): ProxyDetails | undefined;
  setProxyDetails(value?: ProxyDetails): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProxyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProxyResponse): GetProxyResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetProxyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProxyResponse;
  static deserializeBinaryFromReader(message: GetProxyResponse, reader: jspb.BinaryReader): GetProxyResponse;
}

export namespace GetProxyResponse {
  export type AsObject = {
    proxyDetails?: ProxyDetails.AsObject,
  }
}

export class ListProxiesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProxiesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProxiesRequest): ListProxiesRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProxiesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProxiesRequest;
  static deserializeBinaryFromReader(message: ListProxiesRequest, reader: jspb.BinaryReader): ListProxiesRequest;
}

export namespace ListProxiesRequest {
  export type AsObject = {
  }
}

export class ListProxiesResponse extends jspb.Message {
  clearProxyDetailsList(): void;
  getProxyDetailsList(): Array<ProxyDetails>;
  setProxyDetailsList(value: Array<ProxyDetails>): void;
  addProxyDetails(value?: ProxyDetails, index?: number): ProxyDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProxiesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProxiesResponse): ListProxiesResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ListProxiesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProxiesResponse;
  static deserializeBinaryFromReader(message: ListProxiesResponse, reader: jspb.BinaryReader): ListProxiesResponse;
}

export namespace ListProxiesResponse {
  export type AsObject = {
    proxyDetailsList: Array<ProxyDetails.AsObject>,
  }
}
