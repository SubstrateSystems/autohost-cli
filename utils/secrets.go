package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		fmt.Println(fmt.Errorf("no se pudo generar APP_KEY: %w", err))
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)[:length]
}
