package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts the given data using AES encryption with the provided key.
// It returns the encrypted data as a base64 encoded string.
func Encrypt(data string, key []byte) (string, error) {
	// Create a new AES cipher block from the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Convert the data to bytes
	plaintext := []byte(data)

	// Create a slice to hold the ciphertext (AES block size + length of the plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Generate a random initialization vector (IV) for AES encryption
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Create a new AES CFB encrypter stream with the block cipher and the IV
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt the plaintext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Return the base64 encoded encrypted string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given base64 encoded data using AES decryption with the provided key.
// It returns the decrypted data as a string.
func Decrypt(data string, key []byte) (string, error) {
	// Decode the base64 encoded data
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block from the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Check if the ciphertext is long enough to contain the AES block size
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract the IV from the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Create a new AES CFB decrypter stream with the block cipher and the IV
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt the ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	// Return the decrypted data as a string
	return string(ciphertext), nil
}
