package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fault "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/faultinjection"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gloov1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Fault Injection", func() {

	var (
		testClients services.TestClients
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx, _ = context.WithCancel(context.Background())
		t := services.RunGateway(ctx, true)
		testClients = t
	})

	Context("with envoy", func() {

		var (
			envoyInstance *services.EnvoyInstance
			up            *gloov1.Upstream
			opts          clients.WriteOpts
		)

		setupProxy := func(proxy *gloov1.Proxy, up *gloov1.Upstream) {
			proxyCli := testClients.ProxyClient
			_, err := proxyCli.Write(proxy, opts)
			Expect(err).NotTo(HaveOccurred())
		}
		envoyPort := uint32(8080)
		setupInitialProxy := func() {
			proxy := getGlooProxy(nil, nil, envoyPort, up, "")
			setupProxy(proxy, up)
			Eventually(func() error {
				_, err := http.Get(fmt.Sprintf("http://%s:%d/status/200", "localhost", envoyPort))
				if err != nil {
					return err
				}
				return nil
			}, "5s", ".1s").Should(BeNil())
		}

		setupUpstream := func() {
			tu := v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
			// drain channel as we dont care about it
			go func() {
				for range tu.C {
				}
			}()
			var opts clients.WriteOpts
			up = tu.Upstream
			_, err := testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())
		}

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			setupUpstream()
			setupInitialProxy()
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		It("should cause envoy abort fault", func() {
			abort := &fault.RouteAbort{
				HttpStatus: uint32(503),
				Percentage: float32(100),
			}
			proxy := getGlooProxy(abort, nil, envoyPort, up, "3")
			opts.OverwriteExisting = true
			setupProxy(proxy, up)

			Eventually(func() error {
				res, err := http.Get(fmt.Sprintf("http://%s:%d/status/200", "localhost", envoyPort))
				if err != nil {
					return err
				}
				if res.StatusCode != http.StatusServiceUnavailable {
					return errors.New(fmt.Sprintf("%v is not ServiceUnavailable", res.StatusCode))
				}
				return nil
			}, "5s", ".1s").Should(BeNil())
		})

		It("should cause envoy delay fault", func() {
			fixedDelay := time.Second * 3
			delay := &fault.RouteDelay{
				FixedDelay: &fixedDelay,
				Percentage: float32(100),
			}
			proxy := getGlooProxy(nil, delay, envoyPort, up, "3")
			opts.OverwriteExisting = true
			setupProxy(proxy, up)

			Eventually(func() error {
				start := time.Now()
				_, err := http.Get(fmt.Sprintf("http://%s:%d/status/200", "localhost", envoyPort))
				if err != nil {
					return err
				}
				elapsed := time.Now().Sub(start)
				if elapsed < (3 * time.Second) {
					return errors.New(fmt.Sprintf("Elapsed time %s not longer than delay %s", elapsed.String(), fixedDelay.String()))
				}
				return nil
			}, "5s", ".1s").Should(BeNil())

		})
	})
})

func getGlooProxy(abort *fault.RouteAbort, delay *fault.RouteDelay, envoyPort uint32, up *gloov1.Upstream, resourceVersion string) *gloov1.Proxy {
	return &gloov1.Proxy{
		Metadata: core.Metadata{
			Name:            "proxy",
			Namespace:       "default",
			ResourceVersion: resourceVersion,
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "127.0.0.1",
			BindPort:    envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{{
							Matcher: &gloov1.Matcher{
								PathSpecifier: &gloov1.Matcher_Prefix{
									Prefix: "/",
								},
							},
							Action: &gloov1.Route_RouteAction{
								RouteAction: &gloov1.RouteAction{
									Destination: &gloov1.RouteAction_Single{
										Single: &gloov1.Destination{
											Upstream: up.Metadata.Ref(),
										},
									},
								},
							},
							RoutePlugins: &gloov1.RoutePlugins{
								Faults: &fault.RouteFaults{
									Abort: abort,
									Delay: delay,
								},
							},
						}},
						VirtualHostPlugins: &gloov1.VirtualHostPlugins{},
					}},
				},
			},
		}},
	}
}
