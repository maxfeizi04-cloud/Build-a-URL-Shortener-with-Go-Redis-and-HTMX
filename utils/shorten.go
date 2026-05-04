package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GetShortCode() string {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("Failed to generate random bytes:", err)
	}
	code := base64.URLEncoding.EncodeToString(b)
	log.Println("Generated short code:", code)
	return code
}
