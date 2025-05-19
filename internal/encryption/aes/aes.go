package aes

import (
	cryptoAes "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"sync"
)

type AesEncryptionEngine struct {
	key  []byte
	gcm  cipher.AEAD
	once sync.Once
	seal bool
}

func (aes *AesEncryptionEngine) initGcm() {
	var block cipher.Block

	// Init cipher block
	block, err := cryptoAes.NewCipher(aes.key)
	if err != nil {
		panic(err)
	}

	// Init GCM
	aes.gcm, err = cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
}

func (aes *AesEncryptionEngine) EncryptData(data []byte) ([]byte, error) {
	aes.once.Do(aes.initGcm)
	nonce := make([]byte, aes.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error during nonce generation : %v", err)
	}
	encrypted := aes.gcm.Seal(nonce, nonce, data, nil)
	return encrypted, nil
}

func (aes *AesEncryptionEngine) DecryptData(data []byte) ([]byte, error) {
	aes.once.Do(aes.initGcm)
	decrypted, err := aes.gcm.Open(nil, data[:aes.gcm.NonceSize()], data[aes.gcm.NonceSize():], nil)
	if err != nil {
		return nil, fmt.Errorf("error during data decryption : %v", err)
	}
	return decrypted, nil
}

func (aes *AesEncryptionEngine) SetKey(key []byte) {
	aes.key = key
	aes.SetSealStatus(false)
}

func (aes *AesEncryptionEngine) GetSealStatus() bool {
	return aes.seal
}

func (aes *AesEncryptionEngine) SetSealStatus(status bool) {
	aes.seal = status
}

func (aes *AesEncryptionEngine) GenerateKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return key, err
	}
	return key, nil
}

func NewAesEncryptionEngine() *AesEncryptionEngine {
	var engine AesEncryptionEngine
	return &engine
}
