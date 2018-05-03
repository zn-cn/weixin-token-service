package main

import (
	"conf"
	"daemon"
	"util"

	"github.com/jasonlvhit/gocron"
)

func main() {

	mainLog := util.GetLogger("/app/log/main.txt", "[DEBUG]")

	// 初始化配置
	configErr := conf.InitConfig()
	if configErr != nil {
		mainLog.Fatalln(configErr)
	}
	// Database connection
	h, err := daemon.DBInit()

	if err != nil {
		mainLog.Fatalln("redis初始化失败")
	}

	// 先初始化access_token, 再初始化ticket，由于定时任务要隔60分钟再开始，所以先update一次
	h.UpdateToekn()
	h.UpdateTicket()

	// 定时任务
	s := gocron.NewScheduler()
	// s.Every(60).Minutes().Do(h.UpdateToekn)
	// s.Every(60).Minutes().Do(h.UpdateTicket)
	s.Every(60).Minutes().Do(h.Update)
	<-s.Start()
}
