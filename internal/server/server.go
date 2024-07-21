package server

import (
	"fmt"
	"log"
	"net/http"
)

const port = ":4000"

var mux *http.ServeMux

func Start() {
	mux = http.NewServeMux()

	mux.HandleFunc("GET /", getAllArticles)
	mux.HandleFunc("GET /article", newArticle)
	mux.HandleFunc("POST /article", createArticle)
	mux.HandleFunc("GET /article/{articleId}", getArticle)
	mux.HandleFunc("PUT /article/{articleId}", updateArticle)
	mux.HandleFunc("DELETE /article/{articleId}", deleteArticle)
	mux.HandleFunc("GET /article/{articleId}/edit", editArticle)

	log.Printf("Listening on %s...", port)
	if err := http.ListenAndServe(port, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get all articles\n")
}

func newArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "new article\n")
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create article\n")
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "get article\n")
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "update article\n")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "delete article\n")
}

func editArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "edit article\n")
}
