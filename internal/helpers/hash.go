package helpers

import (
	"bytes"

	"golang.org/x/crypto/argon2"
)

const (
	saltLen   = 32
	time      = 3
	memory    = 12288
	threads   = 1
	keyLength = 32
)

// Generate Hash from password and return password and salt.
// Salt is optional and use only for hash generation.
// It should be specified for hash comparison.
func GenerateHash(password []byte, salt []byte) ([]byte, []byte, error) {
	var err error

	if len(salt) == 0 {
		salt, err = GenerateSalt(saltLen)
	}
	if err != nil {
		return nil, nil, err
	}

	hash := argon2.IDKey(password, salt, time, memory, threads, keyLength)
	return hash, salt, nil
}

// Compare password and hash. If hash does not equal, return false, else true.
func Compare(password []byte, hash []byte, salt []byte) (bool, error) {
	pHash, _, err := GenerateHash(password, salt)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(pHash, hash) {
		return false, nil
	}

	return true, nil
}
