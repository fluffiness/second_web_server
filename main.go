package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	// The handler checks if the path of the segment of the url has the prefix "/app/"
	// If so, http.StripPrefix strips "/app" from the path, and concatenates it to filepathRoot
	// The file at the resulting path is served.
	// If the path points to a directory, the index.html file in the directory is served.
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz/", handlerReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
