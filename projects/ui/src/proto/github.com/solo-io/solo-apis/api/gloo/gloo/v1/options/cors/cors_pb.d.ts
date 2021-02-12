/* eslint-disable */
// package: cors.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/options/cors/cors.proto

import * as jspb from "google-protobuf";
import * as extproto_ext_pb from "../../../../../../../../../extproto/ext_pb";

export class CorsPolicy extends jspb.Message {
  clearAllowOriginList(): void;
  getAllowOriginList(): Array<string>;
  setAllowOriginList(value: Array<string>): void;
  addAllowOrigin(value: string, index?: number): string;

  clearAllowOriginRegexList(): void;
  getAllowOriginRegexList(): Array<string>;
  setAllowOriginRegexList(value: Array<string>): void;
  addAllowOriginRegex(value: string, index?: number): string;

  clearAllowMethodsList(): void;
  getAllowMethodsList(): Array<string>;
  setAllowMethodsList(value: Array<string>): void;
  addAllowMethods(value: string, index?: number): string;

  clearAllowHeadersList(): void;
  getAllowHeadersList(): Array<string>;
  setAllowHeadersList(value: Array<string>): void;
  addAllowHeaders(value: string, index?: number): string;

  clearExposeHeadersList(): void;
  getExposeHeadersList(): Array<string>;
  setExposeHeadersList(value: Array<string>): void;
  addExposeHeaders(value: string, index?: number): string;

  getMaxAge(): string;
  setMaxAge(value: string): void;

  getAllowCredentials(): boolean;
  setAllowCredentials(value: boolean): void;

  getDisableForRoute(): boolean;
  setDisableForRoute(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CorsPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: CorsPolicy): CorsPolicy.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CorsPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CorsPolicy;
  static deserializeBinaryFromReader(message: CorsPolicy, reader: jspb.BinaryReader): CorsPolicy;
}

export namespace CorsPolicy {
  export type AsObject = {
    allowOriginList: Array<string>,
    allowOriginRegexList: Array<string>,
    allowMethodsList: Array<string>,
    allowHeadersList: Array<string>,
    exposeHeadersList: Array<string>,
    maxAge: string,
    allowCredentials: boolean,
    disableForRoute: boolean,
  }
}
