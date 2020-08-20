package configsvc_test

import (
	"context"
	"time"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"

	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	mock_gloo "github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	mock_license "github.com/solo-io/solo-projects/pkg/license/mocks"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_namespace "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/configsvc"
)

var (
	apiserver         v1.ConfigApiServer
	mockCtrl          *gomock.Controller
	settingsClient    *mock_gloo.MockSettingsClient
	licenseClient     *mock_license.MockClient
	namespaceClient   *mock_namespace.MockNamespaceClient
	clientCache       *mocks.MockClientCache
	rawGetter         *mock_rawgetter.MockRawGetter
	rbacNamespaced    setup.RbacNamespaced
	podNamespace      = "pod-ns"
	testVersion       = setup.BuildVersion("test-version")
	testOAuthEndpoint = v1.OAuthEndpoint{Url: "test", ClientName: "name"}
	testErr           = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {
	getDetails := func(settings *gloov1.Settings, raw *v1.Raw) *v1.SettingsDetails {
		return &v1.SettingsDetails{
			Settings: settings,
			Raw:      raw,
		}
	}

	getRaw := func(name string) *v1.Raw {
		return &v1.Raw{
			FileName: name + ".yaml",
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		settingsClient = mock_gloo.NewMockSettingsClient(mockCtrl)
		licenseClient = mock_license.NewMockClient(mockCtrl)
		namespaceClient = mock_namespace.NewMockNamespaceClient(mockCtrl)
		clientCache = mocks.NewMockClientCache(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		rbacNamespaced = false
	})

	JustBeforeEach(func() {
		server, err := configsvc.NewConfigGrpcService(context.TODO(), clientCache, licenseClient, namespaceClient, rawGetter, testOAuthEndpoint, testVersion, podNamespace, rbacNamespaced)
		Expect(err).NotTo(HaveOccurred())
		apiserver = server
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetVersion", func() {
		It("works", func() {
			actual, err := apiserver.GetVersion(context.TODO(), &v1.GetVersionRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetVersionResponse{Version: string(testVersion)}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

	Describe("GetOAuthEndpoint", func() {
		It("works", func() {
			actual, err := apiserver.GetOAuthEndpoint(context.TODO(), &v1.GetOAuthEndpointRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetOAuthEndpointResponse{OAuthEndpoint: &testOAuthEndpoint}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

	Describe("GetIsLicenseValid", func() {
		It("works when the license client works", func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)

			actual, err := apiserver.GetIsLicenseValid(context.TODO(), &v1.GetIsLicenseValidRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetIsLicenseValidResponse{IsLicenseValid: true}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("returns reason when license client errors", func() {
			licenseClient.EXPECT().IsLicenseValid().Return(testErr)

			resp, err := apiserver.GetIsLicenseValid(context.TODO(), &v1.GetIsLicenseValidRequest{})
			Expect(err).NotTo(HaveOccurred())
			expectedErr := configsvc.LicenseIsInvalidError(testErr)
			Expect(resp.InvalidReason).To(Equal(expectedErr.Error()))
		})
	})

	Describe("GetSettings", func() {
		It("works when the settings client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			settings := &gloov1.Settings{RefreshRate: &types.Duration{Seconds: 1}, Metadata: metadata}
			raw := getRaw(metadata.Name)
			settingsClient.EXPECT().
				Read(podNamespace, defaults.SettingsName, clients.ReadOpts{Ctx: context.TODO()}).
				Return(settings, nil)
			clientCache.EXPECT().GetSettingsClient().Return(settingsClient)
			rawGetter.EXPECT().
				GetRaw(context.Background(), settings, gloov1.SettingsCrd).
				Return(raw)
			actual, err := apiserver.GetSettings(context.TODO(), &v1.GetSettingsRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetSettingsResponse{SettingsDetails: getDetails(settings, raw)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the settings client errors", func() {
			settingsClient.EXPECT().
				Read(podNamespace, defaults.SettingsName, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)
			clientCache.EXPECT().GetSettingsClient().Return(settingsClient)

			_, err := apiserver.GetSettings(context.TODO(), &v1.GetSettingsRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToReadSettingsError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("UpdateSettings", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})
		Context("with unified input objects", func() {
			buildSettings := func(watchNamespaces []string, refreshRate *types.Duration) *gloov1.Settings {
				namespaces := watchNamespaces

				// the server will change an empty array in the proto object to nil
				if len(watchNamespaces) == 0 {
					namespaces = nil
				}
				return &gloov1.Settings{
					Metadata: core.Metadata{
						Name:      "name",
						Namespace: "ns",
					},
					RefreshRate:     refreshRate,
					WatchNamespaces: namespaces,
				}
			}
			It("works", func() {
				metadata := core.Metadata{
					Namespace: "ns",
					Name:      "name",
				}
				raw := getRaw(metadata.Name)
				refreshRate := types.Duration{Seconds: 1}
				watchNamespaces := []string{"a", "b"}
				request := &v1.UpdateSettingsRequest{
					Settings: buildSettings(watchNamespaces, &refreshRate),
				}
				writeSettings := buildSettings(watchNamespaces, &refreshRate)

				settingsClient.EXPECT().
					Write(writeSettings, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
					Return(writeSettings, nil)
				clientCache.EXPECT().GetSettingsClient().Return(settingsClient)
				rawGetter.EXPECT().
					GetRaw(context.Background(), writeSettings, gloov1.SettingsCrd).
					Return(raw)

				actual, err := apiserver.UpdateSettings(context.TODO(), request)
				Expect(err).NotTo(HaveOccurred())
				expected := &v1.UpdateSettingsResponse{SettingsDetails: getDetails(writeSettings, raw)}
				ExpectEqualProtoMessages(actual, expected)
			})

			It("errors when the provided refresh rate is invalid", func() {
				refreshRate := types.Duration{Nanos: 1}
				request := &v1.UpdateSettingsRequest{
					Settings: buildSettings([]string{}, &refreshRate),
				}

				_, err := apiserver.UpdateSettings(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := configsvc.InvalidRefreshRateError(1 * time.Nanosecond)
				wrapped := configsvc.FailedToUpdateSettingsError(expectedErr)
				Expect(err.Error()).To(ContainSubstring(wrapped.Error()))
			})

			It("errors when the settings client fails to write", func() {
				settings := buildSettings([]string{}, &types.Duration{Seconds: 1})
				request := &v1.UpdateSettingsRequest{Settings: settings}

				settingsClient.EXPECT().
					Write(settings, clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
					Return(nil, testErr)
				clientCache.EXPECT().GetSettingsClient().Return(settingsClient)

				_, err := apiserver.UpdateSettings(context.TODO(), request)
				Expect(err).To(HaveOccurred())
				expectedErr := configsvc.FailedToUpdateSettingsError(testErr)
				Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
			})
		})
	})

	Describe("UpdateSettingsYaml", func() {
		BeforeEach(func() {
			licenseClient.EXPECT().IsLicenseValid().Return(nil)
		})

		It("works on valid input", func() {
			yamlString := "totally-valid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			settings := &gloov1.Settings{
				Metadata: metadata,
			}
			request := &v1.UpdateSettingsYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(nil)
			settingsClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(settings, nil)
			rawGetter.EXPECT().
				GetRaw(context.TODO(), settings, gloov1.SettingsCrd).
				Return(getRaw(metadata.Name))

			clientCache.EXPECT().GetSettingsClient().Return(settingsClient)
			actual, err := apiserver.UpdateSettingsYaml(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			settingsDetails := &v1.SettingsDetails{
				Settings: settings,
				Raw:      getRaw(metadata.Name),
			}
			expected := &v1.UpdateSettingsResponse{
				SettingsDetails: settingsDetails,
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors with an invalid yaml", func() {
			yamlString := "totally-invalid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			request := &v1.UpdateSettingsYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(testErr)

			_, err := apiserver.UpdateSettingsYaml(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToParseSettingsFromYamlError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		It("errors when the settings client errors", func() {
			yamlString := "totally-valid-yaml"
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			request := &v1.UpdateSettingsYamlRequest{
				EditedYamlData: &v1.EditedResourceYaml{
					Ref:        &ref,
					EditedYaml: yamlString,
				},
			}

			rawGetter.EXPECT().
				InitResourceFromYamlString(context.TODO(), yamlString, &ref, gomock.Any()).
				Return(nil)
			settingsClient.EXPECT().
				Write(gomock.Any(), clients.WriteOpts{Ctx: context.TODO(), OverwriteExisting: true}).
				Return(nil, testErr)

			clientCache.EXPECT().GetSettingsClient().Return(settingsClient)
			_, err := apiserver.UpdateSettingsYaml(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToUpdateSettingsError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListNamespaces", func() {
		It("works when the namespace client works", func() {
			namespaceList := []string{"a", "b", "c"}

			namespaceClient.EXPECT().
				ListNamespaces().
				Return(namespaceList, nil)

			actual, err := apiserver.ListNamespaces(context.TODO(), &v1.ListNamespacesRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListNamespacesResponse{Namespaces: namespaceList}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the namespace client errors", func() {
			namespaceClient.EXPECT().
				ListNamespaces().
				Return(nil, testErr)

			_, err := apiserver.ListNamespaces(context.TODO(), &v1.ListNamespacesRequest{})
			Expect(err).To(HaveOccurred())
			expectedErr := configsvc.FailedToListNamespacesError(testErr)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})

		When("installation is namespace-scoped", func() {
			BeforeEach(func() {
				rbacNamespaced = true
			})

			It("does not attempt to list all namespaced and just returns the pod namespace", func() {
				expected := &v1.ListNamespacesResponse{Namespaces: []string{podNamespace}}
				actual, err := apiserver.ListNamespaces(context.TODO(), &v1.ListNamespacesRequest{})
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).To(Equal(expected))
			})
		})
	})

	Describe("GetPodNamespace", func() {
		It("works", func() {
			actual, err := apiserver.GetPodNamespace(context.TODO(), &v1.GetPodNamespaceRequest{})
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetPodNamespaceResponse{Namespace: podNamespace}
			ExpectEqualProtoMessages(actual, expected)
		})
	})
})
