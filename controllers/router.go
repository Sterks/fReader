package controllers

import (
	config2 "github.com/Sterks/fReader/config"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//WebServer ...
type WebServer struct {
	config *config2.Config
	router *mux.Router
}

// New ...
func New(conf *config2.Config) *WebServer {
	return &WebServer{
		router: mux.NewRouter(),
		config: conf,
	}
}

// StartWebServer ...
func (w *WebServer) StartWebServer() {

	r := mux.NewRouter()

	r.HandleFunc("/", w.HomeController)

	srv := &http.Server{
		Handler: r,
		Addr:    w.config.MainSettings.ServerConnect,
	}
	log.Fatal(srv.ListenAndServe())
}