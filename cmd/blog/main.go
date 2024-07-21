package main

import (
	"database/sql"
	"log"

	blogDB "github.com/vpaulo/go_blog/internal/db"
	"github.com/vpaulo/go_blog/internal/server"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	var err error
	db, err = blogDB.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	server.Start()
}
