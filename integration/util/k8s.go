//go:build integration_tests

package util

import (
	"fmt"
	"time"

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

func (itf *IntegrationTestFixture) HelmInstallRelabel(relValuesFile string) {
	_, err := itf.RunCommand(
		"helm",
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
		"-f",
		itf.GetFullPath(relValuesFile),
		"--atomic",
		"--kubeconfig",
		itf.kubeconfig)

	if err != nil {
		itf.t.Errorf("could not install helm chart: %s", err.Error())
	}

	time.Sleep(1 * time.Second)
}

var SleepPodTerminationPeriod int64 = 1
var SleepPodSpec = corev1.PodSpec{
	RestartPolicy:                 corev1.RestartPolicyAlways,
	TerminationGracePeriodSeconds: &SleepPodTerminationPeriod,
	Containers: []corev1.Container{
		corev1.Container{
			Name:  "sleepy",
			Image: "alpine",
			Command: []string{
				"sleep",
				"100000",
			},
		},
	},
}
