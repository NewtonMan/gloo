/* eslint-disable */
// package: gloo.solo.io
// file: gloo/projects/gloo/api/v1/settings.proto

import * as jspb from "google-protobuf";
import * as gogoproto_gogo_pb from "../../../../../gogoproto/gogo_pb";
import * as extproto_ext_pb from "../../../../../protoc-gen-ext/extproto/ext_pb";
import * as solo_kit_api_v1_metadata_pb from "../../../../../solo-kit/api/v1/metadata_pb";
import * as solo_kit_api_v1_status_pb from "../../../../../solo-kit/api/v1/status_pb";
import * as solo_kit_api_v1_solo_kit_pb from "../../../../../solo-kit/api/v1/solo-kit_pb";
import * as gloo_projects_gloo_api_v1_extensions_pb from "../../../../../gloo/projects/gloo/api/v1/extensions_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/extauth/v1/extauth_pb";
import * as gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb from "../../../../../gloo/projects/gloo/api/v1/enterprise/options/rbac/rbac_pb";
import * as gloo_projects_gloo_api_v1_circuit_breaker_pb from "../../../../../gloo/projects/gloo/api/v1/circuit_breaker_pb";
import * as google_protobuf_duration_pb from "google-protobuf/google/protobuf/duration_pb";
import * as google_protobuf_wrappers_pb from "google-protobuf/google/protobuf/wrappers_pb";

export class Settings extends jspb.Message {
  getDiscoveryNamespace(): string;
  setDiscoveryNamespace(value: string): void;

  clearWatchNamespacesList(): void;
  getWatchNamespacesList(): Array<string>;
  setWatchNamespacesList(value: Array<string>): void;
  addWatchNamespaces(value: string, index?: number): string;

  hasKubernetesConfigSource(): boolean;
  clearKubernetesConfigSource(): void;
  getKubernetesConfigSource(): Settings.KubernetesCrds | undefined;
  setKubernetesConfigSource(value?: Settings.KubernetesCrds): void;

  hasDirectoryConfigSource(): boolean;
  clearDirectoryConfigSource(): void;
  getDirectoryConfigSource(): Settings.Directory | undefined;
  setDirectoryConfigSource(value?: Settings.Directory): void;

  hasConsulKvSource(): boolean;
  clearConsulKvSource(): void;
  getConsulKvSource(): Settings.ConsulKv | undefined;
  setConsulKvSource(value?: Settings.ConsulKv): void;

  hasKubernetesSecretSource(): boolean;
  clearKubernetesSecretSource(): void;
  getKubernetesSecretSource(): Settings.KubernetesSecrets | undefined;
  setKubernetesSecretSource(value?: Settings.KubernetesSecrets): void;

  hasVaultSecretSource(): boolean;
  clearVaultSecretSource(): void;
  getVaultSecretSource(): Settings.VaultSecrets | undefined;
  setVaultSecretSource(value?: Settings.VaultSecrets): void;

  hasDirectorySecretSource(): boolean;
  clearDirectorySecretSource(): void;
  getDirectorySecretSource(): Settings.Directory | undefined;
  setDirectorySecretSource(value?: Settings.Directory): void;

  hasKubernetesArtifactSource(): boolean;
  clearKubernetesArtifactSource(): void;
  getKubernetesArtifactSource(): Settings.KubernetesConfigmaps | undefined;
  setKubernetesArtifactSource(value?: Settings.KubernetesConfigmaps): void;

  hasDirectoryArtifactSource(): boolean;
  clearDirectoryArtifactSource(): void;
  getDirectoryArtifactSource(): Settings.Directory | undefined;
  setDirectoryArtifactSource(value?: Settings.Directory): void;

  hasConsulKvArtifactSource(): boolean;
  clearConsulKvArtifactSource(): void;
  getConsulKvArtifactSource(): Settings.ConsulKv | undefined;
  setConsulKvArtifactSource(value?: Settings.ConsulKv): void;

