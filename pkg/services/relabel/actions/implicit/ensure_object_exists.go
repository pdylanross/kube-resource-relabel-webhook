package implicit

import "gomodules.xyz/jsonpatch/v3"

func ensureMapObjectExists(obj map[string]string, path string) []jsonpatch.Operation {
	if obj == nil {
		return []jsonpatch.Operation{
			jsonpatch.NewOperation("add", path, struct{}{}),
		}
	}

	return []jsonpatch.Operation{}
}
