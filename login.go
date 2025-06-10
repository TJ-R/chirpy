package main

import (
	"chirpy/internal/auth"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type loginRequest struct {
	Password string `json:"password"`
	Email string `json:"email"` 
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	login_req := loginRequest{}

	err := decoder.Decode(&login_req)
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

	logged_user := User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, 200, logged_user)
}
