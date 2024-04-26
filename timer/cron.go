package timer

import (
	"github.com/robfig/cron/v3"
	"poem-bot/global"
	"poem-bot/timer/dingtalk"
)

func CronTask() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 10 * * ?", dingtalk.SendPoem)
	c.AddFunc("0 0 13 * * ?", dingtalk.SendPoem)
	c.AddFunc("0 0 17 * * ?", dingtalk.SendPoem)
	c.AddFunc("0 0 19 * * ?", dingtalk.SendPoem)
	c.AddFunc("0 0 22 * * ?", dingtalk.SendPoem)
	c.Start()

	global.LOG.Info("Cron task running")
}
