package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/k1nky/gophkeeper/internal/entity/user"
)

// Возвращает пользователя с указанным логином. Nil - пользователь не найден
func (ps *PostgresStorage) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	u := &user.User{
		Login: login,
	}

	const query = `SELECT user_id, password FROM users WHERE login=$1`
	row := ps.QueryRowContext(ctx, query, login)
	if err := row.Err(); err != nil {
		return nil, NewExecutingQueryError(err)
	}
	if err := row.Scan(&u.ID, &u.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, NewExecutingQueryError(err)
	}
	return u, nil
}

// Добавляет и возвращает нового пользователя
func (ps *PostgresStorage) NewUser(ctx context.Context, u user.User) (*user.User, error) {

	const query = `
		INSERT INTO users AS u (login, password)
		VALUES ($1, $2)
		RETURNING u.user_id
	`

	row := ps.QueryRowContext(ctx, query, u.Login, u.Password)
	if err := row.Err(); err != nil {
		if ps.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("%s %w", u.Login, user.ErrDuplicateLogin)
		}
		return nil, NewExecutingQueryError(err)
	}
	if err := row.Scan(&u.ID); err != nil {
		return nil, NewExecutingQueryError(err)
	}
	return &u, nil
}
