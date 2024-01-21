package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	DefaultMaxKeepaliveConnections = 10
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type PostgresStorage struct {
	dsn string
	*sql.DB
}

func New(dsn string) *PostgresStorage {
	return &PostgresStorage{
		dsn: dsn,
	}
}

// Открывает новое подключение к базе
func (ps *PostgresStorage) Open(ctx context.Context) (err error) {
	if ps.DB, err = sql.Open("pgx", ps.dsn); err != nil {
		return
	}
	ps.DB.SetMaxIdleConns(DefaultMaxKeepaliveConnections)
	ps.DB.SetMaxOpenConns(DefaultMaxKeepaliveConnections)
	err = ps.Initialize(ps.dsn)
	if err == nil {
		// закрывает соединение при отмене контекста
		go func() {
			<-ctx.Done()
			ps.Close()
		}()
	}
	return err
}

// Применяет миграции
func (ps *PostgresStorage) Initialize(dsn string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	// отсутствие изменений при миграции не считаем ошибкой
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (ps *PostgresStorage) hasUniqueViolationError(err error) bool {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerrcode.IsIntegrityConstraintViolation(pgerr.Code) {
			return true
		}
	}
	return false
}

// Закрывает подключение к базе
func (ps *PostgresStorage) Close() error {
	return ps.DB.Close()
}
