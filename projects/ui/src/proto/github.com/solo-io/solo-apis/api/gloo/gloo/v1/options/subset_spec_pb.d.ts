/* eslint-disable */
// package: options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/subset_spec.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../extproto/ext_pb";

export class SubsetSpec extends jspb.Message {
  clearSelectorsList(): void;
  getSelectorsList(): Array<Selector>;
  setSelectorsList(value: Array<Selector>): void;
  addSelectors(value?: Selector, index?: number): Selector;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubsetSpec.AsObject;
  static toObject(includeInstance: boolean, msg: SubsetSpec): SubsetSpec.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SubsetSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubsetSpec;
  static deserializeBinaryFromReader(message: SubsetSpec, reader: jspb.BinaryReader): SubsetSpec;
}

export namespace SubsetSpec {
  export type AsObject = {
    selectorsList: Array<Selector.AsObject>,
  }
}

export class Selector extends jspb.Message {
  clearKeysList(): void;
  getKeysList(): Array<string>;
  setKeysList(value: Array<string>): void;
  addKeys(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Selector.AsObject;
  static toObject(includeInstance: boolean, msg: Selector): Selector.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Selector, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Selector;
  static deserializeBinaryFromReader(message: Selector, reader: jspb.BinaryReader): Selector;
}

export namespace Selector {
  export type AsObject = {
    keysList: Array<string>,
  }
}
