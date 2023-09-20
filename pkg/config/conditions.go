package config

import (
	"fmt"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	conditions2 "github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel/conditions"
)

// GetConcrete deserializes the rule into a concrete ConditionConfig.
func (act *RelabelConfigRuleCondition) GetConcrete() (relabel.ConditionConfig, error) {
	switch act.Type {
	case "has-label":
		var ret HasLabelCondition
		err := act.Value.Decode(&ret)
		return &ret, err
	case "has-annotation":
		var ret HasAnnotationCondition
		err := act.Value.Decode(&ret)
		return &ret, err
	case "is-type":
		var ret IsTypeCondition
		err := act.Value.Decode(&ret)
		return &ret, err
	}

	return nil, fmt.Errorf("unknown action type %s", act.Type)
}

// HasLabelCondition checks existence of labels.
type HasLabelCondition struct {
	// Keys matches on key names
	Keys []string `yaml:"keys,omitempty"`
	// Values matches on values
	Values []string `yaml:"values,omitempty"`
	// Match matches items exactly
	Match map[string]string `yaml:"match,omitempty"`
}

func (h *HasLabelCondition) ToCondition() relabel.Condition {
	return conditions2.NewHasLabelCondition(h.Keys, h.Values, h.Match)
}

// HasAnnotationCondition checks existence of annotations.
type HasAnnotationCondition struct {
	// Keys matches on key names
	Keys []string `yaml:"keys,omitempty"`
	// Values matches on values
	Values []string `yaml:"values,omitempty"`
	// Match matches items exactly
	Match map[string]string `yaml:"match,omitempty"`
}

func (h *HasAnnotationCondition) ToCondition() relabel.Condition {
	return conditions2.NewHasAnnotationCondition(h.Keys, h.Values, h.Match)
}

// IsTypeCondition checks if an object is of a specific type.
type IsTypeCondition struct {
	// Group is the k8s api group to check against
	Group string `yaml:"group,omitempty"`
	// Version is the k8s api version to check against
	Version string `yaml:"version,omitempty"`
	// Kind is the k8s api kind to check against
	Kind string `yaml:"kind,omitempty"`
}

func (i *IsTypeCondition) ToCondition() relabel.Condition {
	return conditions2.NewIsTypeCondition(i.Group, i.Version, i.Kind)
}
