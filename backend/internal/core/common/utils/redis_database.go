package utils

import (
	"context"
	"encoding/json"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisDatabase RedisDatabaseType

type RedisDatabaseType struct {
	Client redis.Client
	Prefix string
}

func ConnectToRedis(ctx context.Context) error {
	log.Info(ctx).Msg("Connecting to Redis")
	redisConfig := config.ApplicationConfig.Data.Redis

	redisOptions := &redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	}
	client := redis.NewClient(redisOptions)

	RedisDatabase = RedisDatabaseType{
		Client: *client,
		Prefix: redisConfig.Prefix,
	}

	cmd := client.Conn().Ping(ctx)
	err := cmd.Err()

	if err == nil {
		log.Info(ctx).Msg("Redis connected!")
	}

	return err
}

func (redisDatabase RedisDatabaseType) Get(ctx context.Context, key string) *redis.StringCmd {
	completeKey := redisDatabase.getCompleteKey(key)
	return redisDatabase.Client.Get(ctx, completeKey)
}

func (redisDatabase RedisDatabaseType) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	completeKey := redisDatabase.getCompleteKey(key)
	return redisDatabase.Client.Set(ctx, completeKey, value, expiration).Err()
}

func (redisDatabase RedisDatabaseType) Del(ctx context.Context, key string) error {
	completeKey := redisDatabase.getCompleteKey(key)
	return redisDatabase.Client.Del(ctx, completeKey).Err()
}

func (redisDatabase RedisDatabaseType) GetStruct(ctx context.Context, key string, value any) error {
	completeKey := redisDatabase.getCompleteKey(key)
	cmd := redisDatabase.Get(ctx, completeKey)

	err := cmd.Err()
	if err != nil {
		return err
	}

	jsonData, _ := cmd.Bytes()

	err = json.Unmarshal(jsonData, value)
	if err != nil {
		return err
	}

	return nil
}

func (redisDatabase RedisDatabaseType) SetStruct(ctx context.Context, key string, value any, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	json := string(jsonData)

	completeKey := redisDatabase.getCompleteKey(key)
	result := redisDatabase.Set(ctx, completeKey, json, expiration)
	return result
}

func (redisDatabase RedisDatabaseType) getCompleteKey(key string) string {
	return redisDatabase.Prefix + constants.COLON + key
}
