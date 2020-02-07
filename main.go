package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/router"
	"github.com/Sterks/FReader/services"
	"github.com/jasonlvhit/gocron"
)

func main() {
	mainRunner()
}

func mainRunner() {
	configPath := "config/config.prod.toml"
	config := config.NewConf()
	toml.DecodeFile(configPath, &config)
	ftpreader := services.New(config)
	f := ftpreader.Start(config)

	now := time.Now()
	y, m, d := now.Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	to := time.Now()

	// str := "2020-12-25"
	// from, _ := time.Parse(time.RFC3339, str)
	// to := time.Now()

	go f.FirstChecherRegions()

	// gocron.Every(1).Minute().Do(testText, f)

	gocron.Every(1).Minute().Do(f.FirstChecherRegions, f)

	// ticker := time.NewTicker(time.Duration(config.Tasks.Notifications) * time.Minute)
	// ticker := time.NewTicker(time.Second * 1)
	// go TaskRun(f, from, to, "notifications", ticker, config)

	// ticker2 := time.NewTicker(time.Duration(config.Tasks.Protocols) * time.Minute)
	// go TaskRun(f, from, to, "protocols", ticker2, config)

	gocron.Every(uint64(config.Tasks.Notifications)).Minutes().Do(f.TaskManager, from, to, "notifications", config)
	gocron.Every(uint64(config.Tasks.Protocols)).Minutes().Do(f.TaskManager, from, to, "protocols", config)

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
	<-gocron.Start()
	r := router.New(config)
	r.StartWebServer()
}

// TaskRun - метод для организации таск
func TaskRun(f *services.FtpReader, from time.Time, to time.Time, tt string, ticker *time.Ticker, config *config.Config) {
	defer ticker.Stop()
	for {
		<-ticker.C
		f.TaskManager(from, to, tt, config)
	}

}

func testText(b *services.FtpReader) {
	fmt.Println("Тест")
	b.FirstChecherRegions()
}
