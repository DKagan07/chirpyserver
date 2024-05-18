package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const port = "8080"
	const filePathRoot = "."
	cfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/*",
		cfg.middlewarMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))),
	)
	mux.Handle("/assets/logo", http.FileServer(http.Dir("./assets/logo.png")))
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits))); err != nil {
			log.Fatalf("writing number of hits: %v", err)
		}
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits = 0
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Starting server
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Fatalf("writing ok: %v", err)
	}
}

func (cfg *apiConfig) middlewarMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
