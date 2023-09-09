package implicit

import (
	"gomodules.xyz/jsonpatch/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EnsureLabelsExistsAction struct{}

func (e *EnsureLabelsExistsAction) Update(obj metaV1.Object) []jsonpatch.Operation {
	return ensureMapObjectExists(obj.GetLabels(), "/metadata/labels")
}
