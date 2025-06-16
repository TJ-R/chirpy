package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	byte_arr := make([]byte, 32)
	rand.Read(byte_arr)

	encoded_key := hex.EncodeToString(byte_arr)

	return encoded_key, nil
} 
