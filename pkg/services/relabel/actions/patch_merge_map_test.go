package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatchMergeMap_AddValue(t *testing.T) {
	path := "/metadata/annotations"
	origin := map[string]string{
		"thing1": "value1",
	}
	newValues := map[string]string{
		"newKey": "newValue",
	}

	results := patchMergeMap(path, origin, newValues)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "/metadata/annotations/newKey", results[0].Path)
	assert.Equal(t, "newValue", results[0].Value)
}

func TestPatchMergeMap_AddMultipleValues(t *testing.T) {
	path := "/metadata/annotations"
	origin := map[string]string{
		"thing1": "value1",
		"thing2": "value2",
	}
	newValues := map[string]string{
		"newKey":  "newValue",
		"newKey2": "newValue",
	}

	results := patchMergeMap(path, origin, newValues)

	assert.Equal(t, 2, len(results))

	for _, v := range results {
		if v.Path == "/metadata/annotations/newKey" {
			assert.Equal(t, "newValue", v.Value)
			assert.Equal(t, "add", v.Operation)
		} else if v.Path == "/metadata/annotations/newKey2" {
			assert.Equal(t, "newValue", v.Value)
			assert.Equal(t, "add", v.Operation)
		} else {
			t.Errorf("unknown path %s", v.Path)
		}
	}
}

func TestPatchMergeMap_ReplaceValue(t *testing.T) {
	path := "/metadata/annotations"
	origin := map[string]string{
		"oldKey": "value1",
	}
	newValues := map[string]string{
		"oldKey": "newValue",
	}

	results := patchMergeMap(path, origin, newValues)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "replace", results[0].Operation)
	assert.Equal(t, "/metadata/annotations/oldKey", results[0].Path)
	assert.Equal(t, "newValue", results[0].Value)
}

func TestPatchMergeMap_ExactMatchNoOp(t *testing.T) {
	path := "/metadata/annotations"
	origin := map[string]string{
		"oldKey": "oldValue",
	}
	newValues := map[string]string{
		"oldKey": "oldValue",
	}

	results := patchMergeMap(path, origin, newValues)

	assert.Equal(t, 0, len(results))
}

func TestPatchMergeMap_ValueWithSlash(t *testing.T) {
	path := "/metadata/annotations"
	origin := map[string]string{}
	newValues := map[string]string{
		"old/Key": "newValue",
	}

	results := patchMergeMap(path, origin, newValues)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "/metadata/annotations/old~1Key", results[0].Path)
	assert.Equal(t, "newValue", results[0].Value)
}

func TestPatchMergeMap_NilOrigin(t *testing.T) {
	path := "/metadata/annotations"
	var origin map[string]string
	newValues := map[string]string{
		"old/Key": "newValue",
	}

	results := patchMergeMap(path, origin, newValues)

	assert.Equal(t, 1, len(results))
	assert.Equal(t, "add", results[0].Operation)
	assert.Equal(t, "/metadata/annotations", results[0].Path)
	assert.Equal(t, newValues, results[0].Value)
}
