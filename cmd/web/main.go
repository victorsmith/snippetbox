package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	// Import internal package
	"snippetbox.victorsmith.dev/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

// for a given DSN.

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "http network address")
	dsn := flag.String("dsn", "root:snippet@/snippetbox?parseTime=true", "Database Connection String")
	// Must call parse, or default value will be used
	flag.Parse()

	// Setup loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// openDB is a helper function which connects our application to a mysql db
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// Closes db connection pool before main exits
	defer db.Close()

	// initialize template cache
	cache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: cache,
	}

	// Establish server so that we can add a logger (instead of using ListenAndServe)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
