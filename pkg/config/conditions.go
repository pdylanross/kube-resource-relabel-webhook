package config

import (
	"fmt"

	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/services/relabel"
	conditions2 "github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/services/relabel/conditions"
)

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
