package postgres

import (
	"context"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
	"github.com/k1nky/gophkeeper/internal/entity/vault"
)

// user_id INT,
// meta_unique_key VARCHAR(100) UNIQUE NOT NULL,
// alias VARCHAR(100),
// type secret_type,
// extra text,

func (ps *PostgresStorage) NewMeta(ctx context.Context, m vault.Meta) (*vault.Meta, error) {

	const query = `
		INSERT INTO meta AS m (user_id, meta_unique_key, alias, type, extra)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING m.meta_id
	`

	row := ps.QueryRowContext(ctx, query, m.UserID, m.ID, m.Alias, m.Type, m.Extra)
	if err := row.Err(); err != nil {
		if ps.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("%s %w", m.ID, user.ErrDuplicateLogin)
		}
		return nil, NewExecutingQueryError(err)
	}
	// if err := row.Scan(&u.ID); err != nil {
	// 	return nil, NewExecutingQueryError(err)
	// }
	return &m, nil

}

func (ps *PostgresStorage) GetMetaByID(ctx context.Context, metaID vault.MetaID, userID user.ID) (*vault.Meta, error) {
	return nil, nil
}

func (ps *PostgresStorage) GetMetaByAlias(ctx context.Context, alias string, userID user.ID) (*vault.Meta, error) {
	if len(alias) == 0 {
		return nil, nil
	}
	return nil, nil
}

func (ps *PostgresStorage) ListMetaByUser(ctx context.Context, userID user.ID) (vault.List, error) {
	list := vault.List{}
	return list, nil
}
