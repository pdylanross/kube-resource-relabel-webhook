//go:build integration_tests

package manual_cert

import (
	"context"
	"log/slog"
	"testing"

	common_tests "github.com/pdylanross/kube-resource-relabel-webhook/integration/common-tests"

	"github.com/stretchr/testify/require"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pdylanross/kube-resource-relabel-webhook/integration/util"
)

func TestManualCert(t *testing.T) {
	fixture := util.NewTestFixture(t)
	fixture.Start()
	defer fixture.Close()

	kubeClientset := fixture.GetKubernetesClient()

	slog.Info("creating cert secret")
	certSecret := fixture.LoadSecret("./integration/test-cases/manual-cert/secret.yaml")
	_, err := kubeClientset.CoreV1().Secrets("default").Create(context.Background(), certSecret, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating secret")

	slog.Info("installing relabeler")
	fixture.HelmInstallRelabel([]string{"./integration/test-cases/manual-cert/relabel-values.yaml", common_tests.CommonValuesFile})

	common_tests.RunCommonTests(t, kubeClientset)
}
