package actions

import (
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/services/relabel"
	"gomodules.xyz/jsonpatch/v3"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ensureAnnotationAction struct {
	annotations map[string]string
}

func (e *ensureAnnotationAction) Update(obj metaV1.Object) []jsonpatch.Operation {
	return patchMergeMap("/metadata/annotations", obj.GetAnnotations(), e.annotations)
}

func NewEnsureAnnotationAction(c map[string]string) relabel.Action {
	return &ensureAnnotationAction{annotations: c}
}
