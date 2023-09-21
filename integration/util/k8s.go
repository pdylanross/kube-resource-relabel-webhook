//go:build integration_tests

package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"k8s.io/client-go/kubernetes"

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
		itf.t.Fatalf("could not load secret from file %s: %s", relPath, err.Error())
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
		"image.pullPolicy=Never",
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
		itf.t.Fatalf("could not install helm chart: %s \n output: %s", err.Error(), output)
	}

	itf.WaitForDeploymentReady("kube-resource-relabel-webhook", "default")
	itf.WaitForEndpointReady("kube-resource-relabel-webhook", "default")
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
	exp.MaxElapsedTime = 2 * time.Minute
	exp.MaxInterval = 10 * time.Second

	deployment, err := getter.Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		itf.t.Fatalf("error getting deployment %s in namespace %s: %s", name, namespace, err.Error())
	}

	for {
		next := exp.NextBackOff()
		if next == exp.Stop {
			itf.t.Fatalf("timed out waiting for deployment %s readiness", name)
		}

		time.Sleep(next)
		slog.Info("wait for deployment", slog.String("name", name), slog.String("namespace", namespace))

		if deployment.Status.UnavailableReplicas == 0 {
			return
		}

		deployment, err = getter.Get(context.Background(), name, metaV1.GetOptions{})
		if err != nil {
			itf.t.Fatalf("error getting deployment %s in namespace %s: %s", name, namespace, err.Error())
		}
	}
}

func (itf *IntegrationTestFixture) WaitForEndpointReady(name string, namespace string) {
	clientset := itf.GetKubernetesClient()
	getter := clientset.CoreV1().Endpoints(namespace)
	exp := backoff.NewExponentialBackOff()
	exp.MaxElapsedTime = 2 * time.Minute
	exp.MaxInterval = 10 * time.Second

	endpoints, err := getter.Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		itf.t.Fatalf("error getting endpoint %s in namespace %s: %s", name, namespace, err.Error())
	}

	for {
		next := exp.NextBackOff()
		if next == exp.Stop {
			itf.t.Fatalf("timed out waiting for endpoint %s", name)
		}

		time.Sleep(next)
		slog.Info("wait for endpoint", slog.String("name", name), slog.String("namespace", namespace))

		if len(endpoints.Subsets) > 0 && len(endpoints.Subsets[0].Addresses) > 0 && len(endpoints.Subsets[0].NotReadyAddresses) == 0 {
			return
		}

		endpoints, err = getter.Get(context.Background(), name, metaV1.GetOptions{})
		if err != nil {
			itf.t.Fatalf("error getting endpoint %s in namespace %s: %s", name, namespace, err.Error())
		}
	}
}

func (itf *IntegrationTestFixture) GetLogsForPodsByLabel(namespace string, selector string, clientset *kubernetes.Clientset) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: selector,
	})

	if err != nil {
		slog.Info("error fetching pods for debug pod dump: %s", slog.String("error", err.Error()))
		return
	}

	for _, pod := range pods.Items {
		itf.GetLogsForPod(namespace, pod.Name, clientset)
	}
}

func (itf *IntegrationTestFixture) GetLogsForPod(namespace string, name string, clientset *kubernetes.Clientset) {
	logsReq := clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{})
	logs, err := logsReq.Stream(context.Background())
	if err != nil {
		slog.Info("error opening pod log stream", slog.String("pod", name), slog.String("error", err.Error()))
		return
	}

	defer logs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		slog.Info("error copying logs to buf", slog.String("pod", name), slog.String("error", err.Error()))
		return
	}

	logStr := buf.String()

	slog.Info("pod logs: ", slog.String("pod", name))
	for _, ln := range strings.Split(logStr, "\n") {
		slog.Info(ln)
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
