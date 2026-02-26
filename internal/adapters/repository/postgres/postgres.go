package postgres

import (
	"context"
	"fmt"
	"shortener/config"
	"shortener/internal/domain"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

const errCodeAlreadyExist = "23505"

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, config config.Postgres) (*PostgresRepository, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = int32(config.MaxConns)
	poolCfg.MinConns = int32(config.MinConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		pool: pool,
	}, nil
}

func (r *PostgresRepository) Save(ctx context.Context, original, shortened string) error {
	query := `
	insert into urls(original, shortened)
	values ($1, $2)
`
	_, err := r.pool.Exec(ctx, query, original, shortened)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errCodeAlreadyExist {
				return domain.ErrAlreadyExist
			}
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) GetByShortened(ctx context.Context, shortened string) (string, error) {
	query := `select original from urls where shortened = $1`

	original := ""
	err := r.pool.QueryRow(ctx, query, shortened).Scan(&original)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return "", domain.ErrNotFound
		}

		return "", err
	}

	return original, nil
}

func (r *PostgresRepository) GetByOriginal(ctx context.Context, origin string) (string, error) {
	query := `select shortened from urls where original = $1`

	shortened := ""
	err := r.pool.QueryRow(ctx, query, origin).Scan(&shortened)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return "", domain.ErrNotFound
		}

		return "", err
	}

	return shortened, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}
