package router

import (
	"log"
	"net/http"

	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/controllers"
	"github.com/gorilla/mux"
)

//WebServer ...
type WebServer struct {
	config *config.Config
	router *mux.Router
}

// New ...
func New(conf *config.Config) *WebServer {
	return &WebServer{
		router: mux.NewRouter(),
		config: conf,
	}
}

// StartWebServer ...
func (w *WebServer) StartWebServer() {

	r := mux.NewRouter()

	r.HandleFunc("/", controllers.HomeController)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:3000",
	}
	log.Fatal(srv.ListenAndServe())
}
