package actions

import (
	"maps"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ensureAnnotationAction struct {
	annotations map[string]string
}

func (e *ensureAnnotationAction) Update(obj metaV1.Object) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	maps.Copy(annotations, e.annotations)
	obj.SetAnnotations(annotations)
}

func NewEnsureAnnotationAction(c map[string]string) relabel.Action {
	return &ensureAnnotationAction{annotations: c}
}
