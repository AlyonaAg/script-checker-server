package checkerserver

import (
	"log"

	"github.com/AlyonaAg/script-checker-server/internal/model"
	"github.com/gin-gonic/gin"
)

type repo interface {
	CreateScript(script model.Script) (int64, error)
}

type Implementation struct {
	router *gin.Engine
	repo   repo
}

func NewCheckerServer(repo repo) *Implementation {
	return &Implementation{
		router: gin.Default(),
		repo:   repo,
	}
}

func (i *Implementation) Start() error {
	i.configureRouter()

	log.Print("starting server.")

	return i.router.Run(":4567")
}

func (i *Implementation) configureRouter() {
	i.router.POST("v1/scripts", i.CreateScript())
}
