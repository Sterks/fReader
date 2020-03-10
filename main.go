package main

import (
	"github.com/Sterks/fReader/controllers"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sterks/fReader/config"
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

	//now := time.Now()

	 str := "2020-03-06"
	 from, _ := time.Parse(time.RFC3339, str)
	 to := time.Now()

	//y, m, d := now.Date()
	//from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	//to := time.Now()
	f.FirstChecherRegions()
	go func() {
		go f.TaskManager(from, to, "protocols", config2)
		go f.TaskManager(from, to, "notifications", config2)
		go gocron.Every(1).Minute().Do(f.FirstChecherRegions, f)
		go gocron.Every(uint64(config2.Tasks.Notifications)).Minutes().Do(f.TaskManager, from, to, "notifications", config2)
		go gocron.Every(uint64(config2.Tasks.Protocols)).Minutes().Do(f.TaskManager, from, to, "protocols", config2)

		<-gocron.Start()
	}()
	s := controllers.NewWebServer(config2, f.Db)
	s.StartWebServer()
}
