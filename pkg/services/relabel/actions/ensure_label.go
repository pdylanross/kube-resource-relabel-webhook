package actions

import (
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/services/relabel"
	"gomodules.xyz/jsonpatch/v3"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ensureLabelAction struct {
	labels map[string]string
}

func (e *ensureLabelAction) Update(obj metaV1.Object) []jsonpatch.Operation {
	return patchMergeMap("/metadata/labels", obj.GetLabels(), e.labels)
}

func NewEnsureLabelAction(c map[string]string) relabel.Action {
	return &ensureLabelAction{labels: c}
}
