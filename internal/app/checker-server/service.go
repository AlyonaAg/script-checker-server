package checkerserver

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Implementation struct {
	router    *gin.Engine
}

func NewCheckerServer() *Implementation{
	return &Implementation{
		router:    gin.Default(),
	}
}

func (i *Implementation) Start() error {
	i.configureRouter()

	log.Print("starting server.")

	return i.router.Run()
}

func (i *Implementation) configureRouter() {
	i.router.POST("v1/scripts", i.CreateScript())
}

type successResponse struct {
	Success bool `json:"success"`
}
