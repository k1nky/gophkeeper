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
	// Версия секрета. Алгоритм повышения версии должен работать с учетом того, что клиенты могут быть на разных хостах.
	Revision int64
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
	return fmt.Sprintf("%s %s %s %d", m.ID, m.Alias, m.Type, m.Revision)
}

// CanUpdated возвращает true если секрет может быть обновлен секретом update.
func (m Meta) CanUpdated(update Meta) bool {
	return m.ID == update.ID && update.Revision > m.Revision
}

// Equal возвращает true если идентификаторы и версии m и target равны.
func (m Meta) Equal(target Meta) bool {
	return m.ID == target.ID && m.Revision == target.Revision
}

// NewRevision возвращает номер для новой версии секрета. Возможен конфликт, если несколько клиентов обновят секрет
// с одним ИД с точностью до секунды. В рамках учебного проекта, данным фактом считаю можно пренебречь.
func NewRevision() int64 {
	// используем UTC чтобы не привязываться к верменной зоне клиента
	return time.Now().UTC().Unix()
}

func (l List) String() string {
	s := strings.Builder{}
	for _, v := range l {
		s.WriteString(v.String() + "\n")
	}
	return s.String()
}
