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
	Weixin      weixin
	Env         string
}

type app struct {
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
func InitConfig() error {

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

	Conf.Env = os.Getenv("ENV")
	if Conf.Env == "prod" {
		Conf.LogLevel = "info"
	} else {
		Conf.LogLevel = "debug"
	}

	return nil
}
