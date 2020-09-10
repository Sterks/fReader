package main

import (
	"github.com/BurntSushi/toml"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/router"
	"github.com/Sterks/fReader/services"
	"github.com/robfig/cron/v3"

	_ "net/http/pprof"
)

func main() {
	mainRunner()
	// secondRunner()
}

// TODO Изменять время каждый день

func mainRunner() {
	configPath := "config/config.prod.toml"
	conf := config.NewConf()
	_, _ = toml.DecodeFile(configPath, &conf)

	c := cron.New()
	fz223Notification := services.NewFtpReader223(conf)
	// services.StartService223(fz223Notification, conf, "protocols223")
	_, _ = c.AddFunc("*/30 * * * *", func() { services.StartService223(fz223Notification, conf, "notifications223") })
	_, _ = c.AddFunc("*/60 * * * *", func() { services.StartService223(fz223Notification, conf, "protocols223") })

	fz44Notification := services.NewFtpReader44(conf)
	_, _ = c.AddFunc("*/20 * * * *", func() { services.StartService(fz44Notification, conf, "notifications44") })
	_, _ = c.AddFunc("*/25 * * * *", func() { services.StartService(fz44Notification, conf, "protocols44") })
	c.Start()

	s := router.NewWebServer(conf, fz44Notification.Db)
	s.Start()
}
