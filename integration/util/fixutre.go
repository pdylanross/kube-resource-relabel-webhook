//go:build integration_tests

package util

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"testing"

	petname "github.com/dustinkirkland/golang-petname"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type IntegrationTestFixture struct {
	t       *testing.T
	tempDir string

	kubeconfig      string
	kindClusterName string

	currentImageTag  string
	currentImageRepo string

	teardown bool
}

func NewTestFixture(t *testing.T) *IntegrationTestFixture {
	teardown := true
	var err error
	if val, ok := os.LookupEnv("TEST_TEARDOWN"); ok {
		teardown, err = strconv.ParseBool(val)
		if err != nil {
			slog.Warn("could not parse TEST_TEARDOWN as bool", slog.String("err", err.Error()))
			teardown = false
		}
	}

	tempDir := t.TempDir()
	return &IntegrationTestFixture{t: t, tempDir: tempDir, currentImageRepo: "relabel", currentImageTag: "dev", teardown: teardown}
}

func (itf *IntegrationTestFixture) Close() {
	if itf.teardown {
		slog.Info("stopping kind cluster", slog.String("name", itf.kindClusterName))
		output, err := itf.RunCommand("kind", "delete", "cluster", "--name", itf.kindClusterName)
		if err != nil {
			itf.t.Errorf("error deleting test cluster: %s \n %s", err, output)
		}
	} else {
		slog.Info("teardown disabled", slog.String("kubeconfig", itf.kubeconfig))
	}
}

func (itf *IntegrationTestFixture) Start() {
	kindConfig := itf.GetFullPath("./integration/kind-config.yaml")
	itf.kindClusterName = petname.Generate(2, "-")
	itf.kubeconfig = itf.GetFullPath(fmt.Sprintf("./integration/kubeconfig-%s", itf.kindClusterName))

	slog.Info("starting kind cluster", slog.String("name", itf.kindClusterName))
	output, err := itf.RunCommand("kind", "create", "cluster", "--config", kindConfig, "--kubeconfig", itf.kubeconfig, "--name", itf.kindClusterName)
	if err != nil {
		itf.t.Fatalf("error creating test cluster: %s \n %s", err, output)
	}

	slog.Info("loading container images")

	kindLoadOutput, err := itf.RunCommand(
		"kind",
		"load",
		"docker-image",
		"--name",
		itf.kindClusterName,
		fmt.Sprintf("%s:%s", itf.currentImageRepo, itf.currentImageTag))

	if err != nil {
		itf.t.Fatalf("could not load current relabel container to kind: %s \n %s", err.Error(), kindLoadOutput)
	}
}

func (itf *IntegrationTestFixture) RunCommand(cmd string, args ...string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("error finding cmd %s: %s", cmd, err.Error())
	}

	outputBytes, err := exec.Command(path, args...).CombinedOutput()

	return string(outputBytes), err
}

func (itf *IntegrationTestFixture) ResolveRepoRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	root := path.Join(path.Dir(filename), "../../")
	root = path.Clean(root)

	return root
}

func (itf *IntegrationTestFixture) GetFullPath(repoRelativeFilename string) string {
	root := itf.ResolveRepoRoot()
	return path.Clean(path.Join(root, repoRelativeFilename))
}

func (itf *IntegrationTestFixture) GetTempFile(filename string) string {
	return path.Join(itf.tempDir, filename)
}

func (itf *IntegrationTestFixture) LoadTestFile(relPath string) []byte {
	fullPath := itf.GetFullPath(relPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		itf.t.Fatalf("Could not load test file: %s", err)
	}

	return data
}

func (itf *IntegrationTestFixture) GetKubernetesClient() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", itf.kubeconfig)
	if err != nil {
		itf.t.Fatalf("error getting kubeconfig %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		itf.t.Fatalf("error building kubernetes client %s", err)
	}

	return clientset
}
