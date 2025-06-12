package auth

import (
	"github.com/google/uuid"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"errors"
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
