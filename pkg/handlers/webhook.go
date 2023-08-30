package handlers

import "github.com/gin-gonic/gin"

type WebhookHandlers struct {
}

func (rh *WebhookHandlers) SetupWebhookHandlers(_ *gin.Engine) error {
	return nil
}
