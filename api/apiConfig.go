package api

import (
	"fmt"
	"log"
	"net/http"
)

type ApiConfig struct {
	FileserverHits int
}

func (cfg *ApiConfig) MiddlewarMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html;")
	adminText := fmt.Sprintf(`
<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.FileserverHits)
	if _, err := w.Write([]byte(adminText)); err != nil {
		log.Fatalf("sending admin html: %v", err)
	}
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits = 0
}
