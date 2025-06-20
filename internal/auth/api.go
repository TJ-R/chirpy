package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Auth Provided")
	}

	if !strings.Contains(authHeader, "ApiKey") {
		return "", errors.New("Auth Headers does not contain ApiKey")
	}

	keys := strings.Fields(authHeader)

	foundApiKey := false
	for _, value := range keys {
		if foundApiKey {
			return value, nil
		}

		if value == "ApiKey" {
			foundApiKey = true
		}
	}

	// Should never reach here
	return "", errors.New("There was an issue retrieving the ApiKey")
}
