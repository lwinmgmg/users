package utils

import "crypto/sha256"

func Hash256(input string) []byte {
	sha := sha256.New()
	sha.Write([]byte(input))
	return sha.Sum(nil)
}
