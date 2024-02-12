// Пакет keeper содержит сервис хранения секретов.
package keeper

import (
	"context"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

// Service сервис хранения секретов. Во всех запросах текущий пользователей берется из контекста (user.NewContextWithClaims).
// Если в контексте пользователь не опеределен, то будет использоваться "локальный" пользователй с ИД 0.
// Если добавлять расширенную модель разграничения доступа, то в можно добавить дополнительный
// аргумент в методы, определяющий для какого пользователя выполняется запрос.
// Остановимся пока на простой моделе, где пользователь вызывающий метод (EffectiveUser) и
// пользователь владелец данных - одно лицо.
type Service struct {
	store storage
	log   logger
}

// New возвращает новый экземпляр сервиса с хранилищем store и логгером log.
func New(store storage, log logger) *Service {
	s := &Service{
		store: store,
		log:   log,
	}
	return s
}

// GetSecretData возвращает данные секрета с ИД metaID для пользователя определенного в контексте или локального пользователя.
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

	return s.store.GetSecretData(ctx, metaID, claims.ID)
}

// GetSecretMeta возвращает мета-данные секрета с ИД metaID для пользователя определенного в контексте или локального пользователя.
func (s *Service) GetSecretMeta(ctx context.Context, metaID vault.MetaID) (*vault.Meta, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	return s.store.GetSecretMetaByID(ctx, metaID, uid)
}

// GetSecretMetaByAlias возвращает мета-данные секрета с псевдонимом alias для пользователя определенного в контексте или локального пользователя.
func (s *Service) GetSecretMetaByAlias(ctx context.Context, alias string) (*vault.Meta, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	return s.store.GetSecretMetaByAlias(ctx, alias, uid)
}

// PutSecret добавляет или обновляет секрет с мета-данными meta и данным в data.
func (s *Service) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	m, err := s.GetSecretMeta(ctx, meta.ID)
	if err != nil {
		return nil, err
	}
	// за правильность ведения версионирования должен отвечать потребитель сервиса Keeper,
	// но если версия по какой-то причине не установлена - определим ее явно.
	if meta.Revision == 0 {
		meta.Revision = vault.NewRevision()
	}
	if m == nil {
		// добавляем новый секрет
		return s.store.PutSecret(ctx, meta, data)
	}

	// обновляем существующий
	return s.store.UpdateSecret(ctx, meta, data)
}

// DeleteSecret удалет секрет с мета-данным meta. Фактически данный вызов ничего не удаляет,
// а просто помечает что секрет должен быть удален.
func (s *Service) DeleteSecret(ctx context.Context, meta vault.Meta) error {
	m, err := s.GetSecretMeta(ctx, meta.ID)
	if err != nil {
		return nil
	}
	if m == nil {
		return nil
	}
	meta.IsDeleted = true
	meta.Revision = vault.NewRevision()
	_, err = s.store.UpdateSecretMeta(ctx, meta)
	return err
}

// ListSecretsByUser возвращает список секретов. В списке перечисляются только мета-данные секретов.
func (s *Service) ListSecretsByUser(ctx context.Context) (vault.List, error) {
	uid := user.LocalUserID
	claims, ok := user.GetEffectiveUser(ctx)
	if ok {
		uid = claims.ID
	}
	return s.store.ListSecretsByUser(ctx, uid)
}
