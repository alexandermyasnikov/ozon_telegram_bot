package kafkareader

import (
	"context"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
)

type MsgCallback = func(ctx context.Context, key, value []byte)

type KafkaReader struct {
	reader *kafka.Reader
}

func New(addr, topic, groupID string) *KafkaReader {
	return &KafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{ //nolint:exhaustruct
			Brokers:                []string{addr},
			GroupID:                groupID,
			Topic:                  topic,
			AllowAutoTopicCreation: true,
		}),
	}
}

func (k *KafkaReader) Read(ctx context.Context, callback MsgCallback) {
	for {
		msg, err := k.reader.ReadMessage(ctx)
		if err != nil {
			logger.Infof("can not read message: %v", err)

			break
		}

		logger.Infof("KafkaReader: read: %v/%v/%v: %s = %s", msg.Topic, msg.Partition, msg.Offset,
			string(msg.Key), string(msg.Value))

		ctx, span := otel.Tracer("KafkaReader").Start(ctx, "Read")
		callback(ctx, msg.Key, msg.Value)
		span.End()
	}
}

func (k *KafkaReader) Close() error {
	err := k.reader.Close()

	return errors.Wrap(err, "KafkaReader.Close")
}
