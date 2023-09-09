package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/logger"
	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/util"

	"github.com/pdylanross/kube-resource-relabel-webhook/pkg/services/mutation"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
)

// WebhookHandlers encapsulates properties required to setup webhook http handlers.
type WebhookHandlers struct {
	logger  *slog.Logger
	mutator *mutation.Mutator
}

func NewWebhookHandlers(logger *slog.Logger, mutator *mutation.Mutator) *WebhookHandlers {
	return &WebhookHandlers{
		logger:  logger.With(slog.String("handler", "webhook")),
		mutator: mutator,
	}
}

// SetupWebhookHandlers registers webhook related http handlers.
func (wh *WebhookHandlers) SetupWebhookHandlers(router *gin.Engine) error {
	setupHealthHandlers(router)

	router.POST("/webhook/mutate", wh.MutateHandler)

	return nil
}

func (wh *WebhookHandlers) MutateHandler(context *gin.Context) {
	body, err := configReader(context.Request)
	if err != nil {
		wh.logger.Warn("error fetching request body", slog.String("error", err.Error()))

		// todo: admission response error
		httpErrorStatus(context, err, 400)
		return
	}

	if len(body) == 0 {
		httpErrorStatus(context, fmt.Errorf("empty request body"), 400)
		return
	}

	review, err := mutation.NewAdmissionReview(body)
	if err != nil {
		httpErrorStatus(context, fmt.Errorf("error deserializing review: %s", err), 400)
		return
	}

	wh.mutator.Mutate(review)

	logger.LogIf(wh.logger, slog.LevelDebug, func(l *slog.Logger) {
		marshalled, err := json.Marshal(review.ToSerializeable())
		util.ErrCheck(err) // this should never fail

		l.DebugContext(context, "webhook response",
			slog.Group("response",
				slog.String("review", string(marshalled)),
			),
		)
	})

	context.JSON(review.GetStatus(), review.ToSerializeable())
}

// MaxRequestBodyBytes represents the max size of Kubernetes objects we read. Kubernetes allows a 2x
// buffer on the max etcd size
// (https://github.com/kubernetes/kubernetes/blob/0afa569499d480df4977568454a50790891860f5/staging/src/k8s.io/apiserver/pkg/server/config.go#L362).
// We allow an additional 2x buffer, as it is still fairly cheap (6mb)
// Taken from https://github.com/istio/istio/commit/6ca5055a4db6695ef5504eabdfde3799f2ea91fd
const MaxRequestBodyBytes = int64(6 * 1024 * 1024)

// configReader is reads an HTTP request, imposing size restrictions aligned with Kubernetes limits.
func configReader(req *http.Request) ([]byte, error) {
	defer req.Body.Close()
	lr := &io.LimitedReader{
		R: req.Body,
		N: MaxRequestBodyBytes + 1,
	}
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if lr.N <= 0 {
		return nil, errors.NewRequestEntityTooLargeError(fmt.Sprintf("limit is %d", MaxRequestBodyBytes))
	}
	return data, nil
}
