package redis

import (
	"context"
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

func (redisClient *RedisStruct) Set(ctx context.Context, fact models.Fact) error {
	if err := redisClient.db.HSet(ctx, fact.IndicatorToMOID, fact, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (redisClient *RedisStruct) GetRandomFact(ctx context.Context) (models.Fact, error) {
	var fact models.Fact
	key := redisClient.db.RandomKey(ctx)
	if key == nil {
		return fact, KeyNilError
	}
	val := redisClient.db.HGetAll(ctx, key.String()).Val()
	fact.PeriodStart = val["PeriodStart"]
	fact.PeriodEnd = val["PeriodEnd"]
	fact.PeriodKeep = val["PeriodKeey"]
	fact.IndicatorToMOID = val["IndicatorToMOID"]
	fact.IndicatorToMOFactID = val["IndicatorToMOFactID"]
	fact.Value = val["Value"]
	fact.FactTime = val["FactTime"]
	fact.IsPlan = val["IsPlan"]
	fact.Comment = val["Comment"]
	return fact, nil
}

func (redisClient *RedisStruct) DeleteFact(ctx context.Context, key string) error {
	err := redisClient.db.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
