package main

import (
	"errors"
	"log"
	"os"

	app "github.com/AlyonaAg/script-checker-server/internal/app/checker-server"
	"github.com/AlyonaAg/script-checker-server/internal/config"
	scriptsdb "github.com/AlyonaAg/script-checker-server/internal/db/scripts"
	"github.com/AlyonaAg/script-checker-server/internal/kafka/consumer"
	"github.com/AlyonaAg/script-checker-server/internal/kafka/producer"
	virustotal "github.com/AlyonaAg/script-checker-server/internal/service/vt"
)

func main() {
	repo, err := scriptsdb.NewRepository()
	if err != nil {
		log.Fatalf("repo error: %v", err)
	}

	apiKey, ok := os.LookupEnv("VT_API_KEY")
	if !ok {
		log.Fatal(errors.New("empty api key"))
	}

	vt, err := virustotal.NewVirusTotal(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	originalScriptTopic, err := config.GetValue(config.OriginalScriptTopic)
	resultTopic, err := config.GetValue(config.ResultTopic)
	brokers, err := config.GetValue(config.Brokers)
	retryMax, err := config.GetValue(config.RetryMax)

	producer, err := producer.NewSyncProducer(
		originalScriptTopic.(string),
		brokers.(string),
		int(retryMax.(int64)),
	)

	consumer, err := consumer.NewConsumer(
		resultTopic.(string),
		brokers.(string),
		repo,
		vt,
	)

	if _, err := os.Stat(".temp"); os.IsNotExist(err) {
		if err := os.Mkdir(".temp", os.ModePerm); err != nil {
			log.Fatalf("create temp dir error: %v", err)
		}
	}

	consumer.Receiver()

	s := app.NewCheckerServer(repo, producer)
	s.Start()
}
