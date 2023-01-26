package checkerserver

import (
	"log"
	"net/http"

	"github.com/AlyonaAg/script-checker-server/internal/model"
	"github.com/gin-gonic/gin"
)

func (i *Implementation) CreateScript() gin.HandlerFunc {
	return func(c *gin.Context) {
		var r = &CreateScriptRequest{}
		if err := c.ShouldBindJSON(r); err != nil {
			newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		for index, s := range r.Scripts {
			_, err := i.repo.CreateScript(model.Script{
				URL:    "test",
				Script: s,
			})
			if err != nil {
				log.Printf("ERROR: url %s, script %d: %v", r.URL, index+1, err)
			}
		}

		c.JSON(http.StatusOK, newSuccessResponse())
	}
}

type CreateScriptRequest struct {
	Count   int64    `json:"count"`
	URL     string   `json:"url"`
	Scripts []string `json:"scripts"`
}
