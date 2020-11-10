package main

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/router"
	"github.com/Sterks/fReader/services"
	"github.com/robfig/cron/v3"
	_ "net/http/pprof"
)

func main() {

	time.Sleep(30 * time.Second)
	mainRunner()

	// secondRunner()
}

func migrationDB(config *config.Config) {
	migration := services.NewMigration(config)
	migration.ConfigureMigration(config)
}

// TODO Изменять время каждый день

func mainRunner() {
	configPath := ""
	getenv := os.Getenv("APPLICATION")
	if getenv == "production" {
		configPath = "config/config.prod.toml"
	} else {
		configPath = "config/config.toml"
	}
	conf := config.NewConf()
	_, _ = toml.DecodeFile(configPath, &conf)

	//migrationDB(conf)

	c := cron.New()
	fz223Notification := services.NewFtpReader223(conf)
	fz44Notification := services.NewFtpReader44(conf)
	go services.StartService(fz44Notification, conf, "notifications44")
	_, _ = c.AddFunc("*/29 * * * *", func() { services.StartService223(fz223Notification, conf, "notifications223") })
	_, _ = c.AddFunc("*/30 * * * *", func() { services.StartService223(fz223Notification, conf, "protocols223") })

	_, _ = c.AddFunc("*/24 * * * *", func() { services.StartService(fz44Notification, conf, "notifications44") })
	_, _ = c.AddFunc("*/20 * * * *", func() { services.StartService(fz44Notification, conf, "protocols44") })
	c.Start()

	s := router.NewWebServer(conf, fz44Notification.Db)
	s.Start()
}
