package main

import (
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/services"
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
	// str := "2020-12-25"
	// from, _ := time.Parse(time.RFC3339, str)
	to := time.Now()
	// ticker := time.NewTicker(time.Second * 1)
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go TaskRun(f, from, to, "notifications", ticker, &wg)
	// wg.Wait()

	f.FirstChecherRegions()

	var wg sync.WaitGroup
	wg.Add(2)

	ticker := time.NewTicker(time.Minute * 3)
	go TaskRun(f, from, to, "notifications", ticker, &wg)

	ticker2 := time.NewTicker(time.Minute * 7)
	go TaskRun(f, from, to, "protocols", ticker2, &wg)
	wg.Wait()
}

// TaskRun - метод для организации таск
func TaskRun(f *services.FtpReader, from time.Time, to time.Time, tt string, ticker *time.Ticker, wg *sync.WaitGroup) {
	defer ticker.Stop()
	defer wg.Done()
	for {
		<-ticker.C
		f.TaskManager(from, to, tt)
	}

}
