package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// To define an application struct to hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog * log.Logger
}

func main() {
	// To define a new command-line flag with the name 'addr',
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	// To create a logger for writing information messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// To create a logger for writing error messages
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
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
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// To call the ListenAndServe method on the new http.Server struct
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}