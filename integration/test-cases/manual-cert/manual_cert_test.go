//go:build integration_tests

package manual_cert

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"

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
	fixture.HelmInstallRelabel("./integration/test-cases/manual-cert/relabel-values.yaml")

	tests := map[string]func(clientset *kubernetes.Clientset, t2 *testing.T){
		"Manual_NoMatch_NoChange":                                              Manual_NoMatch_NoChange,
		"Manual_DagMatch_NoAnnotations_AnnotationsAdded":                       Manual_DagMatch_NoAnnotations_AnnotationsAdded,
		"Manual_DagMatch_ExistingAnnotations_AnnotationsUpdated":               Manual_DagMatch_ExistingAnnotations_AnnotationsUpdated,
		"Manual_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated": Manual_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated,
		"Manual_DagMatch_ExistingAnnotationsExactMatch_NoChange":               Manual_DagMatch_ExistingAnnotationsExactMatch_NoChange,
		"Manual_FluentdMatch_NoLabels_LabelsAdded":                             Manual_FluentdMatch_NoLabels_LabelsAdded,
		"Manual_FluentdMatch_ExistingLabels_LabelsUpdated":                     Manual_FluentdMatch_ExistingLabels_LabelsUpdated,
		"Manual_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated":           Manual_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated,
		"Manual_FluentdMatch_ExistingLabelsExactMatch_NoChange":                Manual_FluentdMatch_ExistingLabelsExactMatch_NoChange,
	}

	for name, f := range tests {
		t.Run(name, func(t *testing.T) {
			f(kubeClientset, t)
		})
	}
}

func Manual_NoMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "no-match-no-change",
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 0, len(result.Annotations))
	assert.Equal(t, 0, len(result.Labels))
}

func Manual_DagMatch_NoAnnotations_AnnotationsAdded(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-noannotations-annotationsadded",
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "true", result.Annotations["karpenter.sh/do-not-evict"])
}

func Manual_DagMatch_ExistingAnnotations_AnnotationsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-update-value-keep-old",
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
			Annotations: map[string]string{
				"other-thing": "false",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 2, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "false", result.Annotations["other-thing"])
	assert.Equal(t, "true", result.Annotations["karpenter.sh/do-not-evict"])
}

func Manual_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-update-valueinplace",
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
			Annotations: map[string]string{
				"karpenter.sh/do-not-evict": "false",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "true", result.Annotations["karpenter.sh/do-not-evict"])
}

func Manual_DagMatch_ExistingAnnotationsExactMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-alreadycorrectannotations-nochange",
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
			Annotations: map[string]string{
				"karpenter.sh/do-not-evict": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "true", result.Annotations["karpenter.sh/do-not-evict"])
}

func Manual_FluentdMatch_NoLabels_LabelsAdded(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-nolabels-labelsadded",
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "true", result.Labels["fluentd.active"])
}

func Manual_FluentdMatch_ExistingLabels_LabelsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-existinglabels-labelsupdated",
			Labels: map[string]string{
				"other-thing": "yes",
			},
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 2, len(result.Labels))

	assert.Equal(t, "true", result.Labels["fluentd.active"])
	assert.Equal(t, "yes", result.Labels["other-thing"])
}

func Manual_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-existinglabelswrongvalue-labelsupdated",
			Labels: map[string]string{
				"fluentd.active": "false",
			},
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "true", result.Labels["fluentd.active"])
}

func Manual_FluentdMatch_ExistingLabelsExactMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "dagmatch-existinglabelsexactmatch-nochange",
			Labels: map[string]string{
				"fluentd.active": "true",
			},
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods("default").Create(context.Background(), &pod, metaV1.CreateOptions{})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, 1, len(result.Annotations))
	assert.Equal(t, 1, len(result.Labels))

	assert.Equal(t, "true", result.Labels["fluentd.active"])
}
