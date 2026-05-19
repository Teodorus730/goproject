package repository

import (
    "github.com/jmoiron/sqlx"
    "goproject/internal/domain"
    "context"
    "log/slog"
	"fmt"
)

type sessionRepo struct { 
    postgresql *sqlx.DB 
    log *slog.Logger
}

func NewSessionRepo(db *sqlx.DB, log *slog.Logger) domain.SessionRepository {
    return &sessionRepo{postgresql: db, log: log}
}

func (r *sessionRepo) CreateSession(ctx context.Context, sessionID, userID string) error {
	_, err := r.postgresql.ExecContext(ctx, createSessionQuery, sessionID, userID)
	if err != nil {
		r.log.Error("failed to create session", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *sessionRepo) GetSessionByUserID(ctx context.Context, userID string) (*domain.Session, error) {
	session := &domain.Session{}
	err := r.postgresql.GetContext(ctx, session, getSessionByUserIDQuery, userID)
	if err != nil {
		r.log.Error("failed to get session by user id", slog.Any("error", err))
		return nil, err
	}
	return session, nil
}

func (r *sessionRepo) UpdateSessionExpiry(ctx context.Context, sessionID string) error {
	result, err := r.postgresql.ExecContext(ctx, updateSessionExpiryQuery, sessionID)
	if err != nil {
		r.log.Error("failed to update session expiry", slog.Any("error", err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Error("failed to get rows affected", slog.Any("error", err))
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no session found with the provided ID")
	}

	return nil
}