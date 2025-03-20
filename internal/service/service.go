package service

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/SaTeR151/TT_Buffer/internal/config"
	"github.com/SaTeR151/TT_Buffer/internal/models"
	"github.com/SaTeR151/TT_Buffer/internal/repository/redis"
	logger "github.com/sirupsen/logrus"
)

var SaveFactError = fmt.Errorf("saving fact error")
var FactError = fmt.Errorf("fact is wrong. fact id: ")

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

// SendFactToSave - функция для отправки факта на сохрание
func SendFactToSave(client *http.Client, fact models.Fact, config config.SFConfig) error {

	buf := &bytes.Buffer{}
	mpw := multipart.NewWriter(buf)

	/// Парсинг данных в form-data
	periodStartW, err := mpw.CreateFormField("period_start")
	if err != nil {
		return err
	}
	_, err = periodStartW.Write([]byte(fact.PeriodStart))
	if err != nil {
		return err
	}

	PeriodEndW, err := mpw.CreateFormField("period_end")
	if err != nil {
		return err
	}
	_, err = PeriodEndW.Write([]byte(fact.PeriodEnd))
	if err != nil {
		return err
	}

	PeriodKeyW, err := mpw.CreateFormField("period_key")
	if err != nil {
		return err
	}
	_, err = PeriodKeyW.Write([]byte(fact.PeriodKey))
	if err != nil {
		return err
	}

	IndicatorToMOIDW, err := mpw.CreateFormField("indicator_to_mo_id")
	if err != nil {
		return err
	}
	_, err = IndicatorToMOIDW.Write([]byte(fact.IndicatorToMOID))
	if err != nil {
		return err
	}

	IndicatorToMOFactIDW, err := mpw.CreateFormField("indicator_to_mo_fact_id")
	if err != nil {
		return err
	}
	_, err = IndicatorToMOFactIDW.Write([]byte(fact.IndicatorToMOFactID))
	if err != nil {
		return err
	}

	ValueW, err := mpw.CreateFormField("value")
	if err != nil {
		return err
	}
	_, err = ValueW.Write([]byte(fact.Value))
	if err != nil {
		return err
	}

	FactTimeW, err := mpw.CreateFormField("fact_time")
	if err != nil {
		return err
	}
	_, err = FactTimeW.Write([]byte(fact.FactTime))
	if err != nil {
		return err
	}

	IsPlanW, err := mpw.CreateFormField("is_plan")
	if err != nil {
		return err
	}
	_, err = IsPlanW.Write([]byte(fact.IsPlan))
	if err != nil {
		return err
	}

	AuthuserIdW, err := mpw.CreateFormField("auth_user_id")
	if err != nil {
		return err
	}
	_, err = AuthuserIdW.Write([]byte(fact.AuthuserId))
	if err != nil {
		return err
	}

	CommentW, err := mpw.CreateFormField("comment")
	if err != nil {
		return err
	}
	_, err = CommentW.Write([]byte(fact.Comment))
	if err != nil {
		return err
	}
	mpw.Close()

	// отправка запроса на сохранение факта
	reqForSave, _ := http.NewRequest("POST", config.SaveFactURL, buf)
	reqForSave.Header.Add("Content-Type", mpw.FormDataContentType())
	reqForSave.Header.Add("Authorization", "Bearer "+config.Token)
	res, err := client.Do(reqForSave)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusBadRequest {
			return FactError
		}
		return SaveFactError
	}
	return nil
}

// функция для отправки данных на сохрание факта, а также его удаление после успешного сохранения
func (s *ServiceStruct) SendFact(ctx context.Context, config config.SFConfig) {
	client := &http.Client{
		Timeout: 6 * time.Second,
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:

			fact, err := s.redisClient.GetRandomFact(ctx)
			if err == redis.KeyNilError {
				time.Sleep(1 * time.Second)
			} else {
				logger.Info("sending fact")
				err = SendFactToSave(client, fact, config)
				if err == nil {
					s.redisClient.DeleteFact(ctx, fact.IndicatorToMOID)
				} else {
					logger.Error(err.Error() + fact.IndicatorToMOID)
				}
			}

		}
	}
}

func (s *ServiceStruct) Insert(ctx context.Context, fact models.Fact) error {
	return s.redisClient.Set(ctx, fact)
}
