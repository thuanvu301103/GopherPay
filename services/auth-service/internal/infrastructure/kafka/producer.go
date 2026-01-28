package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(cfg Config, topic string) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaBrokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			// Other configuration
			RequiredAcks: kafka.RequiredAcks(cfg.KafkaRequiredAcks),
			Async:        cfg.KafkaAsync,

			WriteTimeout: 10 * time.Second,
            ReadTimeout:  10 * time.Second,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	err := p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Printf("Could not write message to Kafka: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}