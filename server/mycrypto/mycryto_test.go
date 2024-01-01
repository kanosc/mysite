package mycrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAndDecrypt(t *testing.T) {
	nonce, _ := GenerateRand(12)
	plainInfo := "Hello world"
	t.Logf("before encrypt[%s]", plainInfo)
	plainText := []byte(plainInfo)
	key, _ := GenerateRand(32)
	secretText, _ := AES256Encrypt(plainText, nonce, key)
	t.Logf("after encrypt[%s]", string(secretText))

	decreptResult, _ := AES256Decrypt(secretText, nonce, key)
	t.Logf("after decrypt[%s]", string(decreptResult))
	assert.Equal(t, plainInfo, string(decreptResult))
}
