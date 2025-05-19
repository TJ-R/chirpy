package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

type chirpValidationRequest struct {
	Body string `json:"body"`
}

type chirpValidationErrorResponse struct {
	Error string `json:"error"`
}

type chirpValidationResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	validation_req := chirpValidationRequest{}
	err := decoder.Decode(&validation_req)
	if err != nil {
		respondWithError(w, "Error when decoding", 500)
		return
	}

	msg := validation_req.Body
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
	
	respondWithJSON(w, 200, chirpValidationResponse{CleanedBody: cleaned_msg,})

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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, "Error Marshalling Json", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)

}
