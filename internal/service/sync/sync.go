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

func (s *Service) PullAll(ctx context.Context) (vault.List, error) {
	pulled := vault.List{}
	list, err := s.client.ListSecrets(ctx)
	if err != nil {
		return pulled, err
	}
	for _, v := range list {
		m, err := s.Pull(ctx, v)
		if err != nil {
			return pulled, err
		}
		pulled = append(pulled, *m)
	}
	return pulled, nil
}

func (s *Service) Push(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
	data, err := s.storage.GetSecretData(ctx, meta.ID)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	m, err := s.client.PutSecret(ctx, meta, data)
	return m, err
}

func (s *Service) PushAll(ctx context.Context) (vault.List, error) {
	pushed := vault.List{}
	list, err := s.storage.ListSecretsByUser(ctx)
	if err != nil {
		return pushed, err
	}
	for _, v := range list {
		m, err := s.Push(ctx, v)
		if err != nil {
			return pushed, err
		}
		pushed = append(pushed, *m)
	}
	return pushed, nil
}
