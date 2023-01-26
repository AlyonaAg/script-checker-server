package producer

import (
	"github.com/Shopify/sarama"
)

type Producer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewSyncProducer(topic string, broker string, retryMax int) (*Producer, error) {
	p := &Producer{
		topic: topic,
	}

	cfg := sarama.NewConfig()

	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = retryMax
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{broker}, cfg)
	if err != nil {
		return nil, err
	}

	p.producer = producer

	return p, nil
}

func (p *Producer) Send(data []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(data),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
