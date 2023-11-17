package main

import (
	"encoding/json"
	"github.com/bielcarpi/RSSAggregator/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiConfig *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	// Parse the request body
	type RequestBody struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	// Parse the request body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil || reqBody.FeedId == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create the user in the database
	feedFollow, err := apiConfig.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    reqBody.FeedId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, DBFeedFollowToFeedFollow(feedFollow))
}

func (apiConfig *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollows, err := apiConfig.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get feed follows from user")
		return
	}

	respondWithJSON(w, http.StatusCreated, DBFeedFollowsToFeedFollows(feedFollows))
}

func (apiConfig *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Feed Follow ID invalid")
		return
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete feed follow from user")
		return
	}

	respondWithJSON(w, http.StatusCreated, struct{}{})
}
