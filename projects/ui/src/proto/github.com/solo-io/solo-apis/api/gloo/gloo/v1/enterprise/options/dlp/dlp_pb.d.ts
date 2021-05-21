/* eslint-disable */
// package: dlp.options.gloo.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/options/dlp/dlp.proto

import * as jspb from "google-protobuf";
import * as github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb from "../../../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/core/matchers/matchers_pb";
import * as github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb from "../../../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/type/percent_pb";
import * as extproto_ext_pb from "../../../../../../../../../../extproto/ext_pb";

export class FilterConfig extends jspb.Message {
  clearDlpRulesList(): void;
  getDlpRulesList(): Array<DlpRule>;
  setDlpRulesList(value: Array<DlpRule>): void;
  addDlpRules(value?: DlpRule, index?: number): DlpRule;

  getEnabledFor(): FilterConfig.EnableForMap[keyof FilterConfig.EnableForMap];
  setEnabledFor(value: FilterConfig.EnableForMap[keyof FilterConfig.EnableForMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FilterConfig.AsObject;
  static toObject(includeInstance: boolean, msg: FilterConfig): FilterConfig.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: FilterConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FilterConfig;
  static deserializeBinaryFromReader(message: FilterConfig, reader: jspb.BinaryReader): FilterConfig;
}

export namespace FilterConfig {
  export type AsObject = {
    dlpRulesList: Array<DlpRule.AsObject>,
    enabledFor: FilterConfig.EnableForMap[keyof FilterConfig.EnableForMap],
  }

  export interface EnableForMap {
    RESPONSE_BODY: 0;
    ACCESS_LOGS: 1;
    ALL: 2;
  }

  export const EnableFor: EnableForMap;
}

export class DlpRule extends jspb.Message {
  hasMatcher(): boolean;
  clearMatcher(): void;
  getMatcher(): github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher | undefined;
  setMatcher(value?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher): void;

  clearActionsList(): void;
  getActionsList(): Array<Action>;
  setActionsList(value: Array<Action>): void;
  addActions(value?: Action, index?: number): Action;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DlpRule.AsObject;
  static toObject(includeInstance: boolean, msg: DlpRule): DlpRule.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DlpRule, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DlpRule;
  static deserializeBinaryFromReader(message: DlpRule, reader: jspb.BinaryReader): DlpRule;
}

export namespace DlpRule {
  export type AsObject = {
    matcher?: github_com_solo_io_solo_apis_api_gloo_gloo_v1_core_matchers_matchers_pb.Matcher.AsObject,
    actionsList: Array<Action.AsObject>,
  }
}

export class Config extends jspb.Message {
  clearActionsList(): void;
  getActionsList(): Array<Action>;
  setActionsList(value: Array<Action>): void;
  addActions(value?: Action, index?: number): Action;

  getEnabledFor(): Config.EnableForMap[keyof Config.EnableForMap];
  setEnabledFor(value: Config.EnableForMap[keyof Config.EnableForMap]): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Config.AsObject;
  static toObject(includeInstance: boolean, msg: Config): Config.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Config, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Config;
  static deserializeBinaryFromReader(message: Config, reader: jspb.BinaryReader): Config;
}

export namespace Config {
  export type AsObject = {
    actionsList: Array<Action.AsObject>,
    enabledFor: Config.EnableForMap[keyof Config.EnableForMap],
  }

  export interface EnableForMap {
    RESPONSE_BODY: 0;
    ACCESS_LOGS: 1;
    ALL: 2;
  }

  export const EnableFor: EnableForMap;
}

export class Action extends jspb.Message {
  getActionType(): Action.ActionTypeMap[keyof Action.ActionTypeMap];
  setActionType(value: Action.ActionTypeMap[keyof Action.ActionTypeMap]): void;

  hasCustomAction(): boolean;
  clearCustomAction(): void;
  getCustomAction(): CustomAction | undefined;
  setCustomAction(value?: CustomAction): void;

  getShadow(): boolean;
  setShadow(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Action.AsObject;
  static toObject(includeInstance: boolean, msg: Action): Action.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Action, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Action;
  static deserializeBinaryFromReader(message: Action, reader: jspb.BinaryReader): Action;
}

export namespace Action {
  export type AsObject = {
    actionType: Action.ActionTypeMap[keyof Action.ActionTypeMap],
    customAction?: CustomAction.AsObject,
    shadow: boolean,
  }

  export interface ActionTypeMap {
    CUSTOM: 0;
    SSN: 1;
    MASTERCARD: 2;
    VISA: 3;
    AMEX: 4;
    DISCOVER: 5;
    JCB: 6;
    DINERS_CLUB: 7;
    CREDIT_CARD_TRACKERS: 8;
    ALL_CREDIT_CARDS: 9;
  }

  export const ActionType: ActionTypeMap;
}

export class CustomAction extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  clearRegexList(): void;
  getRegexList(): Array<string>;
  setRegexList(value: Array<string>): void;
  addRegex(value: string, index?: number): string;

  getMaskChar(): string;
  setMaskChar(value: string): void;

  hasPercent(): boolean;
  clearPercent(): void;
  getPercent(): github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.Percent | undefined;
  setPercent(value?: github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.Percent): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CustomAction.AsObject;
  static toObject(includeInstance: boolean, msg: CustomAction): CustomAction.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: CustomAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CustomAction;
  static deserializeBinaryFromReader(message: CustomAction, reader: jspb.BinaryReader): CustomAction;
}

export namespace CustomAction {
  export type AsObject = {
    name: string,
    regexList: Array<string>,
    maskChar: string,
    percent?: github_com_solo_io_solo_kit_api_external_envoy_type_percent_pb.Percent.AsObject,
  }
}
