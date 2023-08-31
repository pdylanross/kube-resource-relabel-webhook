package handlers

import "github.com/gin-gonic/gin"

// WebhookHandlers encapsulates properties required to setup webhook http handlers.
type WebhookHandlers struct {
}

// SetupWebhookHandlers registers webhook related http handlers.
func (rh *WebhookHandlers) SetupWebhookHandlers(router *gin.Engine) error {
	setupHealthHandlers(router)

	return nil
}
