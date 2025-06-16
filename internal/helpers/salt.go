package helpers

import "crypto/rand"

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)

	_, err := rand.Read(salt)
	if err != nil {
		return salt, err
	}

	return salt, nil
}
