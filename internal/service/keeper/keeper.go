package keeper

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

type Service struct {
	store storage
	log   logger
}

func New(store storage, log logger) *Service {
	s := &Service{
		store: store,
		log:   log,
	}
	return s
}

func (s *Service) GetSecretData(ctx context.Context, uk vault.UniqueKey) (*vault.DataReader, error) {
	claims, ok := user.GetEffectiveUser(ctx)
	if !ok {
		// TODO:
	}
	meta, err := s.store.GetSecretMeta(ctx, uk)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}
	if meta.UserID != claims.ID {
		return nil, user.ErrUnathorized
	}

	return s.store.GetSecretData(ctx, uk)
}

// TODO: verify userID
func (s *Service) GetSecretMeta(ctx context.Context, uk vault.UniqueKey) (*vault.Meta, error) {
	return s.store.GetSecretMeta(ctx, uk)
}

func (s *Service) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	return s.store.PutSecret(ctx, meta, data)
}

func (s *Service) ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	return s.store.ListSecretsByUser(ctx, userID)
}
