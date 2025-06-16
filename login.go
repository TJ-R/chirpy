package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type loginRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
}

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	login_req := loginRequest{}

	err := decoder.Decode(&login_req)

	// Make sure expiration is not over an hours if it is limit it

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

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour * 1)

	if err != nil {
		log.Println(err)
		respondWithError(w, "There was an issue generating your token", 500)
		return
	}

	refresh_token_str, _ := auth.MakeRefreshToken() 

	refresh_token, err := cfg.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token: refresh_token_str,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, "There was an issue creating you refresh token", 500)
		return
	}

	loginResp := LoginResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
		RefreshToken: refresh_token.Token,
	}

	respondWithJSON(w, 200, loginResp)
}
