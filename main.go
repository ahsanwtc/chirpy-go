package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type apiConfig struct {
	filerserverHists atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.filerserverHists.Add(1)
		next.ServeHTTP(w, r)
	})
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
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app",
		apiCfg.middlewareMetricsInc(http.FileServer(fs)),
	))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("Hits: %d", apiCfg.filerserverHists.Load())))
	})
	mux.HandleFunc("POST /reset", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		apiCfg.filerserverHists.Store(0)
	})

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}
