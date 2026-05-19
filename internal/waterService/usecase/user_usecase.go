package usecase

import (
    "fmt"
    "goproject/internal/domain"
    "golang.org/x/crypto/bcrypt"
	"context"
	"log/slog"
)


type userUsecase struct {
    userRepo domain.UserRepository
    sessionRepo domain.SessionRepository
    log *slog.Logger
}

func NewUserUsecase(ur domain.UserRepository, sr domain.SessionRepository, log *slog.Logger) domain.UserUsecase {
    return &userUsecase{userRepo: ur, sessionRepo: sr, log: log}
}


func (u *userUsecase) CreateUser(ctx context.Context, name, email, phone, password string) (*domain.User, error) {
	existingUser, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		u.log.Error("Failed to check existing user", slog.Any("email", email), slog.Any("error", err))
		return nil, err
	}

	if existingUser != nil {
		u.log.Warn("User already exists", slog.Any("email", email))
		return nil, fmt.Errorf("user already exists")
	}

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        u.log.Warn("failed to hash password", slog.Any("email", email))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

    user := &domain.User{Name: name, Email: email, PasswordHash: string(hashedPassword), Phone: phone}

	created_user, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		u.log.Error("Failed to register user", slog.Any("email", email), slog.Any("error", err))
		return nil, err
	}

	return created_user, nil
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		u.log.Error("Failed to get user by login", slog.Any("email", email), slog.Any("error", err))
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) CreateSession(ctx context.Context, sessionID, userID string) error {
	err := u.sessionRepo.CreateSession(ctx, sessionID, userID)
	if err != nil {
		u.log.Error("Failed to create session", slog.String("userID", userID), slog.Any("error", err))
		return err
	}
	return nil
}


