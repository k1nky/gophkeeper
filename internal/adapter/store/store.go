// Пакет store предоставляет адаптер для работы с хранилищем секретов.
package store

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
	"golang.org/x/sync/errgroup"
)

// Adapter адаптер хранилища секретов.
type Adapter struct {
	mstore MetaStore
	ostore ObjectStore
}

var _ Store = new(Adapter)

// Close закрывает хранилище секретов. Возвращает ошибку, если закрытие хранилища мета-данных и/или объектов завершилось с ошибкой.
func (a *Adapter) Close() error {
	g := new(errgroup.Group)
	if a.mstore != nil {
		g.Go(a.mstore.Close)
	}
	if a.ostore != nil {
		g.Go(a.ostore.Close)
	}
	return g.Wait()
}

// New возвращает новый адаптер к хранилищу секретов, где mstore определяет способ хранения мета-данных,
// а ostore хранилище самой секретной информации.
func New(mstore MetaStore, ostore ObjectStore) *Adapter {
	return &Adapter{
		mstore: mstore,
		ostore: ostore,
	}
}

// GetSecretData возвращает данные секрета пользователя userID с ид metaID. Обязательно нужно следить за своевременным
// закрытием полученных данных.
func (a *Adapter) GetSecretData(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.DataReader, error) {
	meta, err := a.GetSecretMetaByID(ctx, metaID, userID)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, nil
	}
	return a.ostore.Get(ctx, meta.DataID)
}

// GetSecretMetaByID возвращает мета-данные секрета с ид metaID пользователя userID.
func (a *Adapter) GetSecretMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error) {
	m, err := a.mstore.GetMetaByID(ctx, metaID, userID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return m, err
}

// GetSecretMetaByID возвращает мета-данные секрета с псевдонимом alias пользователя userID.
func (a *Adapter) GetSecretMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error) {
	m, err := a.mstore.GetMetaByAlias(ctx, alias, userID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return m, err
}

// GetUserByLogin возвращает пользователя с именем login.
func (a *Adapter) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	return a.mstore.GetUserByLogin(ctx, login)
}

// ListSecretsByUser возвращает список мета-данных секретов пользователя userID.
func (a *Adapter) ListSecretsByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	return a.mstore.ListMetaByUser(ctx, userID)
}

// NewUser создает нового пользователя u и возвращает указатель на него.
func (a *Adapter) NewUser(ctx context.Context, u user.User) (*user.User, error) {
	return a.mstore.NewUser(ctx, u)
}

// Open открывает хранилище секретов. Возвращает ошибку, если открытие хранилища мета-данных и/или объектов завершилось с ошибкой.
func (a *Adapter) Open(ctx context.Context) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		return a.mstore.Open(ctx)
	})
	g.Go(func() error {
		return a.ostore.Open(ctx)
	})
	err := g.Wait()
	if err != nil {
		a.Close()
	}
	return err
}

// PutSecret добавляет новый секрет в хранилище. Т.к. секрет может быть перенесен из локального хранилища в удаленное
// в meta должен быть указан уникальный ID в пределах одного пользователя, поэтому лучше генерировать его заранее.
// Метод возвращает ссылку на мета-данные добавленного секрета.
func (a *Adapter) PutSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	if len(meta.ID) == 0 {
		meta.ID = vault.NewMetaID()
	}
	meta.DataID = a.buildDataKey(meta)
	// сначала пробуем записать данные секрета в хранилище
	if err := a.ostore.Put(ctx, meta.DataID, data); err != nil {
		return nil, err
	}
	// потом записываем мета-данные
	if _, err := a.mstore.NewMeta(ctx, meta); err != nil {
		// если их записать не удалось, то удаляем секрет
		a.ostore.Delete(ctx, meta.DataID)
		return nil, err
	}
	return &meta, nil
}

func (a *Adapter) UpdateSecretMeta(ctx context.Context, meta vault.Meta) (*vault.Meta, error) {
	return a.mstore.UpdateMeta(ctx, meta)
}

func (a *Adapter) UpdateSecret(ctx context.Context, meta vault.Meta, data *vault.DataReader) (*vault.Meta, error) {
	cm, err := a.GetSecretMetaByID(ctx, meta.ID, meta.UserID)
	if err != nil {
		return nil, err
	}
	if cm == nil {
		return nil, vault.ErrMetaNotExists
	}
	meta.DataID = a.buildDataKey(meta)
	// сначала пробуем записать данные секрета в хранилище
	if err := a.ostore.Put(ctx, meta.DataID, data); err != nil {
		return nil, err
	}
	// потом записываем мета-данные
	if _, err := a.mstore.UpdateMeta(ctx, meta); err != nil {
		// если их записать не удалось, то удаляем секрет
		a.ostore.Delete(ctx, meta.DataID)
		return nil, err
	}
	a.ostore.Delete(ctx, cm.DataID)
	return &meta, nil
}

func (a *Adapter) DeleteSecret(ctx context.Context, meta vault.Meta) error {
	if err := a.mstore.DeleteMeta(ctx, meta); err != nil {
		return err
	}
	if err := a.ostore.Delete(ctx, meta.DataID); err != nil {
		return err
	}
	return nil
}

func (a *Adapter) buildDataKey(m vault.Meta) string {
	return fmt.Sprintf("%d-%s-%d", m.UserID, m.ID, m.Revision)
}
