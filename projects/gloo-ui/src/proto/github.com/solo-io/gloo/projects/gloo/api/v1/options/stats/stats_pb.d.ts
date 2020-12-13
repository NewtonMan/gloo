/* eslint-disable */
// package: stats.options.gloo.solo.io
// file: github.com/solo-io/gloo/projects/gloo/api/v1/options/stats/stats.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../protoc-gen-ext/extproto/ext_pb";

export class Stats extends jspb.Message {
  clearVirtualClustersList(): void;
  getVirtualClustersList(): Array<VirtualCluster>;
  setVirtualClustersList(value: Array<VirtualCluster>): void;
  addVirtualClusters(value?: VirtualCluster, index?: number): VirtualCluster;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Stats.AsObject;
  static toObject(includeInstance: boolean, msg: Stats): Stats.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Stats, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Stats;
  static deserializeBinaryFromReader(message: Stats, reader: jspb.BinaryReader): Stats;
}

export namespace Stats {
  export type AsObject = {
    virtualClustersList: Array<VirtualCluster.AsObject>,
  }
}

export class VirtualCluster extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getPattern(): string;
  setPattern(value: string): void;

  getMethod(): string;
  setMethod(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VirtualCluster.AsObject;
  static toObject(includeInstance: boolean, msg: VirtualCluster): VirtualCluster.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: VirtualCluster, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VirtualCluster;
  static deserializeBinaryFromReader(message: VirtualCluster, reader: jspb.BinaryReader): VirtualCluster;
}

export namespace VirtualCluster {
  export type AsObject = {
    name: string,
    pattern: string,
    method: string,
  }
}
