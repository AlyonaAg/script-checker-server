package main

import (
	"log"

	app "github.com/AlyonaAg/script-checker-server/internal/app/checker-server"
	"github.com/AlyonaAg/script-checker-server/internal/config"
	scriptsdb "github.com/AlyonaAg/script-checker-server/internal/db/scripts"
	"github.com/AlyonaAg/script-checker-server/internal/kafka/producer"
)

func main() {
	repo, err := scriptsdb.NewRepository()
	if err != nil {
		log.Fatalf("repo error: %v", err)
	}

	originalScriptTopic, err := config.GetValue(config.OriginalScriptTopic)
	brokers, err := config.GetValue(config.Brokers)
	retryMax, err := config.GetValue(config.RetryMax)

	producer, err := producer.NewSyncProducer(
		originalScriptTopic.(string),
		brokers.(string),
		int(retryMax.(int64)),
	)

	s := app.NewCheckerServer(repo, producer)
	s.Start()
}
