package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/thuanvu301103/auth-service/internal/config"
)

type Producer struct {
	Writer *kafka.Writer
}

func parseRequiredAcks(s string) kafka.RequiredAcks {
	switch s {
	case "all":
		return kafka.RequireAll
	case "one":
		return kafka.RequireOne
	case "none":
		return kafka.RequireNone
	default:
		return kafka.RequireAll
	}
}

func NewProducer(cfg config.Config) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.KafkaBrokers...),
			Balancer: &kafka.LeastBytes{},
			// Other configuration
			RequiredAcks: kafka.RequiredAcks(parseRequiredAcks(cfg.KafkaRequiredAcks)),
			Async:        cfg.KafkaAsync,

			WriteTimeout: time.Duration(cfg.KafkaTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.KafkaTimeout) * time.Second,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	err := p.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}
