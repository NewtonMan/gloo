package runner

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/solo-io/go-utils/healthchecker"
	"github.com/solo-io/go-utils/loggingutils"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	_struct "github.com/golang/protobuf/ptypes/struct"
	extauthconfig "github.com/solo-io/ext-auth-service/pkg/config"
	extauth "github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/go-utils/contextutils"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	xdsproto "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/config"

	"github.com/solo-io/go-utils/stats"

	"google.golang.org/grpc"
)

// The extauth server sends xDS discovery requests to Gloo to get its configuration from Gloo. This constant determines
// the value of the nodeInfo.Metadata.role field that the server sends along to retrieve its configuration snapshot,
// similarly to how the regular Gloo gateway-proxies do.
const ExtAuthServerRole = "extauth"

func init() {
	_ = view.Register(ocgrpc.DefaultServerViews...)
}

func Run() {
	clientSettings := NewSettings()
	ctx := context.Background()

	if clientSettings.DebugPort != 0 {
		// TODO(yuval-k): we need to start the stats server before calling contextutils
		// need to think of a better way to express this dependency, or preferably, fix it.
		stats.StartStatsServerWithPort(stats.StartupOptions{Port: clientSettings.DebugPort})
	}

	err := RunWithSettings(ctx, clientSettings)

	if err != nil {
		if ctx.Err() == nil {
			// not a context error - panic
			panic(err)
		}
	}
}

func RunWithSettings(ctx context.Context, clientSettings Settings) error {
	ctx = contextutils.WithLogger(ctx, "extauth")

	err := StartExtAuth(ctx, clientSettings, extauth.NewServer())
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

func StartExtAuth(ctx context.Context, clientSettings Settings, service *extauth.Server) error {
	logger := contextutils.LoggerFrom(ctx)
	callerCtx, cancel := context.WithCancel(ctx) // do not use callerCtx anywhere outside of the interceptors

	madeHealthCheckFail := make(chan struct{}, 1)
	srv := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}), grpc.ChainUnaryInterceptor(
		loggingutils.GrpcUnaryServerLoggerInterceptor(logger),
		healthchecker.GrpcUnaryServerHealthCheckerInterceptor(callerCtx, madeHealthCheckFail)))

	pb.RegisterAuthorizationServer(srv, service)
	hc := healthchecker.NewGrpc(clientSettings.ServiceName, health.NewServer())
	healthpb.RegisterHealthServer(srv, hc.GetServer())
	reflection.Register(srv)

	logger.Infow("Starting ext-auth server")

	err := StartExtAuthWithGrpcServer(ctx, clientSettings, service)
	if err != nil {
		logger.Error("Failed to start ext-auth server: %v", err)
		return err
	}

	var addr, runMode, network string
	if clientSettings.ServerUDSAddr != "" {
		addr = clientSettings.ServerUDSAddr
		runMode = "unixDomainSocket"
		network = "unix"
	} else {
		addr = fmt.Sprintf(":%d", clientSettings.ServerPort)
		runMode = "gRPC"
		network = "tcp"
	}

	logger.Infof("extauth server running in [%s] mode, listening at [%s]", runMode, addr)
	lis, err := net.Listen(network, addr)
	if err != nil {
		logger.Errorw("Failed to announce on network", zap.Any("mode", runMode), zap.Any("address", addr), zap.Error(err))
		return err
	}
	terminationSigs := make(chan os.Signal, 1)
	signal.Notify(terminationSigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-ctx.Done():
			logger.Infof("Context has been cancelled, shutting down.")
		case <-terminationSigs:
			cancel()
			<-callerCtx.Done()
			logger.Infof("Termination signal received, shutting down.")
		}
		select {
		case <-madeHealthCheckFail:
			logger.Infof("Health check will now fail in envoy.")
		case <-time.NewTimer(15 * time.Second).C: // timeout must be > the extauth health check interval
			logger.Infof("Unable to make health check fail in envoy. Shutting down anyway.")
		}
		srv.GracefulStop()
		_ = lis.Close()
	}()

	return srv.Serve(lis)
}

func StartExtAuthWithGrpcServer(ctx context.Context, clientSettings Settings, service extauthconfig.ConfigMutator) error {
	var nodeInfo core.Node
	var err error
	nodeInfo.Id, err = os.Hostname()
	// TODO(yuval-k): unhardcode this
	if err != nil {
		nodeInfo.Id = "extauth-unknown"
	}
	nodeInfo.Cluster = "extauth"
	role := ExtAuthServerRole
	nodeInfo.Metadata = &_struct.Struct{
		Fields: map[string]*_struct.Value{
			"role": {
				Kind: &_struct.Value_StringValue{
					StringValue: role,
				},
			},
		},
	}

	go clientLoop(ctx, clientSettings, nodeInfo, service)
	return nil
}

func clientLoop(ctx context.Context, clientSettings Settings, nodeInfo core.Node, service extauthconfig.ConfigMutator) {

	generator := config.NewGenerator(
		ctx,
		[]byte(clientSettings.SigningKey),
		clientSettings.UserIdHeader,
		plugins.NewPluginLoader(clientSettings.PluginDirectory),
	)

	_ = contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{}).Backoff(ctx, func(ctx context.Context) error {

		client := xdsproto.NewExtAuthConfigClient(
			&nodeInfo,
			func(version string, resources []*xdsproto.ExtAuthConfig) error {

				logger := contextutils.LoggerFrom(ctx)
				logger.Infow("got new config", zap.Any("config", resources))

				serverState, err := generator.GenerateConfig(resources)
				if err != nil {
					logger.Errorw("failed to generate config", zap.Any("err", err))
					return err
				}
				service.UpdateConfig(serverState)
				return nil
			},
		)

		// We are using non secure gRPC to Gloo with the assumption that it will be
		// secured by envoy. if this assumption is not correct this needs to change.
		conn, err := grpc.DialContext(ctx, clientSettings.GlooAddress, grpc.WithInsecure())
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorw("failed to create gRPC client connection to Gloo", zap.Any("error", err))
			return err
		}
		// TODO(yuval-k): a stat that indicates we are connected, with the reverse one deferred.
		// TODO(yuval-k): write a warning log
		err = client.Start(ctx, conn)
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorw("failed to start xDS client", zap.Any("error", err))
		} else {
			contextutils.LoggerFrom(ctx).Info("successfully started xDS client")
		}
		return err
	})
}
