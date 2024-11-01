package security

import (
	"errors"
)

// KeyManager defines the methods for managing encryption keys.
type KeyManager interface {
	GetKey() ([]byte, error)
}

// ErrKeyNotFound is returned when the encryption key is not found.
var ErrKeyNotFound = errors.New("encryption key not found")
