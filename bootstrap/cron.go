package bootstrap

import "github.com/robfig/cron/v3"

var c *cron.Cron

func InitializeCron(f func(c *cron.Cron)) {
	c = cron.New()
	f(c)
	c.Start()
}
