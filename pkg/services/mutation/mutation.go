package mutation

import (
	"log/slog"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/relabel"
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

	patchOps, err := m.rules.Evaluate(req.GetRawObject(), obj)
	if err != nil {
		req.SetError(err)
		return
	}

	if len(patchOps) > 0 {
		if err := req.SetPatches(patchOps); err != nil {
			req.SetError(err)
			return
		}
	}
}
