package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	const port = "8080"
	const filePathRoot = "."

	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.MiddlewareMetricsInc(
			http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))),
		),
	)
	mux.HandleFunc("GET /healthz", handlerHealthz)
	mux.HandleFunc("GET /metrics", cfg.HandlerServerHits)
	mux.HandleFunc("POST /reset", cfg.HandlerResetServerHits)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Starting server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		return
	}
}

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) HandlerServerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	o := fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())
	if _, err := w.Write([]byte(o)); err != nil {
		return
	}
}

func (cfg *apiConfig) HandlerResetServerHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
