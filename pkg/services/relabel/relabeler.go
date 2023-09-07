package relabel

import (
	"log/slog"

	"gomodules.xyz/jsonpatch/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Relabeler is the main entrypoint into the relabeling process
// it keeps a runtime optimized set of rules to relabel objects.
type Relabeler interface {
	Evaluate(obj metaV1.Object) []jsonpatch.Operation
}

func NewRelabeler(rules []Rule) Relabeler {
	return &relabeler{rules: rules, logger: slog.With(slog.String("component", "relabeler"))}
}

type relabeler struct {
	rules  []Rule
	logger *slog.Logger
}

// Evaluate a k8s object against the current ruleset
// returning if the object was modified.
func (r *relabeler) Evaluate(obj metaV1.Object) []jsonpatch.Operation {
	var operations []jsonpatch.Operation
	for _, rule := range r.rules {
		newOps := rule.Evaluate(obj, r.logger)
		operations = append(operations, newOps...)
	}

	return operations
}
