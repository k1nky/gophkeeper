package vault

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/user"
)

type BytesBuffer struct {
	bytes.Buffer
}

func (bb *BytesBuffer) Close() error {
	return nil
}

func NewBytesBuffer(p []byte) *BytesBuffer {
	return &BytesBuffer{
		Buffer: *bytes.NewBuffer(p),
	}
}

type DataReader struct {
	origin io.ReadCloser
	reader *bufio.Reader
}

func NewDataReader(r io.ReadCloser) *DataReader {
	return &DataReader{
		origin: r,
		reader: bufio.NewReader(r),
	}
}

func (d *DataReader) Read(p []byte) (n int, err error) {
	return d.reader.Read(p)
}

func (d *DataReader) WriteTo(w io.Writer) (n int64, err error) {
	return d.reader.WriteTo(w)
}

func (d *DataReader) Close() error {
	return d.origin.Close()
}

type UniqueKey string

type Meta struct {
	UserID user.ID
	Key    UniqueKey
	Extra  string
}

type List map[UniqueKey]Meta

func NewUniqueKey() UniqueKey {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return UniqueKey(hex.EncodeToString(b))
}
