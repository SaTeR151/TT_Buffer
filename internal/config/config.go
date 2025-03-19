package config

import (
	"os"
)

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       string
}

type SFConfig struct {
	Token       string
	SaveFactURL string
	GetFactsURL string
}

func GetServerConfig() ServerConfig {
	var config ServerConfig
	config.Port = os.Getenv("SERVER_PORT")
	return config
}

func GetRedisConfig() RedisConfig {
	var config RedisConfig
	config.Host = os.Getenv("REDIS_HOST")
	config.Port = os.Getenv("REDIS_PORT")
	config.Username = os.Getenv("REDIS_USERNAME")
	config.Password = os.Getenv("REDIS_PASSWORD")
	config.DB = os.Getenv("REDIS_DB")
	return config
}

func GetSFConfig() SFConfig {
	var config SFConfig
	config.Token = os.Getenv("TOKEN")
	config.SaveFactURL = os.Getenv("URL_SAVE_FACT")
	config.GetFactsURL = os.Getenv("URL_GETS_FACT")
	return config
}
