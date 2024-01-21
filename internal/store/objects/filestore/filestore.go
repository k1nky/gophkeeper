package filestore

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

func (fs *FileStore) write(path string, obj vault.Object) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	_, err = obj.WriteTo(f)
	return err
}

func (fs *FileStore) read(path string, obj vault.Object) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0660)
	if err != nil {
		if os.IsNotExist(err) {
			return vault.ErrObjectNotExists
		}
		return err
	}
	defer f.Close()
	_, err = obj.ReadFrom(f)
	return err
}

func (fs *FileStore) path(relative string) string {
	// TODO: sync.Poolhash
	h := md5.New()
	s := hex.EncodeToString(h.Sum([]byte(relative)))
	return path.Join(fs.Path, s)
}

func (fs *FileStore) Open(ctx context.Context) error {
	return os.MkdirAll(fs.Path, 0750)
}

func (fs *FileStore) Close() error {
	return nil
}

func (fs *FileStore) Put(ctx context.Context, key string, obj vault.Object) error {
	path := fs.path(key)
	err := fs.write(path, obj)
	return err
}

func (fs *FileStore) Get(ctx context.Context, key string, obj vault.Object) error {
	path := fs.path(key)
	return fs.read(path, obj)
}

func (fs *FileStore) Delete(ctx context.Context, key string) error {
	path := fs.path(key)
	return os.Remove(path)
}
