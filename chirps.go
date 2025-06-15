package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type chirpRequest struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type chirpValidationErrorResponse struct {
	Error string `json:"error"`
}

type chirpValidationResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	chirp_req := chirpRequest{}
	err := decoder.Decode(&chirp_req)
	if err != nil {
		respondWithError(w, "Error when decoding", 500)
		return
	}

	// Validate token before creating chirp
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when getting token from request", 500)
		return
	}

	userIdFromToken, err := auth.ValidateJWT(token, cfg.secret)

	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when validating token", 401)
		return
	}

	msg := chirp_req.Body
	if len(msg) > 140 {
		respondWithError(w, "Chirp is too long", 400)
		return
	}

	
	words := strings.Split(msg, " ")

	for i, word := range words {
		lowered := strings.ToLower(word)
		if lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			words[i] = "****"	
		}
	}

	cleaned_msg := strings.Join(words, " ")
	
	chirp, err := cfg.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body: cleaned_msg,
		UserID: userIdFromToken,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when creating chirp", 500)
		return
	}

	new_chirp := Chirp {
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: userIdFromToken,
	}
	
	respondWithJSON(w, 201, new_chirp)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(context.Background())
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error when getting chirps", 500)
		return
	}

	var recived_chirps []Chirp
	for _, chirp := range chirps {
		recived_chirps = append(recived_chirps, Chirp {
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})	
	}

	respondWithJSON(w, 200, recived_chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("chirpID")
	log.Println(param)

	id, err := uuid.Parse(param)
	log.Println(id)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to parse chirp id", 500)
	}

	chirp, err := cfg.dbQueries.GetChirp(context.Background(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to find chirp", 404)
		return
	}

	log.Println(chirp)

	recieved_chirp := Chirp {
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, 200, recieved_chirp)
}

func respondWithError (w http.ResponseWriter, message string, code int) {
	errRespBody := chirpValidationErrorResponse {
		Error: message,
	}
	
	log.Println(message)

	dat, err := json.Marshal(errRespBody)
	if err != nil {
		log.Printf("Error marshalling error JSON %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	return
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, "Error Marshalling Json", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)

}
