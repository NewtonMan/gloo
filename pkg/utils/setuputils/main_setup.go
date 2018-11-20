package setuputils

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
)

func Main(loggingPrefix string, setupSyncer v1.SetupSyncer, settingsDir string) error {
	settingsClient, err := KubeOrFileSettingsClient(settingsDir)
	if err != nil {
		return err
	}
	if err := settingsClient.Register(); err != nil {
		return err
	}
	if err := writeDefaultSettings(settingsClient); err != nil {
		return err
	}

	emitter := v1.NewSetupEmitter(settingsClient)
	ctx := contextutils.WithLogger(context.Background(), loggingPrefix)
	eventLoop := v1.NewSetupEventLoop(emitter, setupSyncer)
	errs, err := eventLoop.Run([]string{defaults.GlooSystem}, clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Second,
	})
	if err != nil {
		return err
	}
	for err := range errs {
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}
	return nil
}

// TODO (ilackarms): instead of using an heuristic here, read from a CLI flagg
// first attempt to use kube crd, otherwise fall back to file
func KubeOrFileSettingsClient(settingsDir string) (v1.SettingsClient, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err == nil {
		return v1.NewSettingsClient(&factory.KubeResourceClientFactory{
			Crd:         v1.SettingsCrd,
			Cfg:         cfg,
			SharedCache: kube.NewKubeCache(),
		})
	}
	return v1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: settingsDir,
	})
}

// TODO(ilackarms): remove this or move it to a test package, only use settings watch for production gloo
func writeDefaultSettings(cli v1.SettingsClient) error {
	settings := &v1.Settings{
		ConfigSource: &v1.Settings_KubernetesConfigSource{
			KubernetesConfigSource: &v1.Settings_KubernetesCrds{},
		},
		ArtifactSource: &v1.Settings_KubernetesArtifactSource{
			KubernetesArtifactSource: &v1.Settings_KubernetesConfigmaps{},
		},
		SecretSource: &v1.Settings_KubernetesSecretSource{
			KubernetesSecretSource: &v1.Settings_KubernetesSecrets{},
		},
		BindAddr:        "0.0.0.0:9977",
		RefreshRate:     types.DurationProto(time.Minute),
		DevMode:         true,
		WatchNamespaces: []string{"default", defaults.GlooSystem},
		Metadata:        core.Metadata{Namespace: defaults.GlooSystem, Name: "gloo"},
	}
	if _, err := cli.Write(settings, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
		return err
	}
	return nil
}
