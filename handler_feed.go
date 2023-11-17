package main

import (
	"encoding/json"
	"github.com/bielcarpi/RSSAggregator/internal/db"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiConfig *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	// Parse the request body
	type RequestBody struct {
		Username string `json:"username"`
		URL      string `json:"url"`
	}

	// Parse the request body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.Username == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create the user in the database
	feed, err := apiConfig.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      reqBody.Username,
		Url:       reqBody.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, DBFeedToFeed(feed))
}

func (apiConfig *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not fetch feeds")
		return
	}

	respondWithJSON(w, http.StatusOK, DBFeedsToFeeds(feeds))
}
