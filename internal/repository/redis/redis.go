package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SaTeR151/TT_Buffer/internal/config"
	"github.com/SaTeR151/TT_Buffer/internal/models"
	"github.com/redis/go-redis/v9"
)

var KeyNilError = fmt.Errorf("database is empty")

type RedisInteface interface {
	Set(ctx context.Context, fact models.Fact) error
	GetRandomFact(ctx context.Context) (models.Fact, error)
	DeleteFact(ctx context.Context, key string) error
}

type RedisStruct struct {
	db *redis.Client
}

// / подключение к Redis
func Connect(config config.RedisConfig) (*RedisStruct, error) {
	var rClient RedisStruct
	dbNumber, err := strconv.Atoi(config.DB)
	if err != nil {
		return nil, err
	}

	rClient.db = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       dbNumber,
	})

	if err = rClient.db.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &rClient, nil
}

// Сохранение факта в Redis
func (redisClient *RedisStruct) Set(ctx context.Context, fact models.Fact) error {
	json, err := json.Marshal(fact)
	if err != nil {
		return err
	}

	if err := redisClient.db.Set(ctx, fact.IndicatorToMOID, json, 0).Err(); err != nil {
		return err
	}
	return nil
}

// Получение факта из Redis
func (redisClient *RedisStruct) GetRandomFact(ctx context.Context) (models.Fact, error) {
	var buf []byte
	var fact models.Fact
	key := redisClient.db.RandomKey(ctx)
	if key.Err() == redis.Nil {
		return fact, KeyNilError
	}
	err := redisClient.db.Get(ctx, key.Val()).Scan(&buf)
	if err != nil {
		return fact, err
	}
	err = json.Unmarshal(buf, &fact)
	if err != nil {
		return fact, err
	}
	return fact, nil
}

// Удаление факта из Redis
func (redisClient *RedisStruct) DeleteFact(ctx context.Context, key string) error {
	err := redisClient.db.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
