package create

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/argsutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/constants"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/surveyutils"
	gloov1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/spf13/cobra"
)

func virtualServiceCreate(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "virtualservice",
		Aliases: []string{"vs", "virtualservice", "virtualservices"},
		Short:   "Create a Virtual Service",
		Long: "A virtual service describes the set of routes to match for a set of domains. \n" +
			"Virtual services are containers for routes assigned to a domain or set of domains. \n" +
			"Virtual services must not have overlapping domains, as the virtual service to match a request " +
			"is selected by the Host header (in HTTP1) or :authority header (in HTTP2). " +
			"The routes within a virtual service ",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Top.Interactive {
				if err := surveyutils.AddVirtualServiceFlagsInteractive(&opts.Create.VirtualService); err != nil {
					return err
				}
			}
			err := argsutils.MetadataArgsParse(opts, args)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return createVirtualService(opts, args)
		},
	}
	pflags := cmd.PersistentFlags()
	flagutils.AddMetadataFlags(pflags, &opts.Metadata)
	flagutils.AddVirtualServiceFlags(pflags, &opts.Create.VirtualService)
	return cmd
}

func createVirtualService(opts *options.Options, args []string) error {
	vs, err := virtualServiceFromOpts(opts.Metadata, opts.Create.VirtualService)
	if err != nil {
		return err
	}
	virtualServiceClient := helpers.MustVirtualServiceClient()
	vs, err = virtualServiceClient.Write(vs, clients.WriteOpts{})
	if err != nil {
		return err
	}

	helpers.PrintVirtualServices(v1.VirtualServiceList{vs}, opts.Top.Output)

	return nil
}

func virtualServiceFromOpts(meta core.Metadata, input options.InputVirtualService) (*v1.VirtualService, error) {
	if len(input.Domains) == 0 {
		input.Domains = constants.DefaultDomains
	}
	vs := &v1.VirtualService{
		Metadata: meta,
		VirtualHost: &gloov1.VirtualHost{
			Domains: input.Domains,
		},
	}
	rateLimitOpts := input.RateLimit
	if rateLimitOpts.Enable {
		if vs.VirtualHost.VirtualHostPlugins == nil {
			vs.VirtualHost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{}
		}
		timeUnit, ok := ratelimit.RateLimit_Unit_value[rateLimitOpts.TimeUnit]
		if !ok {
			return nil, errors.Errorf("invalid time unit specified: %v", rateLimitOpts.TimeUnit)
		}
		vs.VirtualHost.VirtualHostPlugins.RateLimits = &ratelimit.IngressRateLimit{
			AnonymousLimits: &ratelimit.RateLimit{
				Unit:            ratelimit.RateLimit_Unit(timeUnit),
				RequestsPerUnit: input.RateLimit.RequestsPerTimeUnit,
			},
		}
	}

	return vs, nil
}
