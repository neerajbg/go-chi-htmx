package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/neerajbg/chi-htmx/database"
	"github.com/neerajbg/chi-htmx/middlewares"
	"github.com/neerajbg/chi-htmx/model"

	_ "github.com/lib/pq"
)

func init() {
	database.ConnectDB()
}
func main() {
	defer database.DBConn.Close()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/user-info", userInfoHandler)

	r.Get("/posts", postHandler)

	// Handler using Context middleware to perform repetitive tasks in middleware
	r.Route("/post/{id}", func(r chi.Router) {
		r.Use(middlewares.PostCtx)

		// post object fetched in the PostCtx middleware. Handlers can perform its own specific set of actions.
		r.Get("/", getPostHandler)

		// r.Post("/", postPostHandler) // Handle Post request
		// r.Put("/", putPostHandler) // Handle Put request
	})

	// r.Get("/post/{id}", GetPostHandler) // Regular approach

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

func getPostHandler(w http.ResponseWriter, r *http.Request) {

	post := r.Context().Value("post")

	// Load template
	t, _ := template.ParseFiles("templates/pages/post_detail.html")

	ctx := make(map[string]interface{})
	ctx["post"] = post
	err := t.Execute(w, ctx)

	if err != nil {
		log.Println("Error in tpl execution", err)
	}

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("OK")
	var posts []model.Post

	sql := "select * from posts"

	rows, err := database.DBConn.Query(sql)

	// defer rows.Close()

	if err != nil {
		log.Println("error in DB execution", err)
	}

	for rows.Next() {
		data := model.Post{}

		err := rows.Scan(&data.Id, &data.Title, &data.Description)

		if err != nil {
			log.Println(err)
		}

		posts = append(posts, data)
	}

	ctx := make(map[string]interface{})

	ctx["posts"] = posts
	ctx["heading"] = "Article List"

	t, _ := template.ParseFiles("templates/pages/post.html")

	err = t.Execute(w, ctx)

	if err != nil {
		log.Println("Error in tpl execution", err)
	}

}
