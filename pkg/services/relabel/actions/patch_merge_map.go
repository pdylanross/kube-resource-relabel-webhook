package actions

import (
	"fmt"
	"strings"

	"gomodules.xyz/jsonpatch/v3"
)

func patchMergeMap(path string, origin map[string]string, newValues map[string]string) []jsonpatch.Operation {
	if origin == nil {
		return []jsonpatch.Operation{
			jsonpatch.NewOperation("add", path, newValues),
		}
	}

	var operations []jsonpatch.Operation

	for newKey, newValue := range newValues {
		var op string
		if oldValue, ok := origin[newKey]; ok {
			if oldValue == newValue {
				continue
			}

			op = "replace"
		} else {
			op = "add"
		}

		operations = append(operations, jsonpatch.NewOperation(op, fmt.Sprintf("%s/%s", path, sanitizeKeyForJSONPatch(newKey)), newValue))
	}

	return operations
}

func sanitizeKeyForJSONPatch(key string) string {
	ret := strings.Replace(key, "~", "~0", -1)
	return strings.Replace(ret, "/", "~1", -1)
}
