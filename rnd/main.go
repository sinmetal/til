package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func main() {
	fmt.Println(base64.StdEncoding.EncodeToString(GenerateRandomKey(32)))
}

func GenerateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
