package conditions

import (
	"strings"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type isTypeCondition struct {
	group   string
	version string
	kind    string
}

func NewIsTypeCondition(group string, version string, kind string) relabel.Condition {
	return &isTypeCondition{
		group:   group,
		version: version,
		kind:    kind,
	}
}

func (i *isTypeCondition) Satisfies(obj metaV1.Object) bool {
	rto, ok := obj.(runtime.Object)
	if !ok {
		return false
	}

	gvk := rto.GetObjectKind().GroupVersionKind()

	return i.checkField(i.group, gvk.Group) &&
		i.checkField(i.version, gvk.Version) &&
		i.checkField(i.kind, gvk.Kind)
}

func (i *isTypeCondition) checkField(field string, value string) bool {
	if field == "" || field == "*" {
		return true
	}

	return strings.EqualFold(field, value)
}
