package controllers

import (
	"html/template"
	"net/http"
)

//HomeController ...
func (s *WebServer) HomeController(w http.ResponseWriter, h *http.Request) {
	type ViewData struct {
		Title string
		LastID int
	}

	last := s.db.LastID()

	data := ViewData{
		Title: "Test",
		LastID: last,
	}


	tmpl := template.Must(template.ParseFiles("views/index.html"))
	_ = tmpl.Execute(w, data)
}
