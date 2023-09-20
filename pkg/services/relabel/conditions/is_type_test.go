package conditions

import (
	"testing"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestIsTypeCondition_Satisfies(t *testing.T) {
	condition := NewIsTypeCondition("", "v1", "pod")
	pod := &corev1.Pod{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "V1",
		},
	}

	assert.True(t, condition.Satisfies(pod))
}
