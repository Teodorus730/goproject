package repository

import (
	"context"
)

type repo struct{}

func New() *repo {
	return &repo{}
}

func (rep *repo) GetTest(ctx context.Context) (string, error) {
	return "QQ", nil
}