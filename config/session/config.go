package session

import (
	"strings"
	"QA-System/config/config"
	WeJHSDK  "github.com/zjutjh/WeJH-SDK"
)

type driver string

const (
	Memory driver = "memory"
	Redis  driver = "redis"
)

var defaultName = "wejh-session"



func getConfig() WeJHSDK.SessionInfoConfig {

	wc := WeJHSDK.SessionInfoConfig{}
	wc.Driver = string(Memory)
	if config.Config.IsSet("session.driver") {
		wc.Driver = strings.ToLower(config.Config.GetString("session.driver"))
	}

	wc.Name = defaultName
	if config.Config.IsSet("session.name") {
		wc.Name = strings.ToLower(config.Config.GetString("session.name"))
	}

	wc.SecretKey = "secret"

	wc.RedisConfig = getRedisConfig()

	return wc
}

func getRedisConfig() WeJHSDK.RedisInfoConfig {
	Info := WeJHSDK.RedisInfoConfig{
		Host:     "localhost",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	if config.Config.IsSet("redis.host") {
		Info.Host = config.Config.GetString("redis.host")
	}
	if config.Config.IsSet("redis.port") {
		Info.Port = config.Config.GetString("redis.port")
	}
	if config.Config.IsSet("redis.db") {
		Info.DB = config.Config.GetInt("redis.db")
	}
	if config.Config.IsSet("redis.pass") {
		Info.Password = config.Config.GetString("redis.pass")
	}
	return Info
}
