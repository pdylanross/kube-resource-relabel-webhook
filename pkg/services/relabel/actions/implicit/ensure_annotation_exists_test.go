package implicit

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnsureAnnotationExistsAction_Update(t *testing.T) {
	action := EnsureAnnotationExistsAction{}

	t.Run("nil annotations", func(t *testing.T) {
		pod := &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{},
		}

		results := action.Update(pod)

		require.Equal(t, 1, len(results))
		assert.Equal(t, "/metadata/annotations", results[0].Path)
		assert.Equal(t, "add", results[0].Operation)
	})

	t.Run("not nil annotations", func(t *testing.T) {
		pod := &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{
				Annotations: map[string]string{},
			},
		}

		results := action.Update(pod)

		require.Equal(t, 0, len(results))
	})
}
