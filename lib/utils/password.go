package utils

import (
	"crypto/sha256"

	"code.google.com/p/go.crypto/pbkdf2"
)

// http://stackoverflow.com/a/19828153/3983047

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

// HashPassword hashes a password with pbkdf2. This is probably good enough.
func HashPassword(password, salt []byte) []byte {
	defer clear(password)
	return pbkdf2.Key(password, salt, 4096, sha256.Size, sha256.New)
}
