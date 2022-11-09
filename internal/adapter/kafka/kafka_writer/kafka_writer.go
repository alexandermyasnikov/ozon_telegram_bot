package kafkawriter

import (
	"context"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func New(addr, topic string) *KafkaWriter {
	return &KafkaWriter{
		writer: &kafka.Writer{ //nolint:exhaustruct
			Addr:                   kafka.TCP(addr),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (k *KafkaWriter) Write(ctx context.Context, key, value []byte) error {
	logger.Infof("kafka.write [%s][%s]", string(key), string(value))

	ctx, span := otel.Tracer("KafkaWriter").Start(ctx, "Write")
	defer span.End()

	err := k.writer.WriteMessages(ctx,
		kafka.Message{ //nolint:exhaustruct
			Key:   key,
			Value: value,
		},
	)

	return errors.Wrap(err, "KafkaWriter.Write")
}

func (k *KafkaWriter) Close() error {
	err := k.writer.Close()

	return errors.Wrap(err, "KafkaWriter.Close")
}
