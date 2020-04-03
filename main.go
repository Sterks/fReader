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
}



func mainRunner() {
	configPath := "config/config.prod.toml"
	config2 := config.NewConf()
	toml.DecodeFile(configPath, &config2)
	ftpreader := services.New(config2)
	f := ftpreader.Start(config2)


	go func() {
		go f.TaskManager("protocols", config2)
		go f.TaskManager("notifications", config2)
		go gocron.Every(1).Minute().Do(f.FirstChecherRegions, f)
		go gocron.Every(uint64(config2.Tasks.Notifications)).Minutes().Do(f.TaskManager, "notifications", config2)
		go gocron.Every(uint64(config2.Tasks.Protocols)).Minutes().Do(f.TaskManager,"protocols", config2)

		<-gocron.Start()
	}()
	s := router.NewWebServer(config2, f.Db)
	s.Start()
}
