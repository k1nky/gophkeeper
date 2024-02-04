// Пакет filestore предоставляет простое хранилище бинарных данных в файлах.
package filestore

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"

	"github.com/k1nky/gophkeeper/internal/adapter/store"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

// FileStore хранилище бинарных данных в файлах.
type FileStore struct {
	// Path путь до каталога, в котором будут располагаться файлы хранилища
	Path string
}

var _ store.ObjectStore = new(FileStore)

// New возвращает экземпляр хранилища.
func New(path string) *FileStore {
	return &FileStore{
		Path: path,
	}
}

// Close закрывает хранилище.
func (fs *FileStore) Close() error {
	return nil
}

// Delete удаляет объект с ключом key из хранилизща.
func (fs *FileStore) Delete(ctx context.Context, key string) error {
	path := fs.path(key)
	return os.Remove(path)
}

// Get возвращает данные по указанному ключу key. Во избежании утечки открытых файлов DataReader следует закрывать Close после прочтения.
func (fs *FileStore) Get(ctx context.Context, key string) (*vault.DataReader, error) {
	path := fs.path(key)
	return fs.read(path)
}

// Open открывает хранилище.
func (fs *FileStore) Open(ctx context.Context) error {
	return os.MkdirAll(fs.Path, 0750)
}

// Put кладет новые данные data c ключом key.
func (fs *FileStore) Put(ctx context.Context, key string, data *vault.DataReader) error {
	if data == nil {
		return nil
	}
	path := fs.path(key)
	err := fs.write(path, data)
	return err
}

func (fs *FileStore) write(path string, data *vault.DataReader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)

	return err
}

func (fs *FileStore) read(path string) (*vault.DataReader, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0660)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, vault.ErrObjectNotExists
		}
		return nil, err
	}

	return vault.NewDataReader(f), nil
}

func (fs *FileStore) path(relative string) string {
	// TODO: sync.Pool hash
	h := sha256.New()
	s := hex.EncodeToString(h.Sum([]byte(relative)))
	return path.Join(fs.Path, s)
}
