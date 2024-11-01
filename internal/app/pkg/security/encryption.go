package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptionService provides methods to encrypt and decrypt data using AES-256.
type EncryptionService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// AES256EncryptionService implements the EncryptionService interface.
type AES256EncryptionService struct {
	key []byte
}

// NewAES256EncryptionService creates a new instance of AES256EncryptionService.
// It retrieves the encryption key using the provided KeyManager.
func NewAES256EncryptionService(km KeyManager) (*AES256EncryptionService, error) {
	key, err := km.GetKey()
	if err != nil {
		return nil, err
	}
	return &AES256EncryptionService{key: key}, nil
}

// Encrypt encrypts the plaintext using AES-256-CBC.
func (s *AES256EncryptionService) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	padded := pkcs7Pad(plaintextBytes, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(padded))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the ciphertext using AES-256-CBC.
func (s *AES256EncryptionService) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	if len(data)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	unpadded, err := pkcs7Unpad(data, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}

// pkcs7Pad pads the plaintext to a multiple of the block size.
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytesRepeat(byte(padding), padding)
	return append(data, padtext...)
}

// pkcs7Unpad removes the padding from the plaintext.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}
	padding := int(data[length-1])
	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid padding")
	}
	for _, v := range data[length-padding:] {
		if int(v) != padding {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:length-padding], nil
}

// bytesRepeat returns a new byte slice consisting of count copies of b.
func bytesRepeat(b byte, count int) []byte {
	result := make([]byte, count)
	for i := range result {
		result[i] = b
	}
	return result
}
