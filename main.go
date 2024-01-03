package main

import (
	"html/template"
	"log"
	"net/http"

	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/neerajbg/chi-htmx/model"

	_ "github.com/lib/pq"
)

var DBConn *sql.DB

func init() {
	dsn := "host=localhost port=5432 user=postgres password=neeraj dbname=chi-htmx-demo sslmode=disable"

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Println("Error in DB connection", err)
	}

	DBConn = db
	log.Println("Database connection successful.")

}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/user-info", userInfoHandler)

	r.Get("/posts", postHandler)

	log.Fatal(http.ListenAndServe(":3000", r))
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

func postHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("OK")
	var posts []model.Post

	sql := "select * from posts"

	rows, err := DBConn.Query(sql)

	defer rows.Close()

	if err != nil {
		log.Println("error in DB execution", err)
	}

	for rows.Next() {
		data := model.Post{}

		err := rows.Scan(&data.Id, &data.Title)

		if err != nil {
			log.Println(err)
		}

		posts = append(posts, data)
	}

	log.Println(posts)

	ctx := make(map[string]interface{})

	ctx["posts"] = posts
	ctx["heading"] = "Article List"

	t, _ := template.ParseFiles("templates/pages/post.html")

	err = t.Execute(w, ctx)

	if err != nil {
		log.Println("Error in tpl execution", err)
	}

}
