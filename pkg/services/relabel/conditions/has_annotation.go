package conditions

import (
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type hasAnnotationCondition struct {
	mc mapCheckCondition
}

func (h *hasAnnotationCondition) Satisfies(obj metaV1.Object) bool {
	return h.mc.isMatch(obj.GetAnnotations())
}

func NewHasAnnotationCondition(keys []string, values []string, match map[string]string) relabel.Condition {
	return &hasAnnotationCondition{mc: mapCheckCondition{
		keys:   keys,
		values: values,
		match:  match,
	}}
}
