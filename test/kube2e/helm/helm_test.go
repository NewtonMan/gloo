package helm_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"

	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"

	"github.com/solo-io/solo-projects/install/helm/gloo-ee/generate"
	osskube2e "github.com/solo-io/solo-projects/test/kube2e"
	admission_v1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	admission_v1_types "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// for testing upgrades from a gloo version before the gloo/gateway merge and
	// before https://github.com/solo-io/gloo/pull/6349 was fixed
	// TODO delete tests once this version is no longer supported https://github.com/solo-io/gloo/issues/6661
	versionBeforeGlooGatewayMerge = "1.11.0"

	versionBeforeCustomReadinessProbeFix = "1.15.2"

	glooChartName  = "gloo"
	glooeeRepoName = "https://storage.googleapis.com/gloo-ee-helm"
)

var _ = Describe("Installing and upgrading GlooEE via helm", func() {

	var (
		chartUri string

		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config
		err    error

		kubeClientset *kubernetes.Clientset

		testHelper *helper.SoloTestHelper

		// if set, the test will install the initial version of gloo from a released version (rather than local version) of the helm chart
		fromRelease string
		// if set, the test will upgrade to a released version of gloo
		targetReleasedVersion string
		// whether to set validation webhook's failurePolicy=Fail
		strictValidation bool
		// additional args to pass into the initial helm install
		additionalInstallArgs []string
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		kubeClientset, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		testHelper, err = osskube2e.GetEnterpriseTestHelper(ctx, namespace)
		Expect(err).NotTo(HaveOccurred())
		targetReleasedVersion = testHelper.ReleasedVersion
		if targetReleasedVersion != "" {
			chartUri = "glooe/gloo-ee"
		} else {
			chartUri = filepath.Join(testHelper.RootDir, testHelper.TestAssetDir, testHelper.HelmChartName+"-"+testHelper.ChartVersion()+".tgz")
		}

		fromRelease = ""
		strictValidation = false

		additionalInstallArgs = []string{}
	})

	JustBeforeEach(func() {
		if fromRelease == "" && targetReleasedVersion != "" {
			fromRelease = targetReleasedVersion
		}
		installGloo(testHelper, chartUri, fromRelease, strictValidation, additionalInstallArgs)
	})

	AfterEach(func() {
		uninstallGloo(testHelper, ctx, cancel)
	})

	// this is a subset of the helm upgrade tests done in the OSS repo
	Context("failurePolicy upgrades", func() {
		var webhookConfigClient admission_v1_types.ValidatingWebhookConfigurationInterface
		// Note: we are using the solo-apis clients instead of the solo-kit ones because the resources returned
		// by the solo-kit clients do not include creation timestamps. In these tests we are using creation timestamps
		// to check that the resources don't get deleted during the helm upgrades.
		var gatewayClientset gatewayv1.Clientset
		var glooClientset gloov1.Clientset

		BeforeEach(func() {
			webhookConfigClient = kubeClientset.AdmissionregistrationV1().ValidatingWebhookConfigurations()

			gatewayClientset, err = newGatewayClientsetFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())
			glooClientset, err = newGlooClientsetFromConfig(cfg)
			Expect(err).NotTo(HaveOccurred())

			fromRelease = versionBeforeGlooGatewayMerge
			strictValidation = false
		})

		getGatewayCreationTimestamp := func(name string) string {
			gw, err := gatewayClientset.Gateways().GetGateway(ctx, client.ObjectKey{
				Namespace: namespace,
				Name:      name,
			})
			Expect(err).NotTo(HaveOccurred())
			return gw.GetCreationTimestamp().String()
		}

		getUpstreamCreationTimestamp := func(name string) string {
			us, err := glooClientset.Upstreams().GetUpstream(ctx, client.ObjectKey{
				Namespace: namespace,
				Name:      name,
			})
			Expect(err).NotTo(HaveOccurred())
			return us.GetCreationTimestamp().String()
		}

		testFailurePolicyUpgrade := func(oldFailurePolicy admission_v1.FailurePolicyType, newFailurePolicy admission_v1.FailurePolicyType) {
			By(fmt.Sprintf("should start with gateway.validation.failurePolicy=%v", oldFailurePolicy))
			webhookConfig, err := webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(*webhookConfig.Webhooks[0].FailurePolicy).To(Equal(oldFailurePolicy))

			// to ensure the default Gateways and Upstreams were not deleted during upgrade, compare their creation timestamps before and after the upgrade
			gwTimestampBefore := getGatewayCreationTimestamp("gateway-proxy")
			gwSslTimestampBefore := getGatewayCreationTimestamp("gateway-proxy-ssl")
			extauthTimestampBefore := getUpstreamCreationTimestamp("extauth")
			extauthSidecarTimestampBefore := getUpstreamCreationTimestamp("extauth-sidecar")
			ratelimitTimestampBefore := getUpstreamCreationTimestamp("rate-limit")

			// upgrade to the new failurePolicy type
			var newStrictValue = false
			if newFailurePolicy == admission_v1.Fail {
				newStrictValue = true
			}
			upgradeGloo(testHelper, chartUri, fromRelease, targetReleasedVersion, newStrictValue, []string{})

			By(fmt.Sprintf("should have updated to gateway.validation.failurePolicy=%v", newFailurePolicy))
			webhookConfig, err = webhookConfigClient.Get(ctx, "gloo-gateway-validation-webhook-"+testHelper.InstallNamespace, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(*webhookConfig.Webhooks[0].FailurePolicy).To(Equal(newFailurePolicy))

			By("Gateway creation timestamps should not have changed")
			gwTimestampAfter := getGatewayCreationTimestamp("gateway-proxy")
			Expect(gwTimestampBefore).To(Equal(gwTimestampAfter))
			gwSslTimestampAfter := getGatewayCreationTimestamp("gateway-proxy-ssl")
			Expect(gwSslTimestampBefore).To(Equal(gwSslTimestampAfter))
			extauthTimestampAfter := getUpstreamCreationTimestamp("extauth")
			Expect(extauthTimestampBefore).To(Equal(extauthTimestampAfter))
			extauthSidecarTimestampAfter := getUpstreamCreationTimestamp("extauth-sidecar")
			Expect(extauthSidecarTimestampBefore).To(Equal(extauthSidecarTimestampAfter))
			ratelimitTimestampAfter := getUpstreamCreationTimestamp("rate-limit")
			Expect(ratelimitTimestampBefore).To(Equal(ratelimitTimestampAfter))
		}

		Context("starting from before the gloo/gateway merge, with failurePolicy=Ignore", func() {
			BeforeEach(func() {
				fromRelease = versionBeforeGlooGatewayMerge
				strictValidation = false
			})
			It("can upgrade to current release, with failurePolicy=Ignore", func() {
				testFailurePolicyUpgrade(admission_v1.Ignore, admission_v1.Ignore)
			})
			It("can upgrade to current release, with failurePolicy=Fail", func() {
				testFailurePolicyUpgrade(admission_v1.Ignore, admission_v1.Fail)
			})
		})
		Context("starting from helm hook release, with failurePolicy=Fail", func() {
			BeforeEach(func() {
				// The original fix for installing with failurePolicy=Fail (https://github.com/solo-io/gloo/issues/6213)
				// went into gloo-ee v1.11.9. It turned the Gloo custom resources into helm hooks to guarantee ordering,
				// however it caused additional issues so we moved away from using helm hooks. This test is to ensure
				// we can successfully upgrade from the helm hook release to the current release.
				// TODO delete tests once this version is no longer supported https://github.com/solo-io/gloo/issues/6661
				fromRelease = "1.11.9"
				strictValidation = true
			})
			It("can upgrade to current release, with failurePolicy=Fail", func() {
				testFailurePolicyUpgrade(admission_v1.Fail, admission_v1.Fail)
			})
		})
	})
})

