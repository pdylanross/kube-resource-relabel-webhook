package relabel

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Relabeler is the main entrypoint into the relabeling process
// it keeps a runtime optimized set of rules to relabel objects.
type Relabeler interface {
	Evaluate(obj metaV1.Object) bool
}

func NewRelabeler(rules []Rule) Relabeler {
	return &relabeler{rules}
}

type relabeler struct {
	rules []Rule
}

// Evaluate a k8s object against the current ruleset
// returning if the object was modified.
func (r *relabeler) Evaluate(obj metaV1.Object) bool {
	modified := false
	for _, rule := range r.rules {
		if rule.Evaluate(obj) {
			modified = true
		}
	}

	return modified
}
