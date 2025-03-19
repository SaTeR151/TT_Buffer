package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/SaTeR151/TT_Buffer/internal/config"
	"github.com/SaTeR151/TT_Buffer/internal/models"
	"github.com/SaTeR151/TT_Buffer/internal/repository/redis"
	logger "github.com/sirupsen/logrus"
)

type ServiceInterface interface {
	SendFact(ctx context.Context, config config.SFConfig)
	Insert(ctx context.Context, fact models.Fact) error
}

type ServiceStruct struct {
	redisClient redis.RedisInteface
}

func New(redisClient *redis.RedisStruct) *ServiceStruct {
	service := &ServiceStruct{redisClient: redisClient}
	return service
}

func (s *ServiceStruct) SendFact(ctx context.Context, config config.SFConfig) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fact, err := s.redisClient.GetRandomFact(ctx)
			if err == redis.KeyNilError {
				time.Sleep(1 * time.Second)
				break
			}
			client := &http.Client{
				Timeout: 6 * time.Second,
			}
			JSONFact, err := json.Marshal(fact)
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			reqForSave, _ := http.NewRequest("POST", config.SaveFactURL, bytes.NewBuffer(JSONFact))
			reqForSave.Header.Add("Authorization", "Bearer"+config.Token)

			for {
				res, err := client.Do(reqForSave)
				if err != nil {
					logger.Error(err)
					panic(err)
				}
				defer res.Body.Close()

				if res.StatusCode != http.StatusOK {
					logger.Error(err)
				}

				var checkFact models.CheckFact
				checkFact.PeriodEnd = fact.PeriodEnd
				checkFact.PeriodKeep = fact.PeriodKeep
				checkFact.IndicatorToMOID = fact.IndicatorToMOID
				checkFact.PeriodStart = fact.PeriodStart

				JSONCheckFact, err := json.Marshal(checkFact)
				if err != nil {
					logger.Error(err)
					panic(err)
				}
				reqCheckFact, _ := http.NewRequest("POST", config.GetFactsURL, bytes.NewBuffer(JSONCheckFact))
				reqCheckFact.Header.Add("Authorization", "Bearer"+config.Token)
				res, err = client.Do(reqCheckFact)
				if err != nil {
					logger.Error(err)
					panic(err)
				}
				defer res.Body.Close()
				if res.StatusCode != http.StatusOK {
					logger.Error(err)
				}

				var buf bytes.Buffer
				_, err = buf.ReadFrom(res.Body)
				if err != nil {
					panic(err)
				}

				m := make(map[string]interface{})

				err = json.Unmarshal(buf.Bytes(), &m)
				if err != nil {
					panic(err)
				}

				msg := m["DATA"].(map[string]interface{})
				ma := msg["rows"].([]interface{})
				GetFact := ma[0].(map[string]string)
				if fact.IndicatorToMOID == GetFact["indicator_to_mo_id"] {
					s.redisClient.DeleteFact(ctx, fact.IndicatorToMOID)
					break
				}
			}
		}
	}
}

func (s *ServiceStruct) Insert(ctx context.Context, fact models.Fact) error {
	return s.redisClient.Set(ctx, fact)
}
