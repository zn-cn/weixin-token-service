package conf

import (
	"errors"
	"os"
	"util"
)

var (
	// Conf 全局配置
	Conf config
)

type config struct {
	Redis    redis
	Weixin   weixin
	LogLevel string
	Env      string
}

type redis struct {
	Host string
	Port string
}

type weixin struct {
	AppID     string
	AppSecret string
}

// InitConfig 初始化
func InitConfig() error {

	res := util.TestOSENV([]string{"APPID", "APP_SECRET", "ENV"})
	if !res {
		return errors.New("Lack environment variables")
	}
	redisHost := "redis"
	redisPort := "6379"
	if os.Getenv("REDIS_HOST") != "" {
		redisHost = os.Getenv("REDIS_HOST")
	}
	if os.Getenv("REDIS_PORT") != "" {
		redisPort = os.Getenv("REDIS_PORT")
	}
	Conf.Redis = redis{
		Host: redisHost,
		Port: redisPort,
	}

	Conf.Weixin = weixin{
		AppID:     os.Getenv("APPID"),
		AppSecret: os.Getenv("APP_SECRET"),
	}

	Conf.Env = os.Getenv("ENV")
	if Conf.Env == "prod" {
		Conf.LogLevel = "info"
	} else {
		Conf.LogLevel = "debug"
	}
	return nil
}
