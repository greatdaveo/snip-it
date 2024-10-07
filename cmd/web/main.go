package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippet-box/pkg/models/mysql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// To define an application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	// To make the SnippetModel object available to the handlers
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template // templateCache field
	users         *mysql.UserModel
}

func main() {
	// To create DB Connection Pool
	dsn := flag.String("dsn", "web:webpassword@/snippetbox?parseTime=true", "MySQL database")
	// To define a command-line flag with the name 'addr',
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	// To define a command-line flag for the session secret
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret Key for session cookies")
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

	// To initialize a new session manager that expires after 12 hours
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// The application dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		// To initialize a mysql.SnippetModel instance & add the application dependencies
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache, // templateCache
		users:         &mysql.UserModel{DB: db},
	}

	// To initialize a tls.Config struct to hold the non-default TLS settings
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
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
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// To add Idle, Read and Write timeouts to the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// To call the ListenAndServe method on the new http.Server struct
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem") // To start the HTTPS server with self-signed TLS certificate
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
