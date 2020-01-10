package main

import (
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/services"
	"github.com/patrickmn/go-cache"
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

	c := cache.New(1*time.Hour, 2*time.Hour)

	var wg sync.WaitGroup
	wg.Add(2)
	go f.TaskManager(from, to, "notifications", &wg, c)
	go f.TaskManager(from, to, "protocols", &wg, c)
	wg.Wait()
}
