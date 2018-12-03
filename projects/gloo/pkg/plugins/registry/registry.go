package registry

import (
	"github.com/solo-io/solo-projects/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/aws"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/azure"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/consul"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/faultinjection"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/grpc"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/prefixrewrite"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/rest"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/static"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/transformation"
)

type registry struct {
	plugins []plugins.Plugin
}

var globalRegistry = func(opts bootstrap.Opts) *registry {
	transformationPlugin := transformation.NewPlugin()
	reg := &registry{}
	// plugins should be added here
	reg.plugins = append(reg.plugins,
		azure.NewPlugin(&transformationPlugin.RequireTransformationFilter),
		aws.NewPlugin(&transformationPlugin.RequireTransformationFilter),
		rest.NewPlugin(&transformationPlugin.RequireTransformationFilter),
		ratelimit.NewPlugin(),
		static.NewPlugin(),
		transformationPlugin,
		consul.NewPlugin(),
		grpc.NewPlugin(&transformationPlugin.RequireTransformationFilter),
		faultinjection.NewPlugin(),
		prefixrewrite.NewPlugin(),
	)
	if opts.KubeClient != nil {
		reg.plugins = append(reg.plugins, kubernetes.NewPlugin(opts.KubeClient))
	}

	return reg
}

func Plugins(opts bootstrap.Opts) []plugins.Plugin {
	return globalRegistry(opts).plugins
}
