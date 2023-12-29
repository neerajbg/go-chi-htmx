package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/user-info", userInfoHandler)

	http.ListenAndServe(":3000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	ctx := make(map[string]string)

	ctx["Name"] = "Neeraj"

	t, _ := template.ParseFiles("templates/index.html")

	err := t.Execute(w, ctx)

	if err != nil {
		log.Println("Erro in template execution.")
	}

}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User info from API server."))
}
