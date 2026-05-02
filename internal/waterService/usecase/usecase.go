package usecase 

import (
	"context"
	waterService "goproject/internal/waterService"
)

type useCase struct {
	repo waterService.Repository
} 

func New(repo waterService.Repository) *useCase {
	return &useCase{repo:repo}
}

func (uc *useCase) GetTest(ctx context.Context) (string, error) {
	return uc.repo.GetTest(ctx)
}

