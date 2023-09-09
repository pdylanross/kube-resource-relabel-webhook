package implicit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnsureLabelsExistsAction_Update(t *testing.T) {
	action := EnsureLabelsExistsAction{}

	t.Run("nil labels", func(t *testing.T) {
		pod := &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{},
		}

		results := action.Update(pod)

		require.Equal(t, 1, len(results))
		assert.Equal(t, "/metadata/labels", results[0].Path)
		assert.Equal(t, "add", results[0].Operation)
	})

	t.Run("not nil labels", func(t *testing.T) {
		pod := &v1.Pod{
			ObjectMeta: metaV1.ObjectMeta{
				Labels: map[string]string{},
			},
		}

		results := action.Update(pod)

		require.Equal(t, 0, len(results))
	})
}
