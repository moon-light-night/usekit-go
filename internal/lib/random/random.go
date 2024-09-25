package random

import (
	"log"
	"math/rand"
)

func NewRandomString(wordLength int) string {
	if wordLength <= 0 {
		log.Fatal("length must be greater than zero")
		return ""
	}
	randomWord := ""
	alphabet := []string{"a", "b", "c", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 1; i <= wordLength; i++ {
		randomWord += alphabet[rand.Intn(len(alphabet))]
	}

	return randomWord
}
