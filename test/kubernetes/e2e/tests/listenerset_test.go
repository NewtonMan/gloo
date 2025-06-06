package tests_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	. "github.com/solo-io/gloo/test/kubernetes/e2e/tests"
	"github.com/solo-io/gloo/test/kubernetes/testutils/gloogateway"
	"github.com/solo-io/gloo/test/testutils"
)

// TestListenerSet is the function which executes a series of tests against a given installation
func TestListenerSet(t *testing.T) {
	ctx := context.Background()
	installNs, nsEnvPredefined := envutils.LookupOrDefault(testutils.InstallNamespace, "ls-test")
	testInstallation := e2e.CreateTestInstallation(
		t,
		&gloogateway.Context{
			InstallNamespace:          installNs,
			ProfileValuesManifestFile: e2e.KubernetesGatewayProfilePath,
			ValuesManifestFile:        e2e.EmptyValuesManifestPath,
			ValidationAlwaysAccept:    false,
			K8sGatewayEnabled:         true,
		},
	)

	if !listenerset.RequiredCrdExists(testInstallation) {
		t.Skip("Skipping as the XListenerSet CRD is not installed")
	}

	testHelper := e2e.MustTestHelper(ctx, testInstallation)

	// Set the env to the install namespace if it is not already set
	if !nsEnvPredefined {
		os.Setenv(testutils.InstallNamespace, installNs)
	}

	// We register the cleanup function _before_ we actually perform the installation.
	// This allows us to uninstall Gloo Gateway, in case the original installation only completed partially
	t.Cleanup(func() {
		if !nsEnvPredefined {
			os.Unsetenv(testutils.InstallNamespace)
		}
		if t.Failed() {
			testInstallation.PreFailHandler(ctx)
		}

		testInstallation.UninstallGlooGatewayWithTestHelper(ctx, testHelper)
	})

	// Install Gloo Gateway
	testInstallation.InstallGlooGatewayWithTestHelper(ctx, testHelper, 5*time.Minute)

	ListenerSetSuiteRunner().Run(ctx, t, testInstallation)
}
