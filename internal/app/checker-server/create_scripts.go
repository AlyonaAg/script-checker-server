package checkerserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func (i *Implementation) CreateScript() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, &successResponse{
			Success: true,
		})
	}
}