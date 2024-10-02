package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippet-box/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// To define an application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// To make the SnippetModel object available to the handlers
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template // templateCache field
}

func main() {
	// Tpo create DB Connection Pool
	dsn := flag.String("dsn", "web:webpassword@/snippetbox?parseTime=true", "MySQL database")
	// To define a new command-line flag with the name 'addr',
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// To create a logger for writing information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// To create a logger for writing error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// To close the connection pool before the main() function exists
	defer db.Close()

	// To initialize a new template cache
	templateCache, err := newTemplateCache("ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// To initialize a mysql.SnippetModel instance & add i the application dependencies
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache, // templateCache
	}

	// Swapped the route declarations to use the application struct's method
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// To create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// To use the mux.Handle() function to register the file server as the handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// To call the ListenAndServe method on the new http.Server struct
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN
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
