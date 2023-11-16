package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithJSON sends a JSON response with the given status code and payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	// Marshal the payload to JSON
	response, err := json.Marshal(payload)
	if err != nil {
		// If there is an error, send an internal server error
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", message)
	}

	// To Marshal the payload to JSON
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{Error: message})
}
