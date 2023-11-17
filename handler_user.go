package main

import (
	"encoding/json"
	"github.com/bielcarpi/RSSAggregator/internal/db"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiConfig *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	type RequestBody struct {
		Username string `json:"username"`
	}

	// Parse the request body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.Username == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create the user in the database
	usr, err := apiConfig.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Username,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, DBUserToUser(usr))
}

func (apiConfig *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, http.StatusOK, DBUserToUser(user))
}
