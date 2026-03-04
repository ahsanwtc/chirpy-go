package main

import (
	"log"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func health(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Header().Add("Content-Type", "text/plain; charset=utf-8")
	response.Write([]byte(http.StatusText(http.StatusOK)))
}

func main () {
	filePathRoot := "."
	fs := http.Dir(filePathRoot)
	port := "8080"


	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(fs)))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/healthz", health)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}
