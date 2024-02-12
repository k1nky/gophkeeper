// Пакет sync предоставляет инструменты для синхронизации секретов с удаленным сервером.
package sync

import (
	"context"
	"errors"
	"io"

	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"golang.org/x/sync/errgroup"
)

// Служба синхронизации секретов.
type Service struct {
	client  client
	storage storage
	log     logger
}

// New возвращает новый экземпляр сервиса.
func New(client client, storage storage, log logger) *Service {
	return &Service{
		client:  client,
		storage: storage,
		log:     log,
	}
}

// Pull забирает секрет с мета-данным meta из удаленного хранилища в локальное.
func (s *Service) Pull(ctx context.Context, meta vault.Meta, force bool) (*vault.Meta, error) {
	var newMeta *vault.Meta

	m, err := s.storage.GetSecretMeta(ctx, meta.ID)
	if err != nil {
		return nil, err
	}
	if m != nil {
		if m.Equal(meta) {
			// данные секрета уже актуальны
			return nil, vault.ErrNothingToUpdate
		}
		if !m.CanUpdated(meta) && !force {
			// конфликт версий
			return nil, vault.ErrConflictVersion
		}
	}

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
	if err := g.Wait(); err != nil {
		newMeta = nil
	}
	return newMeta, nil
}

// PullAll забирает все секреты пользователя из удаленного хранилища в локальное.
func (s *Service) PullAll(ctx context.Context, force bool) error {
	list, err := s.client.ListSecrets(ctx)
	if err != nil {
		return err
	}
	for _, v := range list {
		if _, err := s.Pull(ctx, v, force); err != nil {
			if errors.Is(err, vault.ErrNothingToUpdate) {
				s.log.Debugf("%s %s", err, v)
			}
			if errors.Is(err, vault.ErrConflictVersion) {
				s.log.Errorf("%s %s", err, v)
			}
			continue
		}
	}
	return nil
}

// Push отправляет секрет из локального хранилища в удаленное.
func (s *Service) Push(ctx context.Context, meta vault.Meta, force bool) (*vault.Meta, error) {

	m, _ := s.client.GetSecretMeta(ctx, meta.ID)
	if m != nil {
		if meta.Equal(*m) {
			return nil, vault.ErrNothingToUpdate
		}
		if !m.CanUpdated(meta) && !force {
			return nil, vault.ErrConflictVersion
		}
	}

	data, err := s.storage.GetSecretData(ctx, meta.ID)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return s.client.PutSecret(ctx, meta, data)
}

// PushAll отправляет все секреты пользователя из локального хранилища в удаленное.
func (s *Service) PushAll(ctx context.Context, force bool) error {
	list, err := s.storage.ListSecretsByUser(ctx)
	if err != nil {
		s.log.Errorf("push: %s", err)
		return err
	}
	for _, v := range list {
		if _, err := s.Push(ctx, v, force); err != nil {
			if errors.Is(err, vault.ErrNothingToUpdate) {
				s.log.Debugf("%s %s", err, v)
			}
			if errors.Is(err, vault.ErrConflictVersion) {
				s.log.Errorf("%s %s", err, v)
			}
			continue
		}
	}
	return nil
}
