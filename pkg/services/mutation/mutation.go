package mutation

import (
	"log/slog"

	"github.com/pdylanross/kube-resource-relabel-webhook/v1alpha1/pkg/services/relabel"
)

type Mutator struct {
	rules  relabel.Relabeler
	logger *slog.Logger
}

func NewMutator(rules relabel.Relabeler, logger *slog.Logger) *Mutator {
	return &Mutator{
		rules:  rules,
		logger: logger.With(slog.String("component", "mutator")),
	}
}

func (m *Mutator) Mutate(req AdmissionReview) {
	obj, err := req.GetObject()
	if err != nil {
		req.SetError(err)
		return
	}

	patchOps := m.rules.Evaluate(obj)
	if len(patchOps) > 0 {
		if err := req.SetPatches(patchOps); err != nil {
			req.SetError(err)
			return
		}
	}
}
