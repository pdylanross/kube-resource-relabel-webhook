package relabel

import (
	"log/slog"

	"gomodules.xyz/jsonpatch/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Condition checks if we should modify an object.
type Condition interface {
	Satisfies(obj metaV1.Object) bool
}

// ConditionConfig is a config object that can construct a condition.
type ConditionConfig interface {
	// ToCondition constructs a condition
	ToCondition() Condition
}

// Action applies a change to an object.
type Action interface {
	Update(obj metaV1.Object) []jsonpatch.Operation
}

// ActionConfig is a config object that can construct an action.
type ActionConfig interface {
	// ToAction constructs an action
	ToAction() Action
}

// Rule is a single rule describing if and how we should relabel an object.
type Rule struct {
	Conditions []Condition
	Actions    []Action
	Name       string
}

// Evaluate a k8s object against this rule
// return if the object was modified.
func (r *Rule) Evaluate(obj metaV1.Object, logger *slog.Logger) []jsonpatch.Operation {
	l := logger.With(slog.String("rule-name", r.Name),
		slog.String("namespace", obj.GetNamespace()),
		slog.String("name", obj.GetName()))

	l.Debug("evaluating object")

	for _, c := range r.Conditions {
		if !c.Satisfies(obj) {
			l.Debug("object didn't satisfy preconditions")
			return []jsonpatch.Operation{}
		}
	}

	var operations []jsonpatch.Operation

	for _, a := range r.Actions {
		newPatches := a.Update(obj)
		operations = append(operations, newPatches...)
	}

	l.Debug("obj pending changes", slog.Int("num-changes", len(operations)))
	return operations
}
