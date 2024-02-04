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

func (s *Service) GetSecretData(ctx context.Context, metaID vault.MetaID) (*vault.DataReader, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	meta, err := s.store.GetSecretMetaByID(ctx, metaID, uid)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}
	if meta.UserID != claims.ID {
		return nil, user.ErrUnathorized
	}

	return s.store.GetSecretData(ctx, metaID, claims.ID)
}

func (s *Service) GetSecretMeta(ctx context.Context, metaID vault.MetaID) (*vault.Meta, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	return s.store.GetSecretMetaByID(ctx, metaID, uid)
}

func (s *Service) GetSecretMetaByAlias(ctx context.Context, alias string) (*vault.Meta, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	return s.store.GetSecretMetaByAlias(ctx, alias, uid)
}

func (s *Service) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	return s.store.PutSecret(ctx, meta, data)
}

func (s *Service) ListSecretsByUser(ctx context.Context) (vault.List, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	return s.store.ListSecretsByUser(ctx, uid)
}
