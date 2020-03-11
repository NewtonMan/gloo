/* eslint-disable */
// package: gloo.solo.io
// file: gloo/projects/gloo/api/v1/endpoint.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../protoc-gen-ext/extproto/ext_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as solo_kit_api_v1_ref_pb from "../../../../../solo-kit/api/v1/ref_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";

export class Endpoint extends jspb.Message {
  clearUpstreamsList(): void;
  getUpstreamsList(): Array<solo_kit_api_v1_ref_pb.ResourceRef>;
  setUpstreamsList(value: Array<solo_kit_api_v1_ref_pb.ResourceRef>): void;
  addUpstreams(value?: solo_kit_api_v1_ref_pb.ResourceRef, index?: number): solo_kit_api_v1_ref_pb.ResourceRef;

  getAddress(): string;
  setAddress(value: string): void;

  getPort(): number;
  setPort(value: number): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Endpoint.AsObject;
  static toObject(includeInstance: boolean, msg: Endpoint): Endpoint.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Endpoint, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Endpoint;
  static deserializeBinaryFromReader(message: Endpoint, reader: jspb.BinaryReader): Endpoint;
}

export namespace Endpoint {
  export type AsObject = {
    upstreamsList: Array<solo_kit_api_v1_ref_pb.ResourceRef.AsObject>,
    address: string,
    port: number,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
  }
}
