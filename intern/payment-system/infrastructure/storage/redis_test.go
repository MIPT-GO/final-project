package storage_test

import (
	"context"
	"final-project/intern/payment-system/config"
	"final-project/intern/payment-system/domain"
	"final-project/intern/payment-system/infrastructure/storage"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestRedisStorage(t *testing.T) {
	test_env := os.Getenv("TEST_ENV")
	godotenv.Load(test_env)

	config := config.Config{}
	config.LoadFromEnv()

	storage := storage.RedisStorage{}
	ctx := context.Background()

	err := storage.NewStorage(ctx, &config)
	assert.NoError(t, err)

	test_info := domain.PaymentInfo{
		Amount:  "100",
		Message: "test info",
		WebHook: "test-url",
	}

	err = storage.CreateRecord(ctx, "test-key", test_info)
	assert.NoError(t, err)

	t.Cleanup(func() {
		storage.DeleteRecord(ctx, "test-key")
	})

	extracted_info, err := storage.GetRecord(ctx, "test-key")
	assert.NoError(t, err)
	assert.Equal(t, test_info, extracted_info)
}
