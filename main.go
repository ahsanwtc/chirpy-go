package main

import (
	"log"
	"net/http"
)

func main () {
	filePathRoot := "."
	fs := http.Dir(filePathRoot)
	port := "8080"


	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(fs))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(fs)))

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}
