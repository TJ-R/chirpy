package auth

import (
	"github.com/google/uuid"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"errors"
	"net/http"
	"strings"
)


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject: userID.String(),
		Issuer: "chirpy",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)	
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	log.Println(ss)
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer != string("chirpy") {
		return uuid.Nil, errors.New("invalid issuer")
	}

	user_uuid, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	
	return user_uuid, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Auth Provided")
	}

	if !strings.Contains(authHeader, "Bearer") {
		return "", errors.New("Auth Headers does not container Bearer")
	}

	keys := strings.Fields(authHeader)

	foundBearer := false
	for _, value := range keys {
		if foundBearer {
			return value, nil	
		}

		if value == "Bearer" {
			foundBearer = true
		}
	}

	// Should never get here
	return "", errors.New("There was an issue retrieving the token")
}
