package controllers

import (
	"html/template"
	"net/http"
)

//HomeController ...
func HomeController(w http.ResponseWriter, h *http.Request) {
	type ViewData struct {
		Title string
	}

	data := ViewData{
		Title: "Test",
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))
	tmpl.Execute(w, data)
}
