package sync

import (
	"context"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

type Service struct {
	client  client
	storage storage
}

func New() *Service {
	return &Service{}
}

func (s *Service) Run(ctx context.Context) error {
	list, err := s.client.ListSecrets(ctx)
	if err != nil {
		return err
	}
	for _, v := range list {
		meta, err := s.client.GetSecretMeta(ctx, v.Key)
		if err != nil {
			return err
		}
		r, w := io.Pipe()
		data := vault.NewDataReader(r)
		s.client.GetSecretData(ctx, v.Key, w)
		s.storage.PutSecret(ctx, *meta, data)
	}
	return nil
}
