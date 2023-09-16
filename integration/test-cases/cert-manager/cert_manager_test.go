package cert_manager

import (
	"log/slog"
	"testing"

	common_tests "github.com/pdylanross/kube-resource-relabel-webhook/integration/common-tests"

	"github.com/pdylanross/kube-resource-relabel-webhook/integration/util"
)

func TestCertManager(t *testing.T) {
	fixture := util.NewTestFixture(t)
	fixture.Start()
	defer fixture.Close()

	kubeClientset := fixture.GetKubernetesClient()

	slog.Info("installing cert manager")
	fixture.RunKubectlCmd("apply", "-f", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml")
	fixture.WaitForDeploymentReady("cert-manager-webhook", "cert-manager")
	fixture.WaitForDeploymentReady("cert-manager", "cert-manager")
	fixture.WaitForDeploymentReady("cert-manager-cainjector", "cert-manager")

	slog.Info("installing relabeler")
	fixture.HelmInstallRelabel([]string{"./integration/test-cases/cert-manager/relabel-values.yaml", common_tests.CommonValuesFile})

	common_tests.RunCommonTests(t, kubeClientset)
}
