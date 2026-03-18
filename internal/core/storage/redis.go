package storage

import (
	"context"
	"encoding/json"
	"main/internal/features/models"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_CONN"),
		Password:    "",
		DB:          0,
		ReadTimeout: 2 * time.Second,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

func Set(ctx context.Context, rdb *redis.Client, key string, value *models.TaskModel, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, b, ttl).Err()
}

func SetTasks(ctx context.Context, rdb *redis.Client, key string, value []models.TaskModel, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, b, ttl).Err()
}

func GetTasks(ctx context.Context, rdb *redis.Client, key string) ([]models.TaskModel, error) {
	r, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var models []models.TaskModel
	err = json.Unmarshal([]byte(r), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func Get(ctx context.Context, rdb *redis.Client, key string) (*models.TaskModel, error) {
	r, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var model models.TaskModel
	err = json.Unmarshal([]byte(r), &model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func Delete(ctx context.Context, rdb *redis.Client, key string) error {
	return rdb.Del(ctx, key).Err()
}

func BlockToken(ctx context.Context, rdb *redis.Client, token string, ttl time.Duration) error {
	return rdb.Set(ctx, token, "", ttl).Err()
}

func IsTokenBlocked(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	result, err := rdb.Exists(ctx, token).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
