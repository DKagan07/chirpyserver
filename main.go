package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"chirpyserver/api"
	database "chirpyserver/dbStuff"
	"chirpyserver/types"
)

func main() {
	const port = "8080"
	const filePathRoot = "."
	const dbPath = "./database.json"
	cfg := api.ApiConfig{
		FileserverHits: 0,
	}

	// will need result of NewDb, but will need to implement it too
	_, err := database.NewDb(dbPath)
	if err != nil {
		log.Fatalf("Couldn't create db connection to %s", dbPath)
		return
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/*",
		cfg.MiddlewarMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))),
	)
	mux.Handle("/assets/logo", http.FileServer(http.Dir("./assets/logo.png")))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("/api/reset", cfg.ResetHandler)
	mux.HandleFunc("POST /api/chirps", postChirpsHandler)
	mux.HandleFunc("GET /api/chirps", getChirpsHandler)

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
	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		log.Fatalf("writing ok: %v", err)
	}
}

func postChirpsHandler(w http.ResponseWriter, r *http.Request) {
	incId := 0

	// decode response
	decoder := json.NewDecoder(r.Body)
	param := types.Chirp{}
	if err := decoder.Decode(&param); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// If chirp is greater than 140 characters
	if len(param.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long")
		return
	}

	// Chirp is correct size
	// need to add to db
	msg := cleanMessage(param.Body)
	incId++
	respondWithJSON(w, http.StatusOK, types.ReturnVals{
		Id:          incId,
		CleanedBody: msg,
	})
}

func getChirpsHandler(w http.ResponseWriter, r *http.Request) {}

// *******************************
// ********** HELPERS ************
// *******************************

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		log.Fatalf("error writing data: %v", err)
	}
}

func cleanMessage(msg string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitWords := strings.Split(msg, " ")
	cleanedMsgArr := []string{}

	for _, word := range splitWords {
		if sliceContains(badWords, strings.ToLower(word)) {
			cleanedMsgArr = append(cleanedMsgArr, "****")
		} else {
			cleanedMsgArr = append(cleanedMsgArr, word)
		}
	}
	return strings.Join(cleanedMsgArr, " ")
}

func sliceContains(slc []string, word string) bool {
	for _, v := range slc {
		if v == word {
			return true
		}
	}
	return false
}