func installGloo(testHelper *helper.SoloTestHelper, chartUri string, fromRelease string, strictValidation bool, additionalArgs []string) {
	valueOverrideFile := getHelmValuesFile("helm.yaml")

	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}
	if fromRelease != "" {
		osskube2e.RunAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, glooeeRepoName,
			"--force-update")
		args = append(args, testHelper.HelmChartName+"/gloo-ee",
			"--version", fmt.Sprintf("v%s", fromRelease))
	} else {
		args = append(args, chartUri)
	}
	args = append(args, "-n", testHelper.InstallNamespace,
		// As most CD tools wait for resources to be ready before marking the release as successful,
		// we're emulating that here by passing these two flags.
		// This way we ensure that we indirectly add support for CD tools
		"--wait",
		"--wait-for-jobs",
		// We run our e2e tests on a kind cluster, but kind hasn’t implemented LoadBalancer support.
		// This leads to the service being in a pending state.
		// Since the --wait flag is set, this can cause the upgrade to fail
		// as helm waits until the service is ready and eventually times out.
		// So instead we use the service type as ClusterIP to work around this limitation.
		"--set", "gloo.gatewayProxies.gatewayProxy.service.type=ClusterIP",
		"--create-namespace",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", valueOverrideFile)
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}
	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)
	osskube2e.RunAndCleanCommand("helm", args...)

	// Check that everything is OK
	osskube2e.CheckGlooHealthy(testHelper)
}

