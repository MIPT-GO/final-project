package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type reader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, messages ...kafka.Message) error
	Close() error
}
