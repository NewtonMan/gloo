/* eslint-disable */
// package: pipe.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/pipe/pipe.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_service_spec_pb from "../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/service_spec_pb";

export class UpstreamSpec extends jspb.Message {
  getPath(): string;
  setPath(value: string): void;

  hasServiceSpec(): boolean;
  clearServiceSpec(): void;
  getServiceSpec(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_service_spec_pb.ServiceSpec | undefined;
  setServiceSpec(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_service_spec_pb.ServiceSpec): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpstreamSpec.AsObject;
  static toObject(includeInstance: boolean, msg: UpstreamSpec): UpstreamSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: UpstreamSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpstreamSpec;
  static deserializeBinaryFromReader(message: UpstreamSpec, reader: jspb.BinaryReader): UpstreamSpec;
}

export namespace UpstreamSpec {
  export type AsObject = {
    path: string,
    serviceSpec?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_options_service_spec_pb.ServiceSpec.AsObject,
  }
}
