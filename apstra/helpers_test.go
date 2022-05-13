package apstra

import (
	"math/rand"
	"time"
)

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
