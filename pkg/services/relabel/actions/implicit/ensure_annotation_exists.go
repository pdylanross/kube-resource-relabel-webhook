package implicit

import (
	"gomodules.xyz/jsonpatch/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EnsureAnnotationExistsAction struct{}

func (e *EnsureAnnotationExistsAction) Update(obj metaV1.Object) []jsonpatch.Operation {
	return ensureMapObjectExists(obj.GetAnnotations(), "/metadata/annotations")
}
