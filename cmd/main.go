package main

import (
	"fmt"

	"github.com/SaTeR151/TT_Buffer/internal/config"
	"github.com/SaTeR151/TT_Buffer/internal/handlers"
	logg "github.com/SaTeR151/TT_Buffer/internal/logger"
	"github.com/SaTeR151/TT_Buffer/internal/repository/redis"
	"github.com/SaTeR151/TT_Buffer/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

func main() {
	logg.Init()

	err := godotenv.Load()
	if err != nil {
		logger.Error(err)
		return
	}

	serverConfig := config.GetServerConfig()
	redisConfig := config.GetRedisConfig()
	//SFConfig := config.GetSFConfig()

	redisClient, err := redis.Connect(redisConfig)
	if err != nil {
		logger.Error(err)
		return
	}
	service := service.New(redisClient)
	r := gin.Default()

	r.POST("/facts_to_buffer", handlers.PostFactsToBuffer(service))
	logger.Info("starting server")
	if err := r.Run(":" + serverConfig.Port); err != nil {
		logger.Error(fmt.Sprintf("server starting error: %s\v", err.Error()))
		return
	}
}
