package config

import (
	"fmt"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
	actions2 "github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel/actions"
)

func (act *RelabelConfigRuleAction) GetConcrete() (relabel.ActionConfig, error) {
	switch act.Type {
	case "ensure-annotation":
		var ret EnsureAnnotationAction
		err := act.Value.Decode(&ret)
		return &ret, err
	case "ensure-label":
		var ret EnsureLabelAction
		err := act.Value.Decode(&ret)
		return &ret, err
	}

	return nil, fmt.Errorf("unknown action type %s", act.Type)
}

type EnsureAnnotationAction map[string]string

func (e *EnsureAnnotationAction) ToAction() relabel.Action {
	return actions2.NewEnsureAnnotationAction(*e)
}

type EnsureLabelAction map[string]string

func (e *EnsureLabelAction) ToAction() relabel.Action {
	return actions2.NewEnsureLabelAction(*e)
}
