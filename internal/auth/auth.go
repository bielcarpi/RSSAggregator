package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	// Get the API key from the Authorization header
	apiKey := headers.Get("Authorization")
	log.Println(headers)
	if apiKey == "" {
		return "", errors.New("missing API key")
	}

	// Check if the API key is valid
	key := strings.Split(apiKey, "Bearer ")
	if len(key) != 2 {
		return "", errors.New("invalid API key")
	}

	return key[1], nil
}
