package consumer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/AlyonaAg/script-checker-server/internal/model"
	"github.com/Shopify/sarama"
)

type repo interface {
	GetScript(id int64) (*model.Script, error)
	UpdateResultByID(scriptID int64, result bool) error
	UpdateDangerPercentByID(scriptID int64, dangerPercent float64) error
}

type vt interface {
	GetDengerPercent(path string) (float64, error)
}

type Consumer struct {
	topic    string
	repo     repo
	consumer sarama.PartitionConsumer
	vt       vt
}

func NewConsumer(topic string, broker string, repo repo, vt vt) (*Consumer, error) {
	c := &Consumer{
		topic: topic,
		repo:  repo,
		vt:    vt,
	}

	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	master, err := sarama.NewConsumer([]string{broker}, cfg)
	if err != nil {
		return nil, err
	}

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	c.consumer = consumer

	return c, nil
}

func (c *Consumer) Receiver() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for {
			select {
			case err := <-c.consumer.Errors():
				log.Println(err)
			case msg := <-c.consumer.Messages():
				log.Println("received msg", string(msg.Key), string(msg.Value))
				go c.prepareMsg(msg.Value)
			case <-signals:
				log.Println("Interrupt is detected")
				return
			}
		}
	}()
}

func (c *Consumer) prepareMsg(msg []byte) {
	m := &message{}
	m.Unmarshall(msg)

	// update result in repo
	if err := c.repo.UpdateResultByID(m.ID, m.Result); err != nil {
		log.Printf("prepareMsg: UpdateResultByID err=%v\n", err)
		return
	}

	// get script from repo
	script, err := c.repo.GetScript(m.ID)
	if err != nil {
		log.Printf("prepareMsg: GetScript err=%v\n", err)
		return
	}

	// get script by script url
	resp, err := http.Get(script.Script)
	if err != nil {
		log.Printf("prepareMsg: Get err=%v\n", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("prepareMsg: ReadAll err=%v\n", err)
		return
	}

	path := fmt.Sprintf(".temp/%d.js", script.ID)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("prepareMsg: Create err=%v\n", err)
		return
	}
	defer os.Remove(path)

	_, err = f.Write(body)
	if err != nil {
		log.Printf("prepareMsg: Write err=%v\n", err)
		return
	}

	// virustotal
	danger, err := c.vt.GetDengerPercent(path)
	if err != nil {
		log.Printf("prepareMsg: GetDengerPercent err=%v\n", err)
		return
	}

	if err := c.repo.UpdateDangerPercentByID(script.ID, danger); err != nil {
		log.Printf("prepareMsg: UpdateDangerPercentByID err=%v\n", err)
		return
	}

	return
}
