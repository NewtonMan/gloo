import { grpc } from '@improbable-eng/grpc-web';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';
import {
  GetIsLicenseValidRequest,
  GetIsLicenseValidResponse,
  GetOAuthEndpointRequest,
  GetOAuthEndpointResponse,
  GetPodNamespaceRequest,
  GetPodNamespaceResponse,
  GetSettingsRequest,
  GetSettingsResponse,
  GetVersionRequest,
  GetVersionResponse,
  ListNamespacesRequest,
  ListNamespacesResponse,
  UpdateSettingsRequest,
  UpdateSettingsResponse,
  SettingsDetails,
  UpdateSettingsYamlRequest
} from 'proto/solo-projects/projects/grpcserver/api/v1/config_pb';
import {
  ConfigApiClient,
  ConfigApi
} from 'proto/solo-projects/projects/grpcserver/api/v1/config_pb_service';
import { host } from 'store';
import { Settings } from 'proto/gloo/projects/gloo/api/v1/settings_pb';
import { ResourceRef } from 'proto/solo-kit/api/v1/ref_pb';
import { EditedResourceYaml } from 'proto/solo-projects/projects/grpcserver/api/v1/types_pb';
import { guardByLicense } from './actions';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';

const client = new ConfigApiClient(host, {
  transport: grpc.CrossBrowserHttpTransport({ withCredentials: false }),
  debug: true
});

export const configAPI = {
  getVersion,
  getSettings,
  getOAuthEndpoint,
  getSettingsGrpc,
  updateSettings,
  updateSettingsYaml,
  getIsLicenseValid,
  listNamespaces,
  getPodNamespace,
  updateRefreshRate,
  updateWatchNamespaces
};

function getVersion(): Promise<string> {
  return new Promise((resolve, reject) => {
    let request = new GetVersionRequest();
    client.getVersion(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().version);
      }
    });
  });
}

function getOAuthEndpoint(): Promise<GetOAuthEndpointResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetOAuthEndpointRequest();
    client.getOAuthEndpoint(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function getSettingsGrpc(): Promise<GetSettingsResponse> {
  return new Promise((resolve, reject) => {
    let request = new GetSettingsRequest();
    client.getSettings(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!);
      }
    });
  });
}

function getSettings(): Promise<SettingsDetails.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetSettingsRequest();
    client.getSettings(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().settingsDetails);
      }
    });
  });
}

function updateWatchNamespaces(updateWatchNamespacesRequest: {
  watchNamespacesList: string[];
}): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let currentSettingsReq = await configAPI.getSettingsGrpc();
    let settingsToUpdate = currentSettingsReq
      .getSettingsDetails()
      ?.getSettings();
    let request = new UpdateSettingsRequest();
    let { watchNamespacesList } = updateWatchNamespacesRequest;
    if (settingsToUpdate !== undefined) {
      settingsToUpdate.setWatchNamespacesList(watchNamespacesList);

      request.setSettings(settingsToUpdate);
    }
    client.updateSettings(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function updateRefreshRate(updateRefreshRateRequest: {
  refreshRate: Duration.AsObject;
}): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    let currentSettingsReq = await configAPI.getSettingsGrpc();
    let settingsToUpdate = currentSettingsReq
      .getSettingsDetails()
      ?.getSettings();
    let newRefreshRate = new Duration();
    newRefreshRate.setNanos(updateRefreshRateRequest.refreshRate.nanos);
    newRefreshRate.setSeconds(updateRefreshRateRequest.refreshRate.seconds);

    let updateSettingsRequest = new UpdateSettingsRequest();
    if (settingsToUpdate !== undefined) {
      settingsToUpdate.setRefreshRate(newRefreshRate);
      updateSettingsRequest.setSettings(settingsToUpdate);
    }

    client.updateSettings(updateSettingsRequest, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function updateSettings(
  updateSettingsRequest: UpdateSettingsRequest
): Promise<UpdateSettingsResponse.AsObject> {
  return new Promise(async (resolve, reject) => {
    guardByLicense();

    client.updateSettings(updateSettingsRequest, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function updateSettingsYaml(
  updateSettingsYamlRequest: UpdateSettingsYamlRequest.AsObject
): Promise<SettingsDetails.AsObject> {
  return new Promise(async (resolve, reject) => {
    guardByLicense();

    let request = new UpdateSettingsYamlRequest();
    let settingsRef = new ResourceRef();

    let editedYamlData = new EditedResourceYaml();
    let { ref, editedYaml } = updateSettingsYamlRequest.editedYamlData!;
    settingsRef.setName(ref!.name);
    settingsRef.setNamespace(ref!.namespace);

    editedYamlData.setRef(settingsRef);
    editedYamlData.setEditedYaml(editedYaml);
    request.setEditedYamlData(editedYamlData);

    client.updateSettingsYaml(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().settingsDetails);
      }
    });
  });
}

function getIsLicenseValid(): Promise<GetIsLicenseValidResponse.AsObject> {
  return new Promise((resolve, reject) => {
    let request = new GetIsLicenseValidRequest();
    client.getIsLicenseValid(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject());
      }
    });
  });
}

function listNamespaces(): Promise<string[]> {
  return new Promise((resolve, reject) => {
    let request = new ListNamespacesRequest();
    client.listNamespaces(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().namespacesList);
      }
    });
  });
}

function getPodNamespace(): Promise<string> {
  return new Promise((resolve, reject) => {
    let request = new GetPodNamespaceRequest();
    client.getPodNamespace(request, (error, data) => {
      if (error !== null) {
        console.error('Error:', error.message);
        console.error('Code:', error.code);
        console.error('Metadata:', error.metadata);
        reject(error);
      } else {
        resolve(data!.toObject().namespace);
      }
    });
  });
}
