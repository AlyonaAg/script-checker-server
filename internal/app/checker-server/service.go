package checkerserver

import (
	"log"

	"github.com/AlyonaAg/script-checker-server/internal/kafka/producer"
	"github.com/AlyonaAg/script-checker-server/internal/model"
	"github.com/gin-gonic/gin"
)

type repo interface {
	CreateScript(script model.Script) (int64, error)
}

type Implementation struct {
	router   *gin.Engine
	repo     repo
	producer *producer.Producer
}

func NewCheckerServer(repo repo, producer *producer.Producer) *Implementation {
	return &Implementation{
		router:   gin.Default(),
		repo:     repo,
		producer: producer,
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