  hasRefreshRate(): boolean;
  clearRefreshRate(): void;
  getRefreshRate(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshRate(value?: google_protobuf_duration_pb.Duration): void;

  getDevMode(): boolean;
  setDevMode(value: boolean): void;

  getLinkerd(): boolean;
  setLinkerd(value: boolean): void;

  hasKnative(): boolean;
  clearKnative(): void;
  getKnative(): Settings.KnativeOptions | undefined;
  setKnative(value?: Settings.KnativeOptions): void;

  hasDiscovery(): boolean;
  clearDiscovery(): void;
  getDiscovery(): Settings.DiscoveryOptions | undefined;
  setDiscovery(value?: Settings.DiscoveryOptions): void;

  hasGloo(): boolean;
  clearGloo(): void;
  getGloo(): GlooOptions | undefined;
  setGloo(value?: GlooOptions): void;

  hasGateway(): boolean;
  clearGateway(): void;
  getGateway(): GatewayOptions | undefined;
  setGateway(value?: GatewayOptions): void;

  hasConsul(): boolean;
  clearConsul(): void;
  getConsul(): Settings.ConsulConfiguration | undefined;
  setConsul(value?: Settings.ConsulConfiguration): void;

  hasKubernetes(): boolean;
  clearKubernetes(): void;
  getKubernetes(): Settings.KubernetesConfiguration | undefined;
  setKubernetes(value?: Settings.KubernetesConfiguration): void;

  hasExtensions(): boolean;
  clearExtensions(): void;
  getExtensions(): gloo_projects_gloo_api_v1_extensions_pb.Extensions | undefined;
  setExtensions(value?: gloo_projects_gloo_api_v1_extensions_pb.Extensions): void;

  hasRatelimit(): boolean;
  clearRatelimit(): void;
  getRatelimit(): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.ServiceSettings | undefined;
  setRatelimit(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.ServiceSettings): void;

  hasRatelimitServer(): boolean;
  clearRatelimitServer(): void;
  getRatelimitServer(): gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Settings | undefined;
  setRatelimitServer(value?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Settings): void;

  hasRbac(): boolean;
  clearRbac(): void;
  getRbac(): gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.Settings | undefined;
  setRbac(value?: gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.Settings): void;

  hasExtauth(): boolean;
  clearExtauth(): void;
  getExtauth(): gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.Settings | undefined;
  setExtauth(value?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.Settings): void;

  hasMetadata(): boolean;
  clearMetadata(): void;
  getMetadata(): solo_kit_api_v1_metadata_pb.Metadata | undefined;
  setMetadata(value?: solo_kit_api_v1_metadata_pb.Metadata): void;

  hasStatus(): boolean;
  clearStatus(): void;
  getStatus(): solo_kit_api_v1_status_pb.Status | undefined;
  setStatus(value?: solo_kit_api_v1_status_pb.Status): void;

  getConfigSourceCase(): Settings.ConfigSourceCase;
  getSecretSourceCase(): Settings.SecretSourceCase;
  getArtifactSourceCase(): Settings.ArtifactSourceCase;
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Settings.AsObject;
  static toObject(includeInstance: boolean, msg: Settings): Settings.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Settings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Settings;
  static deserializeBinaryFromReader(message: Settings, reader: jspb.BinaryReader): Settings;
}

export namespace Settings {
  export type AsObject = {
    discoveryNamespace: string,
    watchNamespacesList: Array<string>,
    kubernetesConfigSource?: Settings.KubernetesCrds.AsObject,
    directoryConfigSource?: Settings.Directory.AsObject,
    consulKvSource?: Settings.ConsulKv.AsObject,
    kubernetesSecretSource?: Settings.KubernetesSecrets.AsObject,
    vaultSecretSource?: Settings.VaultSecrets.AsObject,
    directorySecretSource?: Settings.Directory.AsObject,
    kubernetesArtifactSource?: Settings.KubernetesConfigmaps.AsObject,
    directoryArtifactSource?: Settings.Directory.AsObject,
    consulKvArtifactSource?: Settings.ConsulKv.AsObject,
    refreshRate?: google_protobuf_duration_pb.Duration.AsObject,
    devMode: boolean,
    linkerd: boolean,
    knative?: Settings.KnativeOptions.AsObject,
    discovery?: Settings.DiscoveryOptions.AsObject,
    gloo?: GlooOptions.AsObject,
    gateway?: GatewayOptions.AsObject,
    consul?: Settings.ConsulConfiguration.AsObject,
    kubernetes?: Settings.KubernetesConfiguration.AsObject,
    extensions?: gloo_projects_gloo_api_v1_extensions_pb.Extensions.AsObject,
    ratelimit?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.ServiceSettings.AsObject,
    ratelimitServer?: gloo_projects_gloo_api_v1_enterprise_options_ratelimit_ratelimit_pb.Settings.AsObject,
    rbac?: gloo_projects_gloo_api_v1_enterprise_options_rbac_rbac_pb.Settings.AsObject,
    extauth?: gloo_projects_gloo_api_v1_enterprise_options_extauth_v1_extauth_pb.Settings.AsObject,
    metadata?: solo_kit_api_v1_metadata_pb.Metadata.AsObject,
    status?: solo_kit_api_v1_status_pb.Status.AsObject,
  }

