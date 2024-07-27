package db

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type Article struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

func CreateDataFolder(blogDataFolder string, dbPath string) {
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		e := os.MkdirAll(blogDataFolder, 0700) // Create data folder
		if e != nil {
			log.Fatal(e)
		}
		CreateDB(dbPath) // Create db file
	}
}

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateDB(dbPath string) {
	log.Printf("Creating %s...", dbPath)
	file, err := os.Create(dbPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Printf("%s created", dbPath)
}

func CreateTable(db *sql.DB) {
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
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("articles table created")
}

func GetAllArticles(db *sql.DB) ([]*Article, error) {
	query, err := db.Prepare("select id, title, content from articles")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	result, err := query.Query()
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, 0)
	for result.Next() {
		data := new(Article)
		err := result.Scan(
			&data.ID,
			&data.Title,
			&data.Content,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, data)
	}

	return articles, nil
}

func CreateArticle(db *sql.DB, title string, content string) error {
	query, err := db.Prepare("insert into articles(title,content) values (?,?)")
	if err != nil {
		return err
	}
	defer query.Close()

	article := &Article{
		Title:   title,
		Content: template.HTML(content),
	}

	_, err = query.Exec(article.Title, article.Content)
	if err != nil {
		return err
	}

	return nil
}

func GetArticle(db *sql.DB, articleId string) (*Article, error) {
	query, err := db.Prepare("select id, title, content from articles where id = ?")
	if err != nil {
		return nil, err
	}
	defer query.Close()

	result := query.QueryRow(articleId)
	data := new(Article)

	err = result.Scan(&data.ID, &data.Title, &data.Content)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UpdateArticle(db *sql.DB, articleId string, title string, content string) error {
	query, err := db.Prepare("update articles set (title, content) = (?,?) where id=?")
	if err != nil {
		return err
	}
	defer query.Close()

	article := &Article{
		Title:   title,
		Content: template.HTML(content),
	}

	_, err = query.Exec(article.Title, article.Content, articleId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteArticle(db *sql.DB, articleId string) error {
	query, err := db.Prepare("delete from articles where id=?")
	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Exec(articleId)
	if err != nil {
		return err
	}

	return nil
}
