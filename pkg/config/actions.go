package config

import (
	"fmt"

	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/relabel/actions"

	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/relabel"
)

func (act *RelabelConfigRuleAction) GetConcrete() (relabel.ActionConfig, error) {
	switch act.Type {
	case "ensure-annotation":
		var ret EnsureAnnotationAction
		err := act.Value.Decode(&ret)
		return &ret, err
	}

	return nil, fmt.Errorf("unknown action type %s", act.Type)
}

type EnsureAnnotationAction map[string]string

func (e *EnsureAnnotationAction) ToAction() relabel.Action {
	return actions.NewEnsureAnnotationAction(*e)
}
