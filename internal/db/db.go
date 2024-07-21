package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const blogDataFolder = "./.blog"
const dbPath string = blogDataFolder + "/sqlite-database.db"

func Connect() (*sql.DB, error) {
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		e := os.MkdirAll(blogDataFolder, 0700) // Create data folder
		if e != nil {
			log.Fatal(e)
		}
		CreateDatabase() // Create db file
	}

	var err error
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	createTable(db)

	return db, nil
}

func CreateDatabase() {
	log.Printf("Creating %s...", dbPath)
	file, err := os.Create(dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Printf("%s created", dbPath)
}

func createTable(db *sql.DB) {
	createArticlesTableSQL := `CREATE TABLE if not exists articles (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"content" TEXT
	  );`

	log.Println("Create articles table...")
	statement, err := db.Prepare(createArticlesTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("articles table created")
}
