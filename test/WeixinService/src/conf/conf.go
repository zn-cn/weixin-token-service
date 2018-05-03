package conf

import (
	"errors"
	"os"
	"util"
)

var (
	// Conf : holds the global app config.
	Conf config
)

type config struct {
	ReleaseMode bool
	LogLevel    string
	App         app
	Redis       redis
	Weixin      weixin
	Env         string
}

type app struct {
	Host string
	Port string
}

type redis struct {
	Host string
	Port string
}

type weixin struct {
	AppID     string
	AppSecret string
}

func init() {
}

// InitConfig initializes the app configuration by first setting defaults,
// then overriding settings from the app config file, then overriding
// It returns an error if any.
func InitConfig(configFile string) error {

	res := util.TestOSENV([]string{"APP_HOST", "APP_PORT", "ENV", "APPID", "APP_SECRET"})
	if !res {
		return errors.New("Lack environment variables")
	}
	Conf.App = app{
		Host: os.Getenv("APP_HOST"),
		Port: os.Getenv("APP_PORT"),
	}

	Conf.Weixin = weixin{
		AppID:     os.Getenv("APPID"),
		AppSecret: os.Getenv("APP_SECRET"),
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

	Conf.Env = os.Getenv("ENV")
	if Conf.Env == "prod" {
		Conf.LogLevel = "info"
	} else {
		Conf.LogLevel = "debug"
	}

	return nil
}
