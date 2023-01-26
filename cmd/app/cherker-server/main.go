package main

import (
	"log"

	app "github.com/AlyonaAg/script-checker-server/internal/app/checker-server"
	scriptsdb "github.com/AlyonaAg/script-checker-server/internal/db/scripts"
)

func main() {
	repo, err := scriptsdb.NewRepository()
	if err != nil {
		log.Fatalf("repo error: %v", err)
	}

	s := app.NewCheckerServer(repo)
	s.Start()
}
