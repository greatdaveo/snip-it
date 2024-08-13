package main

import (
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// To create a file server which serves files out of the "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// To use the mux.Handle() function to register the file server as the handler
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on: 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}