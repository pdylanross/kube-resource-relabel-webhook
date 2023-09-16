//go:build integration_tests

package common_tests

import (
	"context"
	"testing"

	"github.com/pdylanross/kube-resource-relabel-webhook/integration/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var CommonValuesFile = "./integration/common-tests/values.yaml"

func RunCommonTests(t *testing.T, clientset *kubernetes.Clientset) {
	tests := map[string]func(clientset *kubernetes.Clientset, t2 *testing.T){
		"Common_NoMatch_NoChange":                                              Common_NoMatch_NoChange,
		"Common_DagMatch_NoAnnotations_AnnotationsAdded":                       Common_DagMatch_NoAnnotations_AnnotationsAdded,
		"Common_DagMatch_ExistingAnnotations_AnnotationsUpdated":               Common_DagMatch_ExistingAnnotations_AnnotationsUpdated,
		"Common_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated": Common_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated,
		"Common_DagMatch_ExistingAnnotationsExactMatch_NoChange":               Common_DagMatch_ExistingAnnotationsExactMatch_NoChange,
		"Common_FluentdMatch_NoLabels_LabelsAdded":                             Common_FluentdMatch_NoLabels_LabelsAdded,
		"Common_FluentdMatch_ExistingLabels_LabelsUpdated":                     Common_FluentdMatch_ExistingLabels_LabelsUpdated,
		"Common_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated":           Common_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated,
		"Common_FluentdMatch_ExistingLabelsExactMatch_NoChange":                Common_FluentdMatch_ExistingLabelsExactMatch_NoChange,
	}

	for name, f := range tests {
		t.Run(name, func(t *testing.T) {
			f(clientset, t)
		})
	}
}

func Common_NoMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_DagMatch_NoAnnotations_AnnotationsAdded(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_DagMatch_ExistingAnnotations_AnnotationsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_DagMatch_ExistingAnnotationsExactMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_FluentdMatch_NoLabels_LabelsAdded(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_FluentdMatch_ExistingLabels_LabelsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
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

func Common_FluentdMatch_ExistingLabelsExactMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
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
