// package: glooe.solo.io
// file: github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/ratelimit.proto

var github_com_solo_io_solo_apis_api_gloo_gloo_v1_enterprise_ratelimit_pb = require("../../../../../../../../github.com/solo-io/solo-apis/api/gloo/gloo/v1/enterprise/ratelimit_pb");
var github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb = require("../../../../../../../../github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var RateLimitDiscoveryService = (function () {
  function RateLimitDiscoveryService() {}
  RateLimitDiscoveryService.serviceName = "glooe.solo.io.RateLimitDiscoveryService";
  return RateLimitDiscoveryService;
}());

RateLimitDiscoveryService.StreamRateLimitConfig = {
  methodName: "StreamRateLimitConfig",
  service: RateLimitDiscoveryService,
  requestStream: true,
  responseStream: true,
  requestType: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest,
  responseType: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse
};

RateLimitDiscoveryService.DeltaRateLimitConfig = {
  methodName: "DeltaRateLimitConfig",
  service: RateLimitDiscoveryService,
  requestStream: true,
  responseStream: true,
  requestType: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryRequest,
  responseType: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DeltaDiscoveryResponse
};

RateLimitDiscoveryService.FetchRateLimitConfig = {
  methodName: "FetchRateLimitConfig",
  service: RateLimitDiscoveryService,
  requestStream: false,
  responseStream: false,
  requestType: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryRequest,
  responseType: github_com_solo_io_solo_kit_api_external_envoy_api_v2_discovery_pb.DiscoveryResponse
};

exports.RateLimitDiscoveryService = RateLimitDiscoveryService;

function RateLimitDiscoveryServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

RateLimitDiscoveryServiceClient.prototype.streamRateLimitConfig = function streamRateLimitConfig(metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.client(RateLimitDiscoveryService.StreamRateLimitConfig, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  client.onMessage(function (message) {
    listeners.data.forEach(function (handler) {
      handler(message);
    })
  });
  client.start(metadata);
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

RateLimitDiscoveryServiceClient.prototype.deltaRateLimitConfig = function deltaRateLimitConfig(metadata) {
  var listeners = {
    data: [],
    end: [],
    status: []
  };
  var client = grpc.client(RateLimitDiscoveryService.DeltaRateLimitConfig, {
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport
  });
  client.onEnd(function (status, statusMessage, trailers) {
    listeners.status.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners.end.forEach(function (handler) {
      handler({ code: status, details: statusMessage, metadata: trailers });
    });
    listeners = null;
  });
  client.onMessage(function (message) {
    listeners.data.forEach(function (handler) {
      handler(message);
    })
  });
  client.start(metadata);
  return {
    on: function (type, handler) {
      listeners[type].push(handler);
      return this;
    },
    write: function (requestMessage) {
      client.send(requestMessage);
      return this;
    },
    end: function () {
      client.finishSend();
    },
    cancel: function () {
      listeners = null;
      client.close();
    }
  };
};

RateLimitDiscoveryServiceClient.prototype.fetchRateLimitConfig = function fetchRateLimitConfig(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(RateLimitDiscoveryService.FetchRateLimitConfig, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.RateLimitDiscoveryServiceClient = RateLimitDiscoveryServiceClient;

