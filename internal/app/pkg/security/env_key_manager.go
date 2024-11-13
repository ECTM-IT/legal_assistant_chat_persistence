package security

import (
	"errors"
	"os"
)

// EnvironmentKeyManager retrieves the AES-256 key from environment variables.
type EnvironmentKeyManager struct {
	envVar string
}

// NewEnvironmentKeyManager creates a new instance of EnvironmentKeyManager.
// envVar is the name of the environment variable where the key is stored.
func NewEnvironmentKeyManager(envVar string) *EnvironmentKeyManager {
	return &EnvironmentKeyManager{envVar: envVar}
}

// GetKey retrieves the AES-256 key from the specified environment variable.
func (e *EnvironmentKeyManager) GetKey() ([]byte, error) {
	key := os.Getenv(e.envVar)
	if key == "" {
		return nil, ErrKeyNotFound
	}
	if len(key) != 32 {
		return nil, errors.New("AES-256 key must be 32 bytes long")
	}
	return []byte(key), nil
}
