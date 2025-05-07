package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerResetServerHits)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerServerHits)
	mux.HandleFunc("POST /api/validate_chirp", cfg.HandlerValidateChirp)

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
	w.Header().Add("Content-Type", "text/html")
	o := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())
	if _, err := w.Write([]byte(o)); err != nil {
		return
	}
}

func (cfg *apiConfig) HandlerResetServerHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	o := returnVals{
		CleanedBody: cleaned,
	}

	b, err := json.Marshal(o)
	if err != nil {
		log.Printf("error marshaling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
