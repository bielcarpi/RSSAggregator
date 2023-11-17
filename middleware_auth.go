package main

import (
	"github.com/bielcarpi/RSSAggregator/internal/auth"
	"github.com/bielcarpi/RSSAggregator/internal/db"
	"net/http"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

func (apiConfig *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the API key from header
		api, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusForbidden, err.Error())
			return
		}

		// Get user by API key on DB
		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), api)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Call the handler
		handler(w, r, user)
	}
}
