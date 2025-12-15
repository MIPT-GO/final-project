package kafka

import (
	"context"
	"encoding/json"
	"time"

	"final-project/intern/booking/application/interfaces"

	"github.com/segmentio/kafka-go"
)

type AsyncProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string, topic string) interfaces.Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        topic,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 5 * time.Second,
	}
	return &AsyncProducer{writer: w}
}

func (p *AsyncProducer) Send(message []byte) {
	go func(msg []byte) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.writer.WriteMessages(ctx, kafka.Message{Value: msg})
	}(message)
}

func (p *AsyncProducer) SendStruct(v any) {
	b, _ := json.Marshal(v)
	p.Send(b)
}
