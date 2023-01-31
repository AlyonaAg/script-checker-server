package checkerserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AlyonaAg/script-checker-server/internal/model"
	"github.com/gin-gonic/gin"
)

func (i *Implementation) ListScripts() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {

			page = 0
		}

		page = (page - 1) * limit
		scripts, err := i.repo.ListScripts(model.ListScriptsFilter{
			Limit: int64(limit),
			Page:  int64(page),
		})
		if err != nil {
			log.Printf("ERROR: ListScripts err=%v", err)
		}

		fmt.Println(scripts)

		c.JSON(http.StatusOK, toListScriptsResponse(scripts))
	}
}

type ListScriptsRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"offset"`
}

type ListScriptsResponse struct {
	Scripts []Script `json:"scripts"`
}

type Script struct {
	URL        string `json:"url"`
	Script     string `json:"script"`
	Result     bool   `json:"result"`
	VirusTotal string `json:"virustotal"`
}

func toListScriptsResponse(scripts model.Scripts) *ListScriptsResponse {
	var respScripts []Script
	for _, s := range scripts {
		respScripts = append(respScripts, Script{
			URL:        s.URL,
			Script:     s.Script,
			Result:     s.Result,
			VirusTotal: s.VirusTotal,
		})
	}

	return &ListScriptsResponse{
		Scripts: respScripts,
	}
}
