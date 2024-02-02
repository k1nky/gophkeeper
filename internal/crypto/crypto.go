// Пакет crypto содержит читателей для работы с AES-CBC.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"

	"github.com/enceve/crypto/pad"
)

// EncryptReader читатель для шифрования блока данных алгоритмом AES в режиме CBC из другого читателя.
// Можно считать как middleware для io.Reader. Реализует интерфейс io.ReadCloser.
// Результат шифрования аналогичен openssl enc -aes-256-cbc -nosalt -e -out <file> -K "<key>" -iv 0.
type EncryptReader struct {
	key   []byte
	r     io.Reader
	iv    []byte
	pad   pad.Padding
	block cipher.Block
}

// DecryptReader читатель для расшифровки блока данных алгоритмом AES в режиме CBC из другого читателя.
// Можно считать как middleware для io.Reader. Реализует интерфейс io.ReadCloser.
// Результат расшифрования аналогичен openssl enc -aes-256-cbc -nosalt -d -out <file> -K "<key>" -iv 0.
type DecryptReader struct {
	key   []byte
	r     io.Reader
	iv    []byte
	pad   pad.Padding
	block cipher.Block
}

// NewEncryptReader возвращет новый EncryptReader с ключом `secret` для исходного читателя открытых данных `r` и вектором инициализации `iv`.
// Если iv равен nil, то будет использоваться значение 0.
func NewEncryptReader(secret string, r io.Reader, iv []byte) (*EncryptReader, error) {
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
		iv:    iv,
	}, nil
}

// NewDecryptReader возвращет новый DecryptReader с ключом `secret` для исходного читателя зашифрованных данных `r` и вектором инициализации `iv`.
// Если iv равен nil, то будет использоваться значение 0.
func NewDecryptReader(secret string, r io.Reader, iv []byte) (*DecryptReader, error) {
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
		iv:    iv,
	}, nil
}

func (r *EncryptReader) encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, (len(plaintext)/16+1)*aes.BlockSize)
	// выравниваем до размера блока
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

// Read читает небольше aes.BlockSize из исходного читателя и шифрует их.
func (r *EncryptReader) Read(p []byte) (n int, err error) {
	src := make([]byte, aes.BlockSize)
	// читаем из источника не больше одного блока
	n, err = r.r.Read(src)
	if err != nil {
		return
	}
	if r.iv == nil {
		r.iv = make([]byte, aes.BlockSize)
	}
	// шифруем
	ciphertext := r.encrypt(src[:n])

	n = copy(p, ciphertext)
	// каждый последующий вектор инициализации должен быть равен последнему блоку зашифрованного текста
	r.iv = ciphertext[len(ciphertext)-aes.BlockSize:]

	return n, err
}

// Close для реализации io.ReadCloser
func (r *EncryptReader) Close() error {
	return nil
}

// Read читает небольше 2*aes.BlockSize из исходного читателя и расшифровывает их.
func (r *DecryptReader) Read(p []byte) (n int, err error) {
	// читать будем больше из-за возможного выравнивания
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

// Close для реализации io.ReadCloser
func (r *DecryptReader) Close() error {
	return nil
}
