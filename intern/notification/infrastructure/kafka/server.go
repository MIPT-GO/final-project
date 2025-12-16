package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sync"

	"final-project/intern/notification/config"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(value []byte) error

type Server struct {
	reader  reader
	handler MessageHandler
}

func NewServer(cfg *config.Config, handler MessageHandler) *Server {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.KafkaBrokers,
		Topic:    cfg.KafkaTopic,
		GroupID:  cfg.KafkaGroupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &Server{
		reader:  reader,
		handler: handler,
	}
}

func (kafkaServer *Server) Serve(ctx context.Context) error {
	defer kafkaServer.reader.Close()

	slog.Info("start notification kafka consumer")

	processingConcurrency := runtime.GOMAXPROCS(0)
	if processingConcurrency < 1 {
		processingConcurrency = 1
	}

	workQueue := make(chan kafka.Message, 256)
	commitQueue := make(chan kafka.Message, 256)

	var workersWaitGroup sync.WaitGroup
	var commitWaitGroup sync.WaitGroup

	commitWaitGroup.Add(1)
	go func() {
		defer commitWaitGroup.Done()
		for message := range commitQueue {
			if err := kafkaServer.reader.CommitMessages(ctx, message); err != nil {
				slog.Error("failed to commit message", "error", err.Error())
			}
		}
	}()

	for workerIndex := 0; workerIndex < processingConcurrency; workerIndex++ {
		workersWaitGroup.Add(1)
		go func() {
			defer workersWaitGroup.Done()

			for message := range workQueue {
				if err := kafkaServer.handler(message.Value); err != nil {
					slog.Error("handler error", "error", err.Error())
					continue
				}

				commitQueue <- message
			}
		}()
	}

	for {
		message, fetchErr := kafkaServer.reader.FetchMessage(ctx)
		if fetchErr != nil {
			if ctx.Err() != nil {
				close(workQueue)
				workersWaitGroup.Wait()
				close(commitQueue)
				commitWaitGroup.Wait()
				slog.Info("notification kafka consumer stopped")
				return nil
			}

			close(workQueue)
			workersWaitGroup.Wait()
			close(commitQueue)
			commitWaitGroup.Wait()
			return fmt.Errorf("fetch message: %w", fetchErr)
		}

		select {
		case workQueue <- message:
		case <-ctx.Done():
			close(workQueue)
			workersWaitGroup.Wait()
			close(commitQueue)
			commitWaitGroup.Wait()
			slog.Info("notification kafka consumer stopped")
			return nil
		}
	}
}
