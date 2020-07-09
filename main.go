package main

import (
	"github.com/BurntSushi/toml"
	"github.com/Sterks/fReader/config"
	"github.com/Sterks/fReader/router"
	"github.com/Sterks/fReader/services"
	"github.com/jasonlvhit/gocron"
)

func main() {
	mainRunner()
	// secondRunner()
}

//TODO Изменять время каждый день

func mainRunner() {
	configPath := "config/config.prod.toml"
	config2 := config.NewConf()
	_, _ = toml.DecodeFile(configPath, &config2)
	ftpReader223 := services.NewFtpReader223(config2)
	f223 := ftpReader223.Start223(config2)
	f223.Start223(config2)
	f44 := services.NewFtpReader44(config2)
	f44.Start44(config2)

	go func() {
		go f223.FirstChecherRegions()
		go f44.FirstChecherRegions()
		go gocron.Every(1).Minute().Do(f44.FirstChecherRegions)
		go gocron.Every(1).Minute().Do(f223.FirstChecherRegions)

		go gocron.Every(uint64(config2.Tasks.Notifications44)).Minutes().Do(f44.TaskManager, "notifications44", config2)
		go gocron.Every(uint64(config2.Tasks.Notifications223)).Minutes().Do(f223.TaskManager, "notifications223", config2)
		go gocron.Every(uint64(config2.Tasks.Protocols44)).Minutes().Do(f44.TaskManager, "protocols44", config2)
		go gocron.Every(uint64(config2.Tasks.Protocols223)).Minutes().Do(f223.TaskManager, "protocols223", config2)
		<-gocron.Start()
	}()

	s := router.NewWebServer(config2, f44.Db)
	s.Start()
}
