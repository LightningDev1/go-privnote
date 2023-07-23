package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
)

// The following code replicates https://github.com/mdp/gibberish-aes, which is used by Privnote.

var (
	// ErrInvalidCipherText is returned when the ciphertext is invalid.
	ErrInvalidCipherText = errors.New("invalid ciphertext")

	// ErrInvalidPadding is returned when the padding is invalid.
	// This is usually caused by an incorrect password.
	ErrInvalidPadding = errors.New("invalid padding")
)

// pad adds padding to the data according to the PKCS5 standard.
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padtext...)
}

// unpad removes padding from the text according to the PKCS5 standard.
func unpad(data []byte) ([]byte, error) {
	padding := data[len(data)-1]

	if padding > aes.BlockSize {
		return nil, ErrInvalidPadding
	}

	dataLen := len(data) - int(padding)

	if dataLen < 0 || dataLen > len(data) {
		return nil, ErrInvalidPadding
	}

	return data[:dataLen], nil
}

// md5hash returns the MD5 hash of the given byte slice.
func md5hash(s []byte) []byte {
	hash := md5.Sum(s)
	return hash[:]
}

// openSSLKey returns the key and IV for the given password and salt.
func openSSLKey(password, salt []byte) ([]byte, []byte) {
	passSalt := append(password, salt...)

	result := md5hash(passSalt)
	curHash := result

	for i := 0; i < 2; i++ {
		curHash = md5hash(append(curHash, passSalt...))
		result = append(result, curHash...)
	}

	return result[0 : 4*8], result[4*8 : 4*8+16]
}

// Encrypt encrypts the given data with the given password.
// A random password is generated if none is given.
func Encrypt(data, password string) (string, string, error) {
	plainText := []byte(data)

	// Generate a random password if none is given.
	if password == "" {
		password = randomPassword(9)
	}

	// Generate a random salt.
	salt := randomSalt(8)
	saltBlock := append([]byte("Salted__"), salt...)

	plainText = pad(plainText, aes.BlockSize)

	key, iv := openSSLKey([]byte(password), salt)

	// Encrypt the plaintext using AES-256 in CBC mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)

	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)

	// ciphertext = "Salted__" + salt + cipherText
	cipherText = append(saltBlock, cipherText...)

	// Encode the ciphertext using Base64.
	return base64.StdEncoding.EncodeToString(cipherText), password, nil
}

// Decrypt decrypts the given ciphertext with the given password.
func Decrypt(encodedData, password string) (string, error) {
	// Decode the ciphertext using Base64.
	cipherText, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		// The ciphertext is not valid Base64.
		return "", errors.Join(err, ErrInvalidCipherText)
	}

	// Check if the ciphertext is valid.
	if len(cipherText) < 2*aes.BlockSize {
		return "", ErrInvalidCipherText
	}

	if len(cipherText)%aes.BlockSize != 0 {
		return "", ErrInvalidPadding
	}

	// Extract the salt from the ciphertext.
	salt := cipherText[8:16]
	cipherText = cipherText[16:]

	key, iv := openSSLKey([]byte(password), salt)

	// Decrypt the ciphertext using AES-256 in CBC mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)

	plainText, err = unpad(plainText)
	return string(plainText), err
}
