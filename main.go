package main

import (
	"poem-bot/core"
	"poem-bot/global"
	"poem-bot/timer"
)

func main() {
	core.Run()
	timer.CronTask()
	global.LOG.Info("poem-bot started successfully")

	// 阻塞主线程
	select {}
}
