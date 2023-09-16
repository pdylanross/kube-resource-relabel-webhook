//go:build integration_tests

package util

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

var universalDeserializer = clientsetscheme.Codecs.UniversalDeserializer()

func (itf *IntegrationTestFixture) LoadSecret(relPath string) *corev1.Secret {
	file := itf.LoadTestFile(relPath)

	var secret corev1.Secret

	_, _, err := universalDeserializer.Decode(file, nil, &secret)
	if err != nil {
		itf.t.Errorf("could not load secret from file %s: %s", relPath, err.Error())
	}

	return &secret
}

func (itf *IntegrationTestFixture) HelmInstallRelabel(relValuesFiles []string) {
	args := []string{
		"install",
		"relabel",
		itf.GetFullPath("./chart"),
		"-n",
		"default",
		"--set",
		fmt.Sprintf("image.tag=%s", itf.currentImageTag),
		"--set",
		fmt.Sprintf("image.repository=%s", itf.currentImageRepo),
		"--set",
		"image.pullPolicy=IfNotPresent",
		"--set",
		"fullnameOverride=kube-resource-relabel-webhook",
		"--atomic",
		"--kubeconfig",
		itf.kubeconfig,
	}

	for _, v := range relValuesFiles {
		args = append(args, "-f", itf.GetFullPath(v))
	}

	output, err := itf.RunCommand("helm", args...)

	if err != nil {
		itf.t.Errorf("could not install helm chart: %s \n output: %s", err.Error(), output)
	}

	itf.WaitForDeploymentReady("kube-resource-relabel-webhook", "default")
}

func (itf *IntegrationTestFixture) RunKubectlCmd(args ...string) {
	args = append(args, "--kubeconfig", itf.kubeconfig)

	output, err := itf.RunCommand("kubectl", args...)

	if err != nil {
		itf.t.Errorf("Could not run kubectl cmd: %s \n output: %s", err.Error(), output)
	}
}

func (itf *IntegrationTestFixture) WaitForDeploymentReady(name string, namespace string) {
	clientset := itf.GetKubernetesClient()
	getter := clientset.AppsV1().Deployments(namespace)
	exp := backoff.NewExponentialBackOff()
	exp.MaxElapsedTime = 1 * time.Minute
	exp.MaxInterval = 10 * time.Second

	deployment, err := getter.Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		itf.t.Errorf("error getting deployment %s in namespace %s: %s", name, namespace, err.Error())
	}

	for {
		next := exp.NextBackOff()
		if next == exp.Stop {
			itf.t.Errorf("timed out waiting for deployment %s readiness", name)
		}

		time.Sleep(next)
		slog.Info("wait for deployment", slog.String("name", name), slog.String("namespace", namespace))

		if deployment.Status.UnavailableReplicas == 0 {
			return
		}

		deployment, err = getter.Get(context.Background(), name, metaV1.GetOptions{})
		if err != nil {
			itf.t.Errorf("error getting deployment %s in namespace %s: %s", name, namespace, err.Error())
		}
	}
}

var SleepPodTerminationPeriod int64 = 1
var SleepPodSpec = corev1.PodSpec{
	RestartPolicy:                 corev1.RestartPolicyAlways,
	TerminationGracePeriodSeconds: &SleepPodTerminationPeriod,
	Containers: []corev1.Container{
		{
			Name:  "sleepy",
			Image: "alpine",
			Command: []string{
				"sleep",
				"100000",
			},
		},
	},
}
