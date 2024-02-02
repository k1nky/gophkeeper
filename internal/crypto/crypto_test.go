package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateRandom(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func TestEncrypt(t *testing.T) {
	cases := []int{1, 10, 100, 300, 512, 1000, 10000, 100000}
	for _, size := range cases {
		original := generateRandom(size)
		key := generateRandom(32)

		enc, err := NewEncryptReader(hex.EncodeToString(key), bytes.NewBuffer(original), nil)
		assert.NoError(t, err)
		cipher := bytes.NewBuffer(nil)
		n, err := cipher.ReadFrom(enc)
		assert.Greater(t, n, int64(0))
		assert.NoError(t, err)
		assert.NotEqual(t, original, cipher.Bytes())

		dec, err := NewDecryptReader(hex.EncodeToString(key), cipher, nil)
		assert.NoError(t, err)
		plain := bytes.NewBuffer(nil)
		n, err = plain.ReadFrom(dec)
		assert.Greater(t, n, int64(0))
		assert.NoError(t, err)
		assert.Equal(t, original, plain.Bytes())
	}
}
