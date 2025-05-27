package main

import _ "github.com/lib/pq"

import (
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	"chirpy/internal/database"
	"log"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
}

type metricsHandler struct {
}

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("%v", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig {
		fileserverHits: atomic.Int32{},
		dbQueries: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/",  apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidation)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)


	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}