  export class KubernetesCrds extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesCrds.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesCrds): KubernetesCrds.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesCrds, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesCrds;
    static deserializeBinaryFromReader(message: KubernetesCrds, reader: jspb.BinaryReader): KubernetesCrds;
  }

  export namespace KubernetesCrds {
    export type AsObject = {
    }
  }

  export class KubernetesSecrets extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesSecrets.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesSecrets): KubernetesSecrets.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesSecrets, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesSecrets;
    static deserializeBinaryFromReader(message: KubernetesSecrets, reader: jspb.BinaryReader): KubernetesSecrets;
  }

  export namespace KubernetesSecrets {
    export type AsObject = {
    }
  }

  export class VaultSecrets extends jspb.Message {
    getToken(): string;
    setToken(value: string): void;

    getAddress(): string;
    setAddress(value: string): void;

    getCaCert(): string;
    setCaCert(value: string): void;

    getCaPath(): string;
    setCaPath(value: string): void;

    getClientCert(): string;
    setClientCert(value: string): void;

    getClientKey(): string;
    setClientKey(value: string): void;

    getTlsServerName(): string;
    setTlsServerName(value: string): void;

    hasInsecure(): boolean;
    clearInsecure(): void;
    getInsecure(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setInsecure(value?: google_protobuf_wrappers_pb.BoolValue): void;

    getRootKey(): string;
    setRootKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): VaultSecrets.AsObject;
    static toObject(includeInstance: boolean, msg: VaultSecrets): VaultSecrets.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: VaultSecrets, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): VaultSecrets;
    static deserializeBinaryFromReader(message: VaultSecrets, reader: jspb.BinaryReader): VaultSecrets;
  }

  export namespace VaultSecrets {
    export type AsObject = {
      token: string,
      address: string,
      caCert: string,
      caPath: string,
      clientCert: string,
      clientKey: string,
      tlsServerName: string,
      insecure?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      rootKey: string,
    }
  }

  export class ConsulKv extends jspb.Message {
    getRootKey(): string;
    setRootKey(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConsulKv.AsObject;
    static toObject(includeInstance: boolean, msg: ConsulKv): ConsulKv.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConsulKv, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConsulKv;
    static deserializeBinaryFromReader(message: ConsulKv, reader: jspb.BinaryReader): ConsulKv;
  }

  export namespace ConsulKv {
    export type AsObject = {
      rootKey: string,
    }
  }

  export class KubernetesConfigmaps extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesConfigmaps.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesConfigmaps): KubernetesConfigmaps.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesConfigmaps, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesConfigmaps;
    static deserializeBinaryFromReader(message: KubernetesConfigmaps, reader: jspb.BinaryReader): KubernetesConfigmaps;
  }

  export namespace KubernetesConfigmaps {
    export type AsObject = {
    }
  }

  export class Directory extends jspb.Message {
    getDirectory(): string;
    setDirectory(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Directory.AsObject;
    static toObject(includeInstance: boolean, msg: Directory): Directory.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Directory, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Directory;
    static deserializeBinaryFromReader(message: Directory, reader: jspb.BinaryReader): Directory;
  }

  export namespace Directory {
    export type AsObject = {
      directory: string,
    }
  }

  export class KnativeOptions extends jspb.Message {
    getClusterIngressProxyAddress(): string;
    setClusterIngressProxyAddress(value: string): void;

    getKnativeExternalProxyAddress(): string;
    setKnativeExternalProxyAddress(value: string): void;

    getKnativeInternalProxyAddress(): string;
    setKnativeInternalProxyAddress(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KnativeOptions.AsObject;
    static toObject(includeInstance: boolean, msg: KnativeOptions): KnativeOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KnativeOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KnativeOptions;
    static deserializeBinaryFromReader(message: KnativeOptions, reader: jspb.BinaryReader): KnativeOptions;
  }

  export namespace KnativeOptions {
    export type AsObject = {
      clusterIngressProxyAddress: string,
      knativeExternalProxyAddress: string,
      knativeInternalProxyAddress: string,
    }
  }

  export class DiscoveryOptions extends jspb.Message {
    getFdsMode(): Settings.DiscoveryOptions.FdsModeMap[keyof Settings.DiscoveryOptions.FdsModeMap];
    setFdsMode(value: Settings.DiscoveryOptions.FdsModeMap[keyof Settings.DiscoveryOptions.FdsModeMap]): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DiscoveryOptions.AsObject;
    static toObject(includeInstance: boolean, msg: DiscoveryOptions): DiscoveryOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DiscoveryOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DiscoveryOptions;
    static deserializeBinaryFromReader(message: DiscoveryOptions, reader: jspb.BinaryReader): DiscoveryOptions;
  }

  export namespace DiscoveryOptions {
    export type AsObject = {
      fdsMode: Settings.DiscoveryOptions.FdsModeMap[keyof Settings.DiscoveryOptions.FdsModeMap],
    }

    export interface FdsModeMap {
      BLACKLIST: 0;
      WHITELIST: 1;
      DISABLED: 2;
    }

    export const FdsMode: FdsModeMap;
  }

  export class ConsulConfiguration extends jspb.Message {
    getAddress(): string;
    setAddress(value: string): void;

    getHttpAddress(): string;
    setHttpAddress(value: string): void;

    getDnsAddress(): string;
    setDnsAddress(value: string): void;

    hasDnsPollingInterval(): boolean;
    clearDnsPollingInterval(): void;
    getDnsPollingInterval(): google_protobuf_duration_pb.Duration | undefined;
    setDnsPollingInterval(value?: google_protobuf_duration_pb.Duration): void;

    getDatacenter(): string;
    setDatacenter(value: string): void;

    getUsername(): string;
    setUsername(value: string): void;

    getPassword(): string;
    setPassword(value: string): void;

    getToken(): string;
    setToken(value: string): void;

    getCaFile(): string;
    setCaFile(value: string): void;

    getCaPath(): string;
    setCaPath(value: string): void;

    getCertFile(): string;
    setCertFile(value: string): void;

    getKeyFile(): string;
    setKeyFile(value: string): void;

    hasInsecureSkipVerify(): boolean;
    clearInsecureSkipVerify(): void;
    getInsecureSkipVerify(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setInsecureSkipVerify(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasWaitTime(): boolean;
    clearWaitTime(): void;
    getWaitTime(): google_protobuf_duration_pb.Duration | undefined;
    setWaitTime(value?: google_protobuf_duration_pb.Duration): void;

    hasServiceDiscovery(): boolean;
    clearServiceDiscovery(): void;
    getServiceDiscovery(): Settings.ConsulConfiguration.ServiceDiscoveryOptions | undefined;
    setServiceDiscovery(value?: Settings.ConsulConfiguration.ServiceDiscoveryOptions): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConsulConfiguration.AsObject;
    static toObject(includeInstance: boolean, msg: ConsulConfiguration): ConsulConfiguration.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConsulConfiguration, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConsulConfiguration;
    static deserializeBinaryFromReader(message: ConsulConfiguration, reader: jspb.BinaryReader): ConsulConfiguration;
  }

  export namespace ConsulConfiguration {
    export type AsObject = {
      address: string,
      httpAddress: string,
      dnsAddress: string,
      dnsPollingInterval?: google_protobuf_duration_pb.Duration.AsObject,
      datacenter: string,
      username: string,
      password: string,
      token: string,
      caFile: string,
      caPath: string,
      certFile: string,
      keyFile: string,
      insecureSkipVerify?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      waitTime?: google_protobuf_duration_pb.Duration.AsObject,
      serviceDiscovery?: Settings.ConsulConfiguration.ServiceDiscoveryOptions.AsObject,
    }

    export class ServiceDiscoveryOptions extends jspb.Message {
      clearDataCentersList(): void;
      getDataCentersList(): Array<string>;
      setDataCentersList(value: Array<string>): void;
      addDataCenters(value: string, index?: number): string;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): ServiceDiscoveryOptions.AsObject;
      static toObject(includeInstance: boolean, msg: ServiceDiscoveryOptions): ServiceDiscoveryOptions.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: ServiceDiscoveryOptions, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): ServiceDiscoveryOptions;
      static deserializeBinaryFromReader(message: ServiceDiscoveryOptions, reader: jspb.BinaryReader): ServiceDiscoveryOptions;
    }

    export namespace ServiceDiscoveryOptions {
      export type AsObject = {
        dataCentersList: Array<string>,
      }
    }
  }

  export class KubernetesConfiguration extends jspb.Message {
    hasRateLimits(): boolean;
    clearRateLimits(): void;
    getRateLimits(): Settings.KubernetesConfiguration.RateLimits | undefined;
    setRateLimits(value?: Settings.KubernetesConfiguration.RateLimits): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KubernetesConfiguration.AsObject;
    static toObject(includeInstance: boolean, msg: KubernetesConfiguration): KubernetesConfiguration.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KubernetesConfiguration, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KubernetesConfiguration;
    static deserializeBinaryFromReader(message: KubernetesConfiguration, reader: jspb.BinaryReader): KubernetesConfiguration;
  }

  export namespace KubernetesConfiguration {
    export type AsObject = {
      rateLimits?: Settings.KubernetesConfiguration.RateLimits.AsObject,
    }

    export class RateLimits extends jspb.Message {
      getQps(): number;
      setQps(value: number): void;

      getBurst(): number;
      setBurst(value: number): void;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): RateLimits.AsObject;
      static toObject(includeInstance: boolean, msg: RateLimits): RateLimits.AsObject;
      static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
      static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
      static serializeBinaryToWriter(message: RateLimits, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): RateLimits;
      static deserializeBinaryFromReader(message: RateLimits, reader: jspb.BinaryReader): RateLimits;
    }

    export namespace RateLimits {
      export type AsObject = {
        qps: number,
        burst: number,
      }
    }
  }

  export enum ConfigSourceCase {
    CONFIG_SOURCE_NOT_SET = 0,
    KUBERNETES_CONFIG_SOURCE = 4,
    DIRECTORY_CONFIG_SOURCE = 5,
    CONSUL_KV_SOURCE = 21,
  }

  export enum SecretSourceCase {
    SECRET_SOURCE_NOT_SET = 0,
    KUBERNETES_SECRET_SOURCE = 6,
    VAULT_SECRET_SOURCE = 7,
    DIRECTORY_SECRET_SOURCE = 8,
  }

  export enum ArtifactSourceCase {
    ARTIFACT_SOURCE_NOT_SET = 0,
    KUBERNETES_ARTIFACT_SOURCE = 9,
    DIRECTORY_ARTIFACT_SOURCE = 10,
    CONSUL_KV_ARTIFACT_SOURCE = 23,
  }
}

