package checkerserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlyonaAg/script-checker-server/internal/kafka/producer"
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
			if len(s) == 0 {
				continue
			}
			id, err := i.repo.CreateScript(model.Script{
				URL:    r.URL,
				Script: s,
			})
			if err != nil {
				log.Printf("ERROR: url %s, script %d: %v", r.URL, index+1, err)
				continue
			}

			msg, err := json.Marshal(producer.Message{
				ID:     id,
				Script: s,
			})
			if err != nil {
				log.Printf("ERROR: marshal: %v", err)
				continue
			}
			if err := i.producer.Send(msg); err != nil {
				log.Printf("ERROR: kafka send: %v", err)
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
