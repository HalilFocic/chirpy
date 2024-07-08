package main

import (
	"fmt"
	"github.com/HalilFocic/chirpy/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiConfig := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiConfig.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/api/reset", apiConfig.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /api/validate_chirp", HandleValidateChrip)
	mux.HandleFunc("POST /api/chirps", apiConfig.handleAddChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.handleGetChirpById)
	mux.HandleFunc("GET /api/users", apiConfig.handleGetUsers)
	mux.HandleFunc("POST /api/users", apiConfig.handleAddUser)
	mux.HandleFunc("POST /api/login", apiConfig.handleLogin)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)
	w.Write([]byte(resp))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}
