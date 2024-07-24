package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	blogDB "github.com/vpaulo/go_blog/internal/db"
	_ "modernc.org/sqlite"
)

const port = ":4000"

var db *sql.DB
var mux *http.ServeMux

const blogDataFolder = "./.blog"
const dbPath string = blogDataFolder + "/sqlite-database.db"

func connectDB() {
	var err error

	blogDB.CreateDataFolder(blogDataFolder, dbPath)

	db, err = blogDB.OpenDB(dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	blogDB.CreateTable(db)
}

func Start() {
	connectDB()
	mux = http.NewServeMux()

	mux.HandleFunc("GET /", getAllArticles)
	mux.HandleFunc("GET /article", newArticle)
	mux.HandleFunc("POST /article", createArticle)
	mux.Handle("GET /article/{articleId}", articleCtx(http.HandlerFunc(getArticle)))
	mux.Handle("PUT /article/{articleId}", articleCtx(http.HandlerFunc(updateArticle)))
	mux.Handle("DELETE /article/{articleId}", articleCtx(http.HandlerFunc(deleteArticle)))
	mux.Handle("GET /article/{articleId}/edit", articleCtx(http.HandlerFunc(editArticle)))

	log.Printf("Listening on %s...", port)
	mwMux := changeMethod(mux)
	if err := http.ListenAndServe(port, mwMux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func changeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func articleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleId := r.PathValue("articleId")
		article, err := blogDB.GetArticle(db, articleId)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := blogDB.GetAllArticles(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Articles: ", articles)

	t, _ := template.ParseFiles("web/views/base.html", "web/views/index.html")
	err = t.Execute(w, articles)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func newArticle(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("web/views/base.html", "web/views/new.html")
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	err := blogDB.CreateArticle(db, title, content)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*blogDB.Article)

	t, _ := template.ParseFiles("web/views/base.html", "web/views/article.html")
	err := t.Execute(w, article)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*blogDB.Article)
	title := r.FormValue("title")
	content := r.FormValue("content")

	err := blogDB.UpdateArticle(db, strconv.Itoa(article.ID), title, content)
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Redirect(w, r, fmt.Sprintf("/article/%d", article.ID), http.StatusFound)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*blogDB.Article)

	err := blogDB.DeleteArticle(db, strconv.Itoa(article.ID))
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func editArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*blogDB.Article)

	t, _ := template.ParseFiles("web/views/base.html", "web/views/edit.html")
	err := t.Execute(w, article)
	if err != nil {
		log.Fatal(err.Error())
	}
}
