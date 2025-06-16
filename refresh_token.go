package main

import (
	"chirpy/internal/auth"
	"context"
	"log"
	"net/http"
	"time"
)

type TokenCheckResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println(err)
		respondWithError(w, "Error retrieving bearer token", 500)
		return
	}

	refresh_token, err := cfg.dbQueries.GetRefreshToken(context.Background(), token)
	
	if err != nil { 
		log.Println(err)
		respondWithError(w, "Failed to find token", 401)
		return
	}

	if refresh_token.ExpiresAt.Before(time.Now()) {
		log.Println(err)
		respondWithError(w, "Refresh Token has expired", 401)
		return
	}

	if refresh_token.RevokedAt.Valid {
		log.Println(err)
		respondWithError(w, "Refresh Token has been revoked", 401)
		return
	}

	jwt_token, err := auth.MakeJWT(refresh_token.UserID, cfg.secret, time.Second * 3600)
	
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to create new JWT", 500)
		return
	}

	tokenCheckResponse := TokenCheckResponse{
		Token: jwt_token,
	}

	respondWithJSON(w, 200, tokenCheckResponse)
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearer_token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println(err)
		respondWithError(w, "Error retrieving bearer token", 500)
		return
	}

	err = cfg.dbQueries.RevokeToken(context.Background(), bearer_token)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to revoke refresh token", 500)
		return
	}

	respondWithJSON(w, 204, nil)
}
