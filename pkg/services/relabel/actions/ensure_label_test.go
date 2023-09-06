package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	k8sapiv1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewEnsureLabelAction(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureLabelAction(annotations)

	assert.NotNil(t, action)
}

func TestEnsureLabelAction_UpdateAddsNewField(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureLabelAction(annotations)

	pod := &k8sapiv1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				"should-not": "change",
			},
			Labels: map[string]string{
				"existing": "shouldbeuntouched",
			},
		},
	}

	results := action.Update(pod)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "/metadata/labels/test", results[0].Path)
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "value", results[0].Value)
}

func TestEnsureLabelAction_UpdateAddsNewField_PriorBlank(t *testing.T) {
	annotations := map[string]string{
		"test": "value",
	}

	action := NewEnsureLabelAction(annotations)

	pod := &k8sapiv1.Pod{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				"should-not": "change",
			},
			Labels: nil,
		},
	}

	results := action.Update(pod)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "/metadata/labels/test", results[0].Path)
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "value", results[0].Value)
}
