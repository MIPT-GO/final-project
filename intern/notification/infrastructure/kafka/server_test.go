package kafka

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

type fakeReader struct {
	messages chan kafka.Message

	mu        sync.Mutex
	committed []kafka.Message
}

func newFakeReader(buffer int) *fakeReader {
	return &fakeReader{
		messages: make(chan kafka.Message, buffer),
	}
}

func (reader *fakeReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	select {
	case <-ctx.Done():
		return kafka.Message{}, ctx.Err()
	case msg := <-reader.messages:
		return msg, nil
	}
}

func (reader *fakeReader) CommitMessages(_ context.Context, messages ...kafka.Message) error {
	reader.mu.Lock()
	defer reader.mu.Unlock()
	reader.committed = append(reader.committed, messages...)
	return nil
}

func (reader *fakeReader) Close() error {
	return nil
}

func (reader *fakeReader) committedCount() int {
	reader.mu.Lock()
	defer reader.mu.Unlock()
	return len(reader.committed)
}

func TestServer_CommitsOnSuccess(t *testing.T) {
	reader := newFakeReader(10)

	var handled int32
	handler := func(_ []byte) error {
		atomic.AddInt32(&handled, 1)
		return nil
	}

	server := &Server{
		reader:  reader,
		handler: handler,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- server.Serve(ctx) }()

	reader.messages <- kafka.Message{Value: []byte(`{"event":"booking_created"}`)}
	reader.messages <- kafka.Message{Value: []byte(`{"event":"booking_created"}`)}
	reader.messages <- kafka.Message{Value: []byte(`{"event":"booking_created"}`)}

	assert.Eventually(t, func() bool {
		return reader.committedCount() == 3
	}, 2*time.Second, 10*time.Millisecond)

	cancel()
	assert.NoError(t, <-done)
	assert.Equal(t, int32(3), atomic.LoadInt32(&handled))
}

func TestServer_DoesNotCommitOnHandlerError(t *testing.T) {
	reader := newFakeReader(10)

	errHandler := errors.New("handler failed")
	var handled int32
	handler := func(value []byte) error {
		atomic.AddInt32(&handled, 1)
		if string(value) == "bad" {
			return errHandler
		}
		return nil
	}

	server := &Server{
		reader:  reader,
		handler: handler,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- server.Serve(ctx) }()

	reader.messages <- kafka.Message{Value: []byte("ok")}
	reader.messages <- kafka.Message{Value: []byte("bad")}

	assert.Eventually(t, func() bool {
		return reader.committedCount() == 1
	}, 2*time.Second, 10*time.Millisecond)

	cancel()
	assert.NoError(t, <-done)
	assert.Equal(t, int32(2), atomic.LoadInt32(&handled))
}

func TestServer_StopsOnContextCancel(t *testing.T) {
	reader := newFakeReader(1)
	server := &Server{
		reader:  reader,
		handler: func(_ []byte) error { return nil },
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- server.Serve(ctx) }()

	cancel()
	assert.NoError(t, <-done)
}
