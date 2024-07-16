package generate

import (
	"math/rand"
	"time"
)

const dictionary = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandString is random string generation
func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = dictionary[seededRand.Intn(len(dictionary))]
	}
	return string(b)
}
