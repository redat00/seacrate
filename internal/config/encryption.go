package config

type EncryptionConfiguration struct {
	EncryptionAlgorithm string `yaml:"algorithm" validate:"oneof=aes"`
}
