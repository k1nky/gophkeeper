package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"

	"github.com/enceve/crypto/pad"
)

type EncryptReader struct {
	key   []byte
	r     io.Reader
	iv    []byte
	pad   pad.Padding
	block cipher.Block
}

type DecryptReader struct {
	key   []byte
	r     io.Reader
	iv    []byte
	pad   pad.Padding
	block cipher.Block
}

func NewEncryptReader(secret string, r io.Reader) (*EncryptReader, error) {
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return &EncryptReader{
		pad:   pad.NewPKCS7(aes.BlockSize),
		key:   key[:],
		block: block,
		r:     r,
	}, nil
}

func NewDecryptReader(secret string, r io.Reader) (*DecryptReader, error) {
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	return &DecryptReader{
		pad:   pad.NewPKCS7(aes.BlockSize),
		key:   key[:],
		block: block,
		r:     r,
	}, nil
}

func (r *EncryptReader) encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, (len(plaintext)/16+1)*aes.BlockSize)
	padded := r.pad.Pad(plaintext)
	blk := cipher.NewCBCEncrypter(r.block, r.iv)
	blk.CryptBlocks(ciphertext, padded)
	return ciphertext
}

func (r *DecryptReader) decrypt(ciphertext []byte) ([]byte, error) {
	plaintext := make([]byte, (len(ciphertext)/16)*aes.BlockSize)
	blk := cipher.NewCBCDecrypter(r.block, r.iv)
	blk.CryptBlocks(plaintext, ciphertext)
	return r.pad.Unpad(plaintext)
}

func (r *EncryptReader) Read(p []byte) (n int, err error) {
	src := make([]byte, aes.BlockSize)
	n, err = r.r.Read(src)
	if err != nil {
		return
	}
	if r.iv == nil {
		r.iv = make([]byte, aes.BlockSize)
	}
	ciphertext := r.encrypt(src[:n])

	n = copy(p, ciphertext)
	r.iv = ciphertext[len(ciphertext)-aes.BlockSize:]

	return n, err
}

func (r *EncryptReader) Close() error {
	return nil
}

func (r *DecryptReader) Read(p []byte) (n int, err error) {
	src := make([]byte, aes.BlockSize*2)
	n, err = r.r.Read(src)
	if err != nil {
		return
	}
	if r.iv == nil {
		r.iv = make([]byte, aes.BlockSize)
	}
	plaintext, err := r.decrypt(src[:n])
	if err != nil {
		return
	}
	n = copy(p, plaintext)
	r.iv = src[len(src)-aes.BlockSize:]

	return n, err
}

func (r *DecryptReader) Close() error {
	return nil
}
