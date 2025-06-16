package errors

import (
	"fmt"
)

type ErrSecretDuplicateKey struct {
	Key string
}

func (e ErrSecretDuplicateKey) Error() string {
	return fmt.Sprintf("a secret with the same key (%s) already exist in this folder", e.Key)
}

type ErrOverridingFolder struct {
	Key string
}

func (e ErrOverridingFolder) Error() string {
	return fmt.Sprintf("creating a secret with this key (%s) would override a folder", e.Key)
}

type ErrOverridingSecret struct {
	Key string
}

func (e ErrOverridingSecret) Error() string {
	return fmt.Sprintf("creating a secret with this key (%s) would override another secret", e.Key)
}

type ErrSecretNotFound struct {
	Key string
}

func (e ErrSecretNotFound) Error() string {
	return fmt.Sprintf("no secret found at this key (%s)", e.Key)
}
