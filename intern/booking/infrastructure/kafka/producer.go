package kafka

import (
	"context"
	"encoding/json"
	"time"

	"final-project/intern/booking/application/interfaces"

	"github.com/segmentio/kafka-go"
)

// AsyncProducer writes messages to Kafka in a fire-and-forget fashion.
type AsyncProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string, topic string) interfaces.Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        topic,
		Async:        false, // we'll write in separate goroutine ourselves
		RequiredAcks: kafka.RequireOne,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 5 * time.Second,
	}
	return &AsyncProducer{writer: w}
}

func (p *AsyncProducer) Send(message []byte) {
	// Fire-and-forget: launch goroutine, don't block caller
	go func(msg []byte) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.writer.WriteMessages(ctx, kafka.Message{Value: msg})
		// errors intentionally ignored (fire-and-forget). If desired, this could log via a logger passed at creation.
	}(message)
}

// convenience helper to send arbitrary structs as JSON
func (p *AsyncProducer) SendStruct(v any) {
	b, _ := json.Marshal(v)
	p.Send(b)
}