export class GlooOptions extends jspb.Message {
  getXdsBindAddr(): string;
  setXdsBindAddr(value: string): void;

  getValidationBindAddr(): string;
  setValidationBindAddr(value: string): void;

  hasCircuitBreakers(): boolean;
  clearCircuitBreakers(): void;
  getCircuitBreakers(): gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig | undefined;
  setCircuitBreakers(value?: gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig): void;

  hasEndpointsWarmingTimeout(): boolean;
  clearEndpointsWarmingTimeout(): void;
  getEndpointsWarmingTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setEndpointsWarmingTimeout(value?: google_protobuf_duration_pb.Duration): void;

  hasAwsOptions(): boolean;
  clearAwsOptions(): void;
  getAwsOptions(): GlooOptions.AWSOptions | undefined;
  setAwsOptions(value?: GlooOptions.AWSOptions): void;

  hasInvalidConfigPolicy(): boolean;
  clearInvalidConfigPolicy(): void;
  getInvalidConfigPolicy(): GlooOptions.InvalidConfigPolicy | undefined;
  setInvalidConfigPolicy(value?: GlooOptions.InvalidConfigPolicy): void;

  getDisableKubernetesDestinations(): boolean;
  setDisableKubernetesDestinations(value: boolean): void;

  hasDisableGrpcWeb(): boolean;
  clearDisableGrpcWeb(): void;
  getDisableGrpcWeb(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisableGrpcWeb(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasDisableProxyGarbageCollection(): boolean;
  clearDisableProxyGarbageCollection(): void;
  getDisableProxyGarbageCollection(): google_protobuf_wrappers_pb.BoolValue | undefined;
  setDisableProxyGarbageCollection(value?: google_protobuf_wrappers_pb.BoolValue): void;

  hasRegexMaxProgramSize(): boolean;
  clearRegexMaxProgramSize(): void;
  getRegexMaxProgramSize(): google_protobuf_wrappers_pb.UInt32Value | undefined;
  setRegexMaxProgramSize(value?: google_protobuf_wrappers_pb.UInt32Value): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GlooOptions.AsObject;
  static toObject(includeInstance: boolean, msg: GlooOptions): GlooOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GlooOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GlooOptions;
  static deserializeBinaryFromReader(message: GlooOptions, reader: jspb.BinaryReader): GlooOptions;
}

export namespace GlooOptions {
  export type AsObject = {
    xdsBindAddr: string,
    validationBindAddr: string,
    circuitBreakers?: gloo_projects_gloo_api_v1_circuit_breaker_pb.CircuitBreakerConfig.AsObject,
    endpointsWarmingTimeout?: google_protobuf_duration_pb.Duration.AsObject,
    awsOptions?: GlooOptions.AWSOptions.AsObject,
    invalidConfigPolicy?: GlooOptions.InvalidConfigPolicy.AsObject,
    disableKubernetesDestinations: boolean,
    disableGrpcWeb?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    disableProxyGarbageCollection?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    regexMaxProgramSize?: google_protobuf_wrappers_pb.UInt32Value.AsObject,
  }

  export class AWSOptions extends jspb.Message {
    getEnableCredentialsDiscovey(): boolean;
    setEnableCredentialsDiscovey(value: boolean): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AWSOptions.AsObject;
    static toObject(includeInstance: boolean, msg: AWSOptions): AWSOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AWSOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AWSOptions;
    static deserializeBinaryFromReader(message: AWSOptions, reader: jspb.BinaryReader): AWSOptions;
  }

  export namespace AWSOptions {
    export type AsObject = {
      enableCredentialsDiscovey: boolean,
    }
  }

  export class InvalidConfigPolicy extends jspb.Message {
    getReplaceInvalidRoutes(): boolean;
    setReplaceInvalidRoutes(value: boolean): void;

    getInvalidRouteResponseCode(): number;
    setInvalidRouteResponseCode(value: number): void;

    getInvalidRouteResponseBody(): string;
    setInvalidRouteResponseBody(value: string): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvalidConfigPolicy.AsObject;
    static toObject(includeInstance: boolean, msg: InvalidConfigPolicy): InvalidConfigPolicy.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvalidConfigPolicy, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvalidConfigPolicy;
    static deserializeBinaryFromReader(message: InvalidConfigPolicy, reader: jspb.BinaryReader): InvalidConfigPolicy;
  }

  export namespace InvalidConfigPolicy {
    export type AsObject = {
      replaceInvalidRoutes: boolean,
      invalidRouteResponseCode: number,
      invalidRouteResponseBody: string,
    }
  }
}

export class GatewayOptions extends jspb.Message {
  getValidationServerAddr(): string;
  setValidationServerAddr(value: string): void;

  hasValidation(): boolean;
  clearValidation(): void;
  getValidation(): GatewayOptions.ValidationOptions | undefined;
  setValidation(value?: GatewayOptions.ValidationOptions): void;

  getReadGatewaysFromAllNamespaces(): boolean;
  setReadGatewaysFromAllNamespaces(value: boolean): void;

  getAlwaysSortRouteTableRoutes(): boolean;
  setAlwaysSortRouteTableRoutes(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GatewayOptions.AsObject;
  static toObject(includeInstance: boolean, msg: GatewayOptions): GatewayOptions.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GatewayOptions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GatewayOptions;
  static deserializeBinaryFromReader(message: GatewayOptions, reader: jspb.BinaryReader): GatewayOptions;
}

export namespace GatewayOptions {
  export type AsObject = {
    validationServerAddr: string,
    validation?: GatewayOptions.ValidationOptions.AsObject,
    readGatewaysFromAllNamespaces: boolean,
    alwaysSortRouteTableRoutes: boolean,
  }

  export class ValidationOptions extends jspb.Message {
    getProxyValidationServerAddr(): string;
    setProxyValidationServerAddr(value: string): void;

    getValidationWebhookTlsCert(): string;
    setValidationWebhookTlsCert(value: string): void;

    getValidationWebhookTlsKey(): string;
    setValidationWebhookTlsKey(value: string): void;

    getIgnoreGlooValidationFailure(): boolean;
    setIgnoreGlooValidationFailure(value: boolean): void;

    hasAlwaysAccept(): boolean;
    clearAlwaysAccept(): void;
    getAlwaysAccept(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setAlwaysAccept(value?: google_protobuf_wrappers_pb.BoolValue): void;

    hasAllowWarnings(): boolean;
    clearAllowWarnings(): void;
    getAllowWarnings(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setAllowWarnings(value?: google_protobuf_wrappers_pb.BoolValue): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ValidationOptions.AsObject;
    static toObject(includeInstance: boolean, msg: ValidationOptions): ValidationOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ValidationOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ValidationOptions;
    static deserializeBinaryFromReader(message: ValidationOptions, reader: jspb.BinaryReader): ValidationOptions;
  }

  export namespace ValidationOptions {
    export type AsObject = {
      proxyValidationServerAddr: string,
      validationWebhookTlsCert: string,
      validationWebhookTlsKey: string,
      ignoreGlooValidationFailure: boolean,
      alwaysAccept?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      allowWarnings?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }
}
