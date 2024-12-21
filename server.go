package main

import (
	"fmt"
	"net/http"

	"github.com/adilanuarbek/emojitranslate/controllers"
	"github.com/adilanuarbek/emojitranslate/templates"
	"github.com/adilanuarbek/emojitranslate/views"
	"github.com/go-chi/chi/v5"
)

func main() {
	handlers := controllers.Handlers{}
	r := chi.NewRouter()
	handlers.Templates.Index = views.Must(views.ParseFS(templates.FS, "home.html"))
	r.Get("/", handlers.IndexHandler)

	fmt.Println("Starting the server on 8080...")
	http.ListenAndServe(":8080", r)
}
