package storage

import (
	"context"
	"encoding/json"
	"final-project/intern/payment-system/config"
	"final-project/intern/payment-system/domain"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client  *redis.Client
	timeout int
}

func (storage *RedisStorage) NewStorage(ctx context.Context, config *config.Config) error {
	storage.client = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + strconv.Itoa(config.RedisPort),
		Username: config.RedisUser,
		Password: config.RedisPassword,
		DB:       config.RedisDb,
	})

	err := storage.client.Ping(ctx).Err()

	if err != nil {
		return err
	}

	storage.timeout = config.StorageTimeout

	return nil
}

func (storage *RedisStorage) CreateRecord(ctx context.Context, key string, info domain.PaymentInfo) error {
	infoJson, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = storage.client.Set(ctx, key, infoJson, time.Minute*time.Duration(storage.timeout)).Err()
	return err
}

func (storage *RedisStorage) GetRecord(ctx context.Context, key string) (domain.PaymentInfo, error) {
	value, err := storage.client.Get(ctx, key).Result()
	if err != nil {
		return domain.PaymentInfo{}, err
	}

	var info domain.PaymentInfo
	err = json.Unmarshal([]byte(value), &info)

	if err != nil {
		return domain.PaymentInfo{}, err
	}

	return info, nil
}
