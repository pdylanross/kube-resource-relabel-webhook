package relabel

import (
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
	Conditions []Condition
	Actions    []Action
	Name       string
}

// Evaluate a k8s object against this rule
// return if the object was modified.
func (r *Rule) Evaluate(obj metaV1.Object) bool {
	for _, c := range r.Conditions {
		if !c.Satisfies(obj) {
			return false
		}
	}

	modified := false
	for _, a := range r.Actions {
		a.Update(obj)
		modified = true
	}

	return modified
}
