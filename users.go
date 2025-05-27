package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type newUserRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"create_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	new_user_req := newUserRequest{}

	err := decoder.Decode(&new_user_req)
	if err != nil {
		respondWithError(w, "Error when decoding", 500)
		return
	}

	user, err := cfg.dbQueries.CreateUser(context.Background(), new_user_req.Email)
	if err != nil {
		respondWithError(w, "Error when creating user", 500)	
		return
	}

	respondWithJSON(w, 201, user)
}
