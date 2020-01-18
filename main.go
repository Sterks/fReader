package main

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/services"
	"github.com/Sterks/FReader/services/router"
)

func main() {
	mainRunner()
}

func mainRunner() {
	configPath := "config/config.toml"
	config := config.NewConf()
	toml.DecodeFile(configPath, &config)
	ftpreader := services.New(config)
	f := ftpreader.Start()

	now := time.Now()
	y, m, d := now.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Now()

	// str := "2020-12-25"
	// from, _ := time.Parse(time.RFC3339, str)
	// to := time.Now()
	go f.FirstChecherRegions()

	ticker := time.NewTicker(time.Minute * 1)
	go TaskRun(f, from, to, "notifications", ticker)

	ticker2 := time.NewTicker(time.Minute * 2)
	go TaskRun(f, from, to, "protocols", ticker2)

	// ticker := time.NewTicker(time.Second * 1)
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go TaskRun(f, from, to, "notifications", ticker, &wg)
	// wg.Wait()

	// var wg sync.WaitGroup
	// wg.Add(2)

	// ticker := time.NewTicker(time.Minute * 3)
	// go TaskRun(f, from, to, "notifications", ticker, &wg)

	// ticker2 := time.NewTicker(time.Minute * 7)
	// go TaskRun(f, from, to, "protocols", ticker2, &wg)
	// wg.Wait()
	r := router.New(config)
	r.StartWebServer()
}

// TaskRun - метод для организации таск
func TaskRun(f *services.FtpReader, from time.Time, to time.Time, tt string, ticker *time.Ticker) {
	defer ticker.Stop()
	for {
		<-ticker.C
		f.TaskManager(from, to, tt)
	}

}
