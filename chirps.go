package main

import (
	"net/http"
	"encoding/json"
	"log"
)

type chirpValidationRequest struct {
	Body string `json:"body"`
}

type chirpValidationErrorResponse struct {
	Error string `json:"error"`
}

type chirpValidationResponse struct {
	Valid bool `json:"valid"`
}

func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	validation_req := chirpValidationRequest{}
	err := decoder.Decode(&validation_req)
	if err != nil {
		returnError(w, "Error when decoding", 500)
		return
	}

	if len(validation_req.Body) > 140 {
		returnError(w, "Chirp is too long", 400)
		return
	}

	respBody := chirpValidationResponse{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		returnError(w, "Error Marshalling Json", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func returnError (w http.ResponseWriter, message string, code int) {
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
