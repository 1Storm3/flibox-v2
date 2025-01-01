package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Produce(message map[string]interface{}) error {
	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("4"),
		Value: value,
	})
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Printf("Ошибка закрытия Kafka producer: %s", err.Error())
	}
}
