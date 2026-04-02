package connection

import (
	"context"
	"fmt"

	"github.com/koer/koer-module/pkg/config"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaProducer(cfg config.KafkaProducerConfig) *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: w}
}

func (p *KafkaProducer) Publish(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("publishing kafka message: %w", err)
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("closing kafka producer: %w", err)
	}
	return nil
}

func NewKafkaConsumer(cfg config.KafkaConsumerConfig) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})
	return &KafkaConsumer{reader: r}
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("reading kafka message: %w", err)
	}
	return msg, nil
}

func (c *KafkaConsumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("closing kafka consumer: %w", err)
	}
	return nil
}
