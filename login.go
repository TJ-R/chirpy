package main

import (
	"chirpy/internal/auth"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type loginRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoggedUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	login_req := loginRequest{
		ExpiresInSeconds: 3600,
	}

	err := decoder.Decode(&login_req)

	// Make sure expiration is not over an hours if it is limit it
	if login_req.ExpiresInSeconds > 3600 {
		login_req.ExpiresInSeconds = 3600
	}

	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when decoding", 500)
		return
	}

	user, err := cfg.dbQueries.GetUser(context.Background(), login_req.Email)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Cannot find user", 401)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, login_req.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Incorrect Password", 401)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second*time.Duration(login_req.ExpiresInSeconds))

	if err != nil {
		log.Println(err)
		respondWithError(w, "There was an issue generating your token", 500)
		return
	}

	logged_user := LoggedUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}

	respondWithJSON(w, 200, logged_user)
}
