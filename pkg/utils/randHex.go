package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func RandomHex() string {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Println("can`t make a hex string")
	}
	return hex.EncodeToString(bytes)
}
