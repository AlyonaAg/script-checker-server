package checkerserver

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Success bool `json:"success"`
}

func newErrorResponse(c *gin.Context, httpError int, err error) {
	c.JSON(httpError, errorResponse{Error: err.Error()})
}

func newSuccessResponse() *successResponse {
	return &successResponse{
		Success: true,
	}
}
