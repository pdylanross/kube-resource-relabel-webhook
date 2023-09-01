package actions

import (
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	k8sapiv1 "k8s.io/api/core/v1"
)

func TestNewEnsureAnnotationAction(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureAnnotationAction(annotations)

	assert.NotNil(t, action)
}

func TestEnsureAnnotationAction_UpdateAddsNewField(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureAnnotationAction(annotations)

	pod := &k8sapiv1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"should-not": "change",
			},
			Annotations: map[string]string{
				"existing": "shouldbeuntouched",
			},
		},
	}

	action.Update(pod)

	assert.NotNil(t, pod.ObjectMeta.Annotations)
	assert.Equal(t, "test-pod", pod.Name)
	assert.Equal(t, "default", pod.Namespace)
	assert.Equal(t, map[string]string{"should-not": "change"}, pod.Labels)
	assert.Equal(t, map[string]string{"existing": "shouldbeuntouched", "test": "value"}, pod.Annotations)
}

func TestEnsureAnnotationAction_UpdateAddsNewField_PriorBlank(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureAnnotationAction(annotations)

	pod := &k8sapiv1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"should-not": "change",
			},
			Annotations: nil,
		},
	}

	action.Update(pod)

	assert.NotNil(t, pod.ObjectMeta.Annotations)
	assert.Equal(t, "test-pod", pod.Name)
	assert.Equal(t, "default", pod.Namespace)
	assert.Equal(t, map[string]string{"should-not": "change"}, pod.Labels)
	assert.Equal(t, map[string]string{"test": "value"}, pod.Annotations)
}

func TestEnsureAnnotationAction_UpdateExistingField(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureAnnotationAction(annotations)

	pod := &k8sapiv1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"should-not": "change",
			},
			Annotations: map[string]string{
				"test": "wontbethisnomore",
			},
		},
	}

	action.Update(pod)

	assert.NotNil(t, pod.ObjectMeta.Annotations)
	assert.Equal(t, "test-pod", pod.Name)
	assert.Equal(t, "default", pod.Namespace)
	assert.Equal(t, map[string]string{"should-not": "change"}, pod.Labels)
	assert.Equal(t, map[string]string{"test": "value"}, pod.Annotations)
}
