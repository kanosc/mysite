package mycrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

func GenerateRand(byteSize int) ([]byte, error) {
	result := make([]byte, byteSize)
	_, err := rand.Read(result)
	return result, err
}

func AES256Encrypt(plainText, nonce, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid aes key size")
	}
	if len(nonce) != 12 {
		return nil, errors.New("invalid aes nonce size")
	}
	aes, err := aes.NewCipher(key)
	aesgcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}
	return aesgcm.Seal(nil, nonce, plainText, nil), nil
}

func AES256Decrypt(secretText, nonce, key []byte) ([]byte, error) {
	aes, err := aes.NewCipher(key)
	aesgcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}
	return aesgcm.Open(nil, nonce, secretText, nil)
}

func ScriptGenerateKey(password, nonce []byte) ([]byte, error) {
	return scrypt.Key(password, nonce, 1<<15, 8, 1, 32)
}

func EncodeBytesToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func DecodeBase64ToBytes(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
