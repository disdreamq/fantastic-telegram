package kafka

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	consumer *kafka.Reader
}

func NewConsumer(consumer *kafka.Reader) *Consumer {
	return &Consumer{
		consumer: consumer,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	logger := log.Ctx(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.consumer.ReadMessage(ctx)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("Failed to update posts in outbox db")

			}
			logger.Info().
				Str("Topic", msg.Topic).
				Int("Partition", msg.Partition).
				Str("Value", string(msg.Value)).
				Msg("Received message")
		}
	}
}
