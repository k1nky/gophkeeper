package vault

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/user"
)

type Object interface {
	Bytes() []byte
	Read(p []byte) (n int, err error)
	ReadFrom(r io.Reader) (n int64, err error)
	Write(p []byte) (n int, err error)
	WriteTo(w io.Writer) (n int64, err error)
}
type Meta struct {
	UserID user.ID
	Extra  string
}

type UniqueKey string

type Secret struct {
	Key  UniqueKey
	Data Object
	Meta Meta
}

type List map[UniqueKey]Meta

func NewUniqueKey() UniqueKey {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return UniqueKey(hex.EncodeToString(b))
}
