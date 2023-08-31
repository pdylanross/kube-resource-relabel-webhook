package conditions

import (
	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/relabel"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type hasLabelCondition struct {
	mc mapCheckCondition
}

func (h *hasLabelCondition) Satisfies(obj metaV1.Object) bool {
	return h.mc.isMatch(obj.GetLabels())
}

func NewHasLabelCondition(keys []string, values []string, match map[string]string) relabel.Condition {
	return &hasLabelCondition{mc: mapCheckCondition{
		keys:   keys,
		values: values,
		match:  match,
	}}
}
