package filestore

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/k1nky/gophkeeper/internal/adapter/store"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

type FileStore struct {
	Path string
}

var _ store.ObjectStore = new(FileStore)

func New(path string) *FileStore {
	return &FileStore{
		Path: path,
	}
}

func (fs *FileStore) write(path string, obj *vault.DataReader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, obj)

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
	// TODO: sync.Poolhash
	h := sha256.New()
	s := fmt.Sprintf("%x", h.Sum([]byte(relative)))
	return path.Join(fs.Path, s)
}

func (fs *FileStore) Open(ctx context.Context) error {
	return os.MkdirAll(fs.Path, 0750)
}

func (fs *FileStore) Close() error {
	return nil
}

func (fs *FileStore) Put(ctx context.Context, key string, obj *vault.DataReader) error {
	if obj == nil {
		return nil
	}
	path := fs.path(key)
	err := fs.write(path, obj)
	return err
}

func (fs *FileStore) Get(ctx context.Context, key string) (*vault.DataReader, error) {
	path := fs.path(key)
	return fs.read(path)
}

func (fs *FileStore) Delete(ctx context.Context, key string) error {
	path := fs.path(key)
	return os.Remove(path)
}
