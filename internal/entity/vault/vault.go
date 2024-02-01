package vault

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/user"
)

// BytesBuffer представляет собой bytes.Buffer дополнительно реализующий io.Closer.
// Удобно использовать для создания DataReader.
type BytesBuffer struct {
	bytes.Buffer
}

// DataReader читатель данных из хранилища секретов. Позволяет читать большой объем данных.
// После использоваения обязательно следует его закрывать. Ответственность за закрытие
// лежит на потребителе данных, т.е. того, кто их запросил.
type DataReader struct {
	// читатель исходных данных, сохраняем его, чтобы была возможность потом закрыть
	origin io.ReadCloser
	// буферизированный читатель
	reader *bufio.Reader
}

func NewBytesBuffer(p []byte) *BytesBuffer {
	return &BytesBuffer{
		Buffer: *bytes.NewBuffer(p),
	}
}

func NewDataReader(r io.ReadCloser) *DataReader {
	return &DataReader{
		origin: r,
		reader: bufio.NewReader(r),
	}
}

func (bb *BytesBuffer) Close() error {
	return nil
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

type List []Meta

func NewUniqueKey() UniqueKey {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return UniqueKey(hex.EncodeToString(b))
}
