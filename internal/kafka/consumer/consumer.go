package consumer

import (
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	topic    string
	consumer sarama.PartitionConsumer
}

func NewConsumer(topic string, broker string, retryMax int) (*Consumer, error) {
	c := &Consumer{
		topic: topic,
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
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-c.consumer.Errors():
				log.Println(err)
			case msg := <-c.consumer.Messages():
				log.Println("Received messages", string(msg.Key), string(msg.Value))
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
}
