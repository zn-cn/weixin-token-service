package main

import (
	"cache"
	"conf"
	"restful/routes"
	"util"

	"github.com/jasonlvhit/gocron"
)

func main() {

	mainLog := util.GetLogger("/app/log/main.txt", "[DEBUG]")

	// InitConfig
	if err := conf.InitConfig(); err != nil {
		mainLog.Fatalln(err)
	}
	go func() {
		// 首先更新一次
		cache.UpdateToekn()
		cache.UpdateTicket()

		// 定时任务
		s := gocron.NewScheduler()
		s.Every(60).Minutes().Do(cache.UpdateToekn)
		s.Every(60).Minutes().Do(cache.UpdateTicket)
		<-s.Start()
	}()

	routes.Init()

}
