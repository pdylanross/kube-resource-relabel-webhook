package handlers

import "github.com/gin-gonic/gin"

type HTTPErrorModel struct {
	Status  int
	Message string
}

func httpErrorStatus(c *gin.Context, err error, status int) {
	model := HTTPErrorModel{
		Status:  status,
		Message: err.Error(),
	}

	c.PureJSON(status, model)
}
