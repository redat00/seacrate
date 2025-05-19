package encryption

import (
	"github.com/redat00/seacrate/internal/config"
	"github.com/redat00/seacrate/internal/encryption/aes"
)

type EncryptionEngine interface {
	EncryptData([]byte) ([]byte, error)
	DecryptData([]byte) ([]byte, error)
	GenerateKey(int) ([]byte, error)
	SetKey([]byte)
	GetSealStatus() bool
	SetSealStatus(bool)
}

func NewEncryptionEngine(config config.EncryptionConfiguration) (EncryptionEngine, error) {
	var encryptionEngine EncryptionEngine
	var err error

	switch algorithm := config.EncryptionAlgorithm; algorithm {
	case "aes":
		encryptionEngine = aes.NewAesEncryptionEngine()
		encryptionEngine.SetSealStatus(true)
	}

	return encryptionEngine, err
}
