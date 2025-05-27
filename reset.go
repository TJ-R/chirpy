package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, "Not Authorized", 403)
		return
	}

	err := cfg.dbQueries.DeleteUsers(context.Background())

	if err != nil {
		respondWithError(w, "Error when deleting users", 500)
		return
	}

	respondWithJSON(w, 200, "OK")
}
