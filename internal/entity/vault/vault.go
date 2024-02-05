package vault

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

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
	// буферизированный читатель, позволяет читать исходные данные через буфер
	reader *bufio.Reader
}

type SecretType int

const (
	TypeText SecretType = iota
	TypeLoginPassword
	TypeCreditCard
	TypeFile
)

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

type MetaID string

type Meta struct {
	UserID user.ID
	ID     MetaID
	Alias  string
	Type   SecretType
	Extra  string
}

type List []Meta

func NewMetaID() MetaID {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return MetaID(hex.EncodeToString(b))
}

func (t SecretType) String() string {
	switch t {
	case TypeText:
		return "TEXT"
	case TypeLoginPassword:
		return "LOGIN_PASSWORD"
	case TypeCreditCard:
		return "CREDIT_CARD"
	case TypeFile:
		return "FILE"
	}
	return "UNKNOWN"
}

func (m Meta) String() string {
	return fmt.Sprintf("%s %s %s", m.ID, m.Alias, m.Type)
}

func (l List) String() string {
	s := strings.Builder{}
	for _, v := range l {
		s.WriteString(v.String() + "\n")
	}
	return s.String()
}
