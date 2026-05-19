package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"goproject/internal/domain"
)

type userRepo struct {
	postgresql *sqlx.DB
	log        *slog.Logger
}

func NewUserRepo(db *sqlx.DB, log *slog.Logger) domain.UserRepository {
	return &userRepo{postgresql: db, log: log}
}

func (r *userRepo) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var u domain.User
	err := r.postgresql.QueryRowContext(ctx, createUserQuery, user.Name, user.Email, user.Phone, user.PasswordHash).Scan(
		&u.ID, &u.Name, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		r.log.Error("failed to create user", slog.Any("error", err))
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	err := r.postgresql.QueryRowContext(ctx, getUserByIDQuery, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		// Если юзер просто не найден — это не паника системы, возвращаем nil без ошибки
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.Error("failed to get user by id", slog.Any("error", err)) // Поправил текст лога
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.postgresql.QueryRowContext(ctx, getUserByEmailQuery, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Phone, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		// Вот тут происходил затык! Исправляем:
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Возвращаем nil, nil -> юзкейс поймет, что email свободен
		}
		r.log.Error("failed to get user by email", slog.Any("error", err)) // Поправил текст лога
		return nil, err
	}

	return &u, nil
}