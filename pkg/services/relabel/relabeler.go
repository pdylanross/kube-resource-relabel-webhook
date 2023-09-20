package relabel

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"gomodules.xyz/jsonpatch/v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Relabeler is the main entrypoint into the relabeling process
// it keeps a runtime optimized set of rules to relabel objects.
type Relabeler interface {
	Evaluate(originalRawObject []byte, obj metaV1.Object) ([]jsonpatch.Operation, error)
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
func (r *relabeler) Evaluate(originalRawObject []byte, obj metaV1.Object) ([]jsonpatch.Operation, error) {
	changed := false
	for _, rule := range r.rules {
		changed = rule.Evaluate(obj, r.logger) || changed
	}

	if !changed {
		return []jsonpatch.Operation{}, nil
	}

	modifiedJSON, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("error serializing modified object: %s", err.Error())
	}

	patch, err := jsonpatch.CreatePatch(originalRawObject, modifiedJSON)
	if err != nil {
		return nil, fmt.Errorf("error creating patch: %s", err.Error())
	}

	return patch, nil
}
