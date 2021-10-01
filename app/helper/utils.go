package helper

import (
	"math/rand"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GetRandomString generating random string with length 8
func GetRandomString() string {
	b := make([]byte, 8)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}

	return string(b)
}
