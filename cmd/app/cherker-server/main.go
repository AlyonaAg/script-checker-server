package main

import (
	app "github.com/AlyonaAg/script-checker-server/internal/app/checker-server"
)

func main() {
	s := app.NewCheckerServer()
	s.Start()
}