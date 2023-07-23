package util

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
)

var (
	charList    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	charListLen = len(charList)
)

// randomInt returns a random integer between 0 and max.
// If generating a secure random integer fails, a pseudo-random integer is generated instead.
func randomInt(max int) int {
	n, err := crand.Int(crand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return rand.Intn(max)
	}

	return int(n.Int64())
}

// randomPassword returns a random password of the given length.
func randomPassword(length int) (password string) {
	for i := 0; i < length; i++ {
		randLetter := charList[randomInt(charListLen)]

		password += string(randLetter)
	}

	return
}

// randomSalt returns a random salt of the given length.
func randomSalt(length int) (salt []byte) {
	for i := 0; i < length; i++ {
		salt = append(salt, byte(randomInt(256)))
	}

	return
}
