package vault

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"

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

type MetaID string

// Meta мета-данные секрета.
type Meta struct {
	// Псевдоним
	Alias string
	// Идентификатор данных секрета
	DataID string
	// Поле для дополнительных данных
	Extra string
	// ИД секрета
	ID MetaID
	// Метка удаления, если true, то секрет можно считать удаленным
	IsDeleted bool
	// Тип секрета
	Type SecretType
	// Временная метка версии секрета. Временная зона должна быть UTC.
	UpdatedAt int64
	// ИД пользователя владельца секрета
	UserID user.ID
}

// Список мета-данных секретов.
type List []Meta

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

// NewMetaID возвращает новый уникальный ИД секрета.
func NewMetaID() MetaID {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return MetaID(hex.EncodeToString(b))
}

func (m Meta) String() string {
	return fmt.Sprintf("%s %s %s %s", m.ID, m.Alias, m.Type, m.UpdateAtLocalTime())
}

// CanUpdated возвращает true если секрет может быть обновлен секретом `update`.
func (m Meta) CanUpdated(update Meta) bool {
	return m.ID == update.ID && update.UpdatedAt > m.UpdatedAt
}

// UpdateAtLocalTime возвращает временную метку секрета как тип Time.
func (m Meta) UpdateAtLocalTime() time.Time {
	return time.Unix(m.UpdatedAt, 0)
}

func (l List) String() string {
	s := strings.Builder{}
	for _, v := range l {
		s.WriteString(v.String() + "\n")
	}
	return s.String()
}

// Возвращает текущее время в unix формате в зоне UTC.
func Now() int64 {
	return time.Now().UTC().Unix()
}
