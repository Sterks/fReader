package controllers

import (
	database "github.com/Sterks/Pp.Common.Db/db"
	config2 "github.com/Sterks/fReader/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//WebServer ...
type WebServer struct {
	config *config2.Config
	router *mux.Router
	db *database.Database
}

func NewWebServer(config *config2.Config, db *database.Database) *WebServer {
	return &WebServer{config: config, router: nil, db: db}
}

func (w *WebServer) StartWebServer() {

	r := mux.NewRouter()
	r.HandleFunc("/", w.HomeController)

	srv := &http.Server{
		Handler: r,
		Addr:    w.config.MainSettings.ServerConnect,
	}
	log.Fatal(srv.ListenAndServe())
}