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

	r.Get("/post/create", createPostHandler)
	r.Post("/post/create", createPostHandler)

	// Handler using Context middleware to perform repetitive tasks in middleware
	r.Route("/post/{id}", func(r chi.Router) {
		r.Use(middlewares.PostCtx)

		// post object fetched in the PostCtx middleware. Handlers can perform its own specific set of actions.
		r.Get("/", getPostHandler)

		r.Get("/edit", editPostHandler)
		r.Post("/edit", editPostHandler)

		r.Delete("/delete", deletePostHandler)

		// r.Post("/", postPostHandler) // Handle Post request
		// r.Put("/", putPostHandler) // Handle Put request
	})

	// r.Get("/post/{id}", GetPostHandler) // Regular approach

	log.Fatal(http.ListenAndServe(":3000", r))
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := make(map[string]interface{})

	//post part

	if r.Method == "POST" {

		r.ParseForm()

		title := r.PostForm.Get("title")
		description := r.PostForm.Get("description")

		stmt := "insert into posts (title, description) VALUES ($1, $2)"

		q, err := database.DBConn.Prepare(stmt)

		if err != nil {
			log.Println(err)
		}

		res, err := q.Exec(title, description)

		if err != nil {
			log.Println(err)
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 1 {
			ctx["success"] = "Post successfully created."
		}

		log.Println("Rows affected - ", rowsAffected)

	}

	// Get part
	t, _ := template.ParseFiles("templates/pages/post_form.html")

	err := t.Execute(w, ctx)

	if err != nil {
		log.Println("Error in template execution.")
	}

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

// Helper function
func catchErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value("post").(model.Post)

	stmt := "DELETE from posts where id=$1"

	query, err := database.DBConn.Prepare(stmt)

	catchErr(err)

	res, err := query.Exec(post.Id)

	catchErr(err)

	rowsAffected, _ := res.RowsAffected()

	log.Println("Rows affected: ", rowsAffected)

	w.Write([]byte("Record deleted successfully."))

}

func editPostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := make(map[string]interface{})
	p := r.Context().Value("post")

	post := p.(model.Post)

	// Post Part

	if r.Method == "POST" {

		r.ParseForm()

		title := r.PostForm.Get("title")
		description := r.PostForm.Get("description")

		stmt := "UPDATE posts set title=$1, description=$2 where id=$3"

		query, err := database.DBConn.Prepare(stmt)

		if err != nil {
			log.Println(err)
		}

		res, err := query.Exec(title, description, post.Id)

		if err != nil {
			log.Println(err)
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 1 {
			ctx["success"] = "Post successully updated."
		}

		log.Println(rowsAffected)

	}

	// Get Part

	// Load template
	t, _ := template.ParseFiles("templates/pages/post_form.html")

	ctx["post"] = post
	err := t.Execute(w, ctx)

	if err != nil {
		log.Println("Error in tpl execution", err)
	}

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
