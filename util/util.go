package util

import (
	"math/rand"
	"time"
)

var (
	randomPool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	poolLength = len(randomPool)
)

// RandomString generates random string of given length
// It sets rand see on each call and returns generated string.
func RandomString(length int) string {
	str := make([]byte, length)
	rand.Seed(time.Now().UTC().UnixNano())

	for i := range str {
		str[i] = randomPool[rand.Intn(poolLength)]
	}

	return string(str)
}
