//go:build integration_tests

package common_tests

import (
	"context"
	"testing"

	networkingv1 "k8s.io/api/networking/v1"

	"github.com/pdylanross/kube-resource-relabel-webhook/integration/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var CommonValuesFile = "./integration/common-tests/values.yaml"
var TestNamespace = "common-tests"

func RunCommonTests(t *testing.T, clientset *kubernetes.Clientset) {
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
		ObjectMeta: metaV1.ObjectMeta{Name: TestNamespace},
	}, metaV1.CreateOptions{})
	require.Nil(t, err)

	tests := map[string]func(clientset *kubernetes.Clientset, t2 *testing.T){
		// pod tests
		"Common_NoMatch_NoChange":                                              Common_NoMatch_NoChange,
		"Common_DagMatch_NoAnnotations_AnnotationsAdded":                       Common_DagMatch_NoAnnotations_AnnotationsAdded,
		"Common_DagMatch_ExistingAnnotations_AnnotationsUpdated":               Common_DagMatch_ExistingAnnotations_AnnotationsUpdated,
		"Common_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated": Common_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated,
		"Common_DagMatch_ExistingAnnotationsExactMatch_NoChange":               Common_DagMatch_ExistingAnnotationsExactMatch_NoChange,
		"Common_FluentdMatch_NoLabels_LabelsAdded":                             Common_FluentdMatch_NoLabels_LabelsAdded,
		"Common_FluentdMatch_ExistingLabels_LabelsUpdated":                     Common_FluentdMatch_ExistingLabels_LabelsUpdated,
		"Common_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated":           Common_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated,
		"Common_FluentdMatch_ExistingLabelsExactMatch_NoChange":                Common_FluentdMatch_ExistingLabelsExactMatch_NoChange,
		"Common_FluentdMatch_NoLabels_MultiMatch":                              Common_FluentdMatch_NoLabels_MultiMatch,

		// ingress tests
		"Common_IngressResource_AddsAnnotation": Common_IngressResource_AddsAnnotation,
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
			Name:      "no-match-no-change",
			Namespace: TestNamespace,
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Nil(t, result.Annotations)

	assert.Nil(t, result.Labels)
}

func Common_DagMatch_NoAnnotations_AnnotationsAdded(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "dagmatch-noannotations-annotationsadded",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"karpenter.sh/do-not-evict": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"dag_id": "is-dag",
	}, result.Labels)
}

func Common_DagMatch_ExistingAnnotations_AnnotationsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "dagmatch-update-value-keep-old",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
			Annotations: map[string]string{
				"other-thing": "false",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"karpenter.sh/do-not-evict": "true",
		"other-thing":               "false",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"dag_id": "is-dag",
	}, result.Labels)
}

func Common_DagMatch_ExistingAnnotationsDifferentValue_AnnotationsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "dagmatch-update-valueinplace",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
			Annotations: map[string]string{
				"karpenter.sh/do-not-evict": "false",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"karpenter.sh/do-not-evict": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"dag_id": "is-dag",
	}, result.Labels)
}

func Common_DagMatch_ExistingAnnotationsExactMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "dagmatch-alreadycorrectannotations-nochange",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"dag_id": "is-dag",
			},
			Annotations: map[string]string{
				"karpenter.sh/do-not-evict": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"karpenter.sh/do-not-evict": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"dag_id": "is-dag",
	}, result.Labels)

	assert.Equal(t, "true", result.Annotations["karpenter.sh/do-not-evict"])
}

func Common_FluentdMatch_NoLabels_LabelsAdded(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "fluentdmatch-nolabels-labelsadded",
			Namespace: TestNamespace,
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Labels)

	assert.Equal(t, "true", result.Labels["fluentd.active"])
}

func Common_FluentdMatch_NoLabels_MultiMatch(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "fluentdmatch-nolabels-labelsadded-multi",
			Namespace: TestNamespace,
			Annotations: map[string]string{
				"fluentd.active": "true",
				"test2":          "value",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
		"test2":          "value",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
		"test2":          "value",
	}, result.Labels)
}

func Common_FluentdMatch_ExistingLabels_LabelsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "fluentdmatch-existinglabels-labelsupdated",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"other-thing": "yes",
			},
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
		"other-thing":    "yes",
	}, result.Labels)
}

func Common_FluentdMatch_ExistingLabelsWrongValue_LabelsUpdated(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "fluentdmatch-existinglabelswrongvalue-labelsupdated",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"fluentd.active": "false",
			},
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Labels)
}

func Common_FluentdMatch_ExistingLabelsExactMatch_NoChange(clientset *kubernetes.Clientset, t *testing.T) {
	pod := corev1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "fluentdmatch-existinglabelsexactmatch-nochange",
			Namespace: TestNamespace,
			Labels: map[string]string{
				"fluentd.active": "true",
			},
			Annotations: map[string]string{
				"fluentd.active": "true",
			},
		},
		Spec: *util.SleepPodSpec.DeepCopy(),
	}

	result, err := clientset.CoreV1().Pods(TestNamespace).Create(context.Background(), &pod, metaV1.CreateOptions{FieldValidation: "Strict"})
	require.NoError(t, err, "error creating test pod")

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Annotations)

	assert.Equal(t, map[string]string{
		"fluentd.active": "true",
	}, result.Labels)
}

func Common_IngressResource_AddsAnnotation(clientset *kubernetes.Clientset, t *testing.T) {
	pathType := networkingv1.PathTypePrefix
	ingress := networkingv1.Ingress{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "ingress-resource-add-annotation",
			Namespace: TestNamespace,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "test.example.com",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "test",
											Port: networkingv1.ServiceBackendPort{
												Name: "http",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	results, err := clientset.NetworkingV1().Ingresses(TestNamespace).Create(context.Background(), &ingress, metaV1.CreateOptions{})
	require.NoError(t, err)

	assert.Equal(t, map[string]string{"nginx.ingress.kubernetes.io/default-backend": "some-svc"}, results.Annotations)
}
