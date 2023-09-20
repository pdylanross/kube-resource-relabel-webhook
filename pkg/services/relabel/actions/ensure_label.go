package actions

import (
	"maps"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ensureLabelAction struct {
	labels map[string]string
}

func (e *ensureLabelAction) Update(obj metaV1.Object) {
	labels := obj.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	maps.Copy(labels, e.labels)
	obj.SetLabels(labels)
}

func NewEnsureLabelAction(c map[string]string) relabel.Action {
	return &ensureLabelAction{labels: c}
}
