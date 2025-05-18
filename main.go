package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type metricsHandler struct {
}

func main() {
	apiCfg := apiConfig {
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/",  apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerHealthCheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidation)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)


	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}