func upgradeGloo(testHelper *helper.SoloTestHelper, chartUri string, fromRelease string, targetReleasedVersion string, strictValidation bool, additionalArgs []string) {
	upgradeCrds(fromRelease, chartUri, targetReleasedVersion)

	valueOverrideFile := getHelmUpgradeValuesOverrideFile()

	var args = []string{"upgrade", testHelper.HelmChartName, chartUri,
		"-n", testHelper.InstallNamespace,
		// As most CD tools wait for resources to be ready before marking the release as successful,
		// we're emulating that here by passing these two flags.
		// This way we ensure that we indirectly add support for CD tools
		"--wait",
		"--wait-for-jobs",
		// We run our e2e tests on a kind cluster, but kind hasn’t implemented LoadBalancer support.
		// This leads to the service being in a pending state.
		// Since the --wait flag is set, this can cause the upgrade to fail
		// as helm waits until the service is ready and eventually times out.
		// So instead we use the service type as ClusterIP to work around this limitation.
		"--set", "gloo.gatewayProxies.gatewayProxy.service.type=ClusterIP",
		"--set-string", "license_key=" + testHelper.LicenseKey,
		"--values", valueOverrideFile}
	if targetReleasedVersion != "" {
		args = append(args, "--version", targetReleasedVersion)
	}
	if strictValidation {
		args = append(args, strictValidationArgs...)
	}
	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)
	osskube2e.RunAndCleanCommand("helm", args...)

	// Check that everything is OK
	osskube2e.CheckGlooHealthy(testHelper)
}

func uninstallGloo(testHelper *helper.SoloTestHelper, ctx context.Context, cancel context.CancelFunc) {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGlooAll()
	Expect(err).NotTo(HaveOccurred())
	_, err = kube2e.MustKubeClient().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
	cancel()
}

// returns repository and version of the Gloo OSS dependency
func getGlooOSSDep(reqTemplateUri string) (string, string, error) {
	bytes, err := os.ReadFile(reqTemplateUri)
	if err != nil {
		return "", "", err
	}

	var dl generate.DependencyList
	err = yaml.Unmarshal(bytes, &dl)
	if err != nil {
		return "", "", err
	}

	for _, v := range dl.Dependencies {
		if v.Name == glooChartName {
			return v.Repository, v.Version, nil
		}
	}
	return "", "", eris.New("could not get gloo dependency info")
}

func upgradeCrds(fromRelease string, localChartUri string, publishedChartVersion string) {
	// if we're just upgrading within the same release, no need to reapply crds
	if fromRelease == "" {
		return
	}

	// untar the chart into a temp dir
	dir, err := os.MkdirTemp("", "unzipped-chart")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)
	if publishedChartVersion != "" {
		// Download the crds from the released chart
		osskube2e.RunAndCleanCommand("helm", "repo", "add", "glooe", glooeeRepoName, "--force-update")
		osskube2e.RunAndCleanCommand("helm", "pull", "glooe/gloo-ee", "--version", publishedChartVersion, "--untar", "--untardir", dir)
	} else {
		//untar the local chart to get the crds
		osskube2e.RunAndCleanCommand("tar", "-xvf", localChartUri, "--directory", dir)
	}
	// apply the crds
	crdDir := dir + "/gloo-ee/charts/gloo/crds"
	osskube2e.RunAndCleanCommand("kubectl", "apply", "-f", crdDir)
	// allow some time for the new crds to take effect
	time.Sleep(time.Second * 5)
}

var strictValidationArgs = []string{
	"--set", "gloo.gateway.validation.failurePolicy=Fail",
	"--set", "gloo.gateway.validation.allowWarnings=false",
	"--set", "gloo.gateway.validation.alwaysAcceptResources=false",
}

func getHelmValuesFile(filename string) string {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred(), "working dir could not be retrieved")
	helmUpgradeValuesFile := filepath.Join(cwd, "artifacts", filename)
	return helmUpgradeValuesFile

}

func getHelmUpgradeValuesOverrideFileForCustomReadinessProbe() (filename string) {
	return getHelmValuesFile("custom-readiness-probe.yaml")
}

func getHelmUpgradeValuesOverrideFile() (filename string) {
	return getHelmValuesFile("upgrade-override.yaml")
}

// calling NewClientsetFromConfig multiple times results in a race condition due to the use of the global scheme.Scheme.
// to avoid this, make a copy of the function here but use runtime.NewScheme instead of the global scheme
func newGatewayClientsetFromConfig(cfg *rest.Config) (gatewayv1.Clientset, error) {
	scheme := runtime.NewScheme()
	if err := gatewayv1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	client, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return gatewayv1.NewClientset(client), nil
}

func newGlooClientsetFromConfig(cfg *rest.Config) (gloov1.Clientset, error) {
	scheme := runtime.NewScheme()
	if err := gloov1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, err
	}
	client, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	return gloov1.NewClientset(client), nil
}
