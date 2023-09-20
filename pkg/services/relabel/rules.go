package relabel

import (
	"log/slog"

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
	Update(obj metaV1.Object)
}

// ActionConfig is a config object that can construct an action.
type ActionConfig interface {
	// ToAction constructs an action
	ToAction() Action
}

// Rule is a single rule describing if and how we should relabel an object.
type Rule struct {
	conditions []Condition
	actions    []Action
	name       string
}

func NewRelabelRule(name string, conditions []Condition, actions []Action) *Rule {
	return &Rule{
		name:       name,
		actions:    actions,
		conditions: conditions,
	}
}

// Evaluate a k8s object against this rule
// return if the object was modified.
func (r *Rule) Evaluate(obj metaV1.Object, logger *slog.Logger) bool {
	l := logger.With(slog.String("rule-name", r.name),
		slog.String("namespace", obj.GetNamespace()),
		slog.String("name", obj.GetName()))

	l.Debug("evaluating object")

	for _, c := range r.conditions {
		if !c.Satisfies(obj) {
			l.Debug("object didn't satisfy preconditions")
			return false
		}
	}

	l.Debug("updating object")
	for _, a := range r.actions {
		a.Update(obj)
	}

	return true
}
