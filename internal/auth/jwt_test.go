package auth

import (
	"log"
	"testing"
	"time"
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	goodHeaders := http.Header {
		"Authorization": {"Bearer TestToken"},
	}

    noAuthHeaders := http.Header {
		
	}

	tests := []struct {
		name    string
		headers http.Header
		wantErr bool
	} {
		{
			name:    "Valid token",
			headers: goodHeaders,
			wantErr: false,
		},
		{
			name:    "No Auth in header",
			headers: noAuthHeaders,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wanterr %v", err, tt.wantErr)
			}
		})		
	}
}
