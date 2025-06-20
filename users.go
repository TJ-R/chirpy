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
	"os"
)

type UserRequest struct {
	Password string `json:"password"`
	Email string `json:"email"`
}

type ChirpyRedData struct {
	UserID uuid.UUID `json:"user_id"`	
}

type ChirpyRedRequest struct {
	Event string `json:"event"`
	Data  ChirpyRedData `json:"data"` 
}

type User struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	IsChirpyRed bool     `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	new_user_req := UserRequest{}

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
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, 201, new_user)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	bearer_token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error getting token", 401)
		return
	}

	userID, err := auth.ValidateJWT(bearer_token, cfg.secret)
	
	if err != nil {
		log.Println(err)
		respondWithError(w, "Token is invalid", 401)
		return
	}

	decoder := json.NewDecoder(req.Body)
	updateUserRequest := UserRequest{}

	err = decoder.Decode(&updateUserRequest)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when decoding", 500)
		return
	}

	hashedPassword, err := auth.HashPassword(updateUserRequest.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when hashing password", 500)
		return
	}

	user, err := cfg.dbQueries.UpdateUser(context.Background(), database.UpdateUserParams{
		Email: updateUserRequest.Email,
		HashedPassword: hashedPassword,
		ID: userID,
	})

	updatedUser := User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}	

	respondWithJSON(w, 200, updatedUser)
}

func (cfg *apiConfig) handlerUpdateUserChirpyRed(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)

	if err != nil {
		log.Println(err)
		respondWithError(w, "Authorization failed", 401)
		return
	}

	if apiKey != os.Getenv("POLKA_KEY") {
		log.Println("Authorization failed")
		respondWithError(w, "Permission denied: Not a valid API key", 401)
		return
	}
	
	decoder := json.NewDecoder(req.Body)
	chirpyRedRequest := ChirpyRedRequest{}

	err = decoder.Decode(&chirpyRedRequest)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when decoding", 500)
		return
	}

	if chirpyRedRequest.Event != "user.upgraded" {
		respondWithJSON(w, 204, "")
		return
	}

	_, err = cfg.dbQueries.UpdateUserChirpyRed(context.Background(), chirpyRedRequest.Data.UserID)
	if err != nil {
		log.Println(err)
		respondWithError(w, "User not found", 404)
		return
	}

	respondWithJSON(w, 204, "")
	return
}
