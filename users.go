package main

import (
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"chirpy/internal/auth"
	"github.com/google/uuid"
)

type newUserRequest struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	new_user_req := newUserRequest{}

	err := decoder.Decode(&new_user_req)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when decoding", 500)
		return
	}

	hashed_password, err := auth.HashPassword(new_user_req.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when hashing password", 500)
		return
	}

	user, err := cfg.dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email: new_user_req.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when creating user", 500)	
		return
	}

	new_user := User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, 201, new_user)
}
