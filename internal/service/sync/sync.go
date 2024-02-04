package sync

import (
	"context"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	client  client
	storage storage
}

func New(client client, storage storage) *Service {
	return &Service{
		client:  client,
		storage: storage,
	}
}

func (s *Service) Pull(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
	var newMeta *vault.Meta

	g := new(errgroup.Group)
	r, w := io.Pipe()
	data := vault.NewDataReader(r)
	g.Go(func() error {
		err := s.client.GetSecretData(ctx, meta.ID, w)
		w.Close()
		return err
	})
	g.Go(func() error {
		m, err := s.storage.PutSecret(ctx, meta, data)
		newMeta = m
		return err
	})
	err := g.Wait()
	if err != nil {
		newMeta = nil
	}
	return newMeta, nil
}

func (s *Service) Push(ctx context.Context, meta vault.Meta) error {
	data, err := s.storage.GetSecretData(ctx, meta.ID)
	if err != nil {
		return err
	}
	defer data.Close()
	_, err = s.client.PutSecret(ctx, meta, data)
	return err
}
