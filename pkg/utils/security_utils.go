package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/baywiggins/qIt-backend/internal/config"

	"golang.org/x/crypto/bcrypt"
)

var aesKey []byte

func init () {
	aesKey = []byte(config.AESKey)
}

// Hash password
func HashPassword(password string) (string, error) {
	var err error;

	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min const
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

// Check if passwords match
func DoPasswordsMatch(hashedPassword, curPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(curPassword))

	return err == nil
}

// Encrypt plain text using AES-GCM
func Encrypt(s string) (string, error) {
	var err error;

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}
	// Generate a nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Create GCM mode instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt
	cipherText := aesGCM.Seal(nil, nonce, []byte(s), nil)

	// Combine nonce and ciphertext for easier storage
	final := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(final), err
}

// Decrypt decrypts the base64-encoded cipher text using AES-GCM with the given key
func Decrypt(cipherText string) (string, error) {
	decodedCipherText, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	// Extract nonce (first 12 bytes)
	if len(decodedCipherText) < 12 {
		return "", errors.New("invalid ciphertext: too short")
	}
	nonce, cipherTextBytes := decodedCipherText[:12], decodedCipherText[12:]

	// Create GCM mode instance
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}