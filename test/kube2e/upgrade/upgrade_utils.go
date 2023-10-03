package upgrade

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/gomega"

	"github.com/solo-io/k8s-utils/testutils/helper"
	. "github.com/solo-io/solo-projects/test/kube2e"
)

func UpgradeGloo(testHelper *helper.SoloTestHelper, chartUri string, helmOverrideFilePath string, additionalArgs []string) {
	UpgradeCrds(chartUri, testHelper.ReleasedVersion)
	var args = []string{"upgrade", testHelper.HelmChartName, chartUri,
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
		"-n", testHelper.InstallNamespace,
		"--set-string", "license_key=" + testHelper.LicenseKey,
		"--values", helmOverrideFilePath}
	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)
	RunAndCleanCommand("helm", args...)

	//Check that everything is OK
	CheckGlooHealthy(testHelper)
}

func UpgradeGlooWithArgs(testHelper *helper.SoloTestHelper, chartUri string, helmOverrideFilePath string, additionalArgs []string) {
	UpgradeCrds(chartUri, testHelper.ReleasedVersion)
	var args = []string{"upgrade", testHelper.HelmChartName, chartUri,
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
		"-n", testHelper.InstallNamespace,
		"--set-string", "license_key=" + testHelper.LicenseKey,
		"--values", helmOverrideFilePath}

	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)
	RunAndCleanCommand("helm", args...)

	//Check that everything is OK
	CheckGlooHealthy(testHelper)
}

// UpgradeCrds first applies CRDs to a cluster when performing a `helm install` operation
// However, `helm upgrade` intentionally does not apply CRDs (https://helm.sh/docs/topics/charts/#limitations-on-crds)
// Before performing the upgrade, we must manually apply any CRDs that were introduced since v1.9.0
func UpgradeCrds(localChartUri string, publishedChartVersion string) {
	// untar the chart into a temp dir
	dir, err := os.MkdirTemp("", "unzipped-chart")
	Expect(err).NotTo(HaveOccurred())
	defer os.RemoveAll(dir)
	if publishedChartVersion != "" {
		// Download the crds from the released chart
		RunAndCleanCommand("helm", "repo", "add", "glooe", GlooeRepoName, "--force-update")
		RunAndCleanCommand("helm", "pull", "glooe/gloo-ee", "--version", publishedChartVersion, "--untar", "--untardir", dir)
	} else {
		//untar the local chart to get the crds
		RunAndCleanCommand("tar", "-xvf", localChartUri, "--directory", dir)
	}
	// apply the crds
	crdDir := dir + "/gloo-ee/charts/gloo/crds"
	RunAndCleanCommand("kubectl", "apply", "-f", crdDir)
	// allow some time for the new crds to take effect
	time.Sleep(time.Second * 5)
}
