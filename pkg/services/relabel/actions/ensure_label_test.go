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

	action.Update(pod)

	expected := map[string]string{
		"test":     "value",
		"existing": "shouldbeuntouched",
	}
	assert.Equal(t, expected, pod.GetLabels())
}

func TestEnsureLabelAction_UpdateAddsNewField_PriorBlank(t *testing.T) {
	labels := map[string]string{
		"test": "value",
	}

	action := NewEnsureLabelAction(labels)

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

	action.Update(pod)

	expected := map[string]string{
		"test": "value",
	}
	assert.Equal(t, expected, pod.GetLabels())
}
