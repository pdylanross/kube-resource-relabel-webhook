package actions

import (
	"fmt"

	"gomodules.xyz/jsonpatch/v3"
)

func patchMergeMap(path string, origin map[string]string, newValues map[string]string) []jsonpatch.Operation {
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

		operations = append(operations, jsonpatch.NewOperation(op, fmt.Sprintf("%s/%s", path, newKey), newValue))
	}

	return operations
}
