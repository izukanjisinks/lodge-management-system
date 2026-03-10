package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateSessionToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Failed to generate token: ", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
