package auth

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	uuid1 := uuid.New()	
	secret1 := "supersecret123"
	bad_secret := "notcorrect"
	expires1, err := time.ParseDuration("1h")
	token1, _ := MakeJWT(uuid1, secret1, expires1)

	if err != nil {
		log.Println("Failed to parse duration")
		return
	}

	tests := []struct {
		name     string
		secret   string
		token   string
		wantErr  bool
	}{
		{
			name: "Correct secret && token",
			secret: secret1,
			token: token1,
			wantErr: false,
		},
		{ 
			name: "Incorrect secret",
			secret: bad_secret,
			token: token1,
			wantErr: true,
		},
		{
			name: "Invalid token",
			secret: secret1,
			token: "invaid.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
