package repository

import (
    "github.com/jmoiron/sqlx"
    "goproject/internal/domain"
    "context"
    "log/slog"
)

type tariffRepo struct { 
    postgresql *sqlx.DB 
    log *slog.Logger
}

func NewTariffRepo(db *sqlx.DB, log *slog.Logger) domain.TariffRepository {
    return &tariffRepo{postgresql: db, log: log}
}

func (r *tariffRepo) GetTariffByID(ctx context.Context, id int) (*domain.Tariff, error) {
    var t domain.Tariff
	err := r.postgresql.QueryRowContext(ctx, createUserQuery, id).Scan(
		&t.ID, &t.Name, &t.Price,
	)
	if err != nil {
		r.log.Error("failed to load tariff", slog.Any("error", err))
		return nil, err
	}

	return &t, nil
}