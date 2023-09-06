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

	results := action.Update(pod)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "/metadata/annotations/test", results[0].Path)
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "value", results[0].Value)
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

	results := action.Update(pod)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "/metadata/annotations/test", results[0].Path)
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "value", results[0].Value)
}
