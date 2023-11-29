package bootstrap

import (
	"github.com/robfig/cron/v3"
	"github.com/succko/hera/global"
)

var c *cron.Cron

func InitializeCron() {
	c = cron.New()
	f := global.App.RunConfig.Cron
	if f != nil {
		f(c)
	}
	c.Start()
}
