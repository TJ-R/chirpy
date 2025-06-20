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
	"sort"
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
	authorID := r.URL.Query().Get("author_id")
	
	// No optional param
	if authorID == "" {
		chirps, err := cfg.dbQueries.GetChirps(context.Background())
		if err != nil {
			log.Println(err)
			respondWithError(w, "Error when getting chirps", 500)
			return
		}

		var recivedChirps []Chirp
		for _, chirp := range chirps {
			recivedChirps = append(recivedChirps, Chirp {
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,
			})	
		}

		sortOrder := r.URL.Query().Get("sort")	
		if sortOrder == "desc" {
			sort.Slice(recivedChirps, func(i, j int) bool { return recivedChirps[i].CreatedAt.After(recivedChirps[j].CreatedAt)})
		}

		respondWithJSON(w, 200, recivedChirps)
		return
	}


	// optional param
	authorUserID, err := uuid.Parse(authorID)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to parse author id", 500)
		return
	}

	chirps, err := cfg.dbQueries.GetChirpsForUser(context.Background(), authorUserID)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to find chirps", 404)
		return
	}

	var recivedChirps []Chirp
	for _, chirp := range chirps {
		recivedChirps = append(recivedChirps, Chirp {
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body, UserID: chirp.UserID,
		})	
	}

	sortOrder := r.URL.Query().Get("sort")	
	log.Println(sortOrder)
	if sortOrder == "desc" {
		sort.Slice(recivedChirps, func(i, j int) bool { return recivedChirps[i].CreatedAt.After(recivedChirps[j].CreatedAt)})
	}

	respondWithJSON(w, 200, recivedChirps)
	return
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("chirpID")

	id, err := uuid.Parse(param)
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

	recieved_chirp := Chirp {
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, 200, recieved_chirp)
}

func (cfg *apiConfig) handlerDeleteChrip(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("chirpID")
	
	id, err := uuid.Parse(param)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to parse chipr id", 500)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Error getting token", 401)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Token is invalid", 401)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(context.Background(), id)

	if chirp.UserID != userID {
		log.Println("UserID does not match chirp's userID")
		respondWithError(w, "Chirp does not belong to user", 403)
		return
	}

	err = cfg.dbQueries.DeleteChirp(context.Background(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, "Failed to find chirp", 404)
		return
	}

	respondWithJSON(w, 204, "")
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
