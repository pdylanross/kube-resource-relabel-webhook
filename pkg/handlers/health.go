package handlers

import "github.com/gin-gonic/gin"

func setupHealthHandlers(router *gin.Engine) {
	router.GET("/livez", func(context *gin.Context) {
		context.Status(200)
	})
}
