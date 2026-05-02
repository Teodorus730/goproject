package waterService

import (
	"context"
	"net/http"
)

type Handler interface {
	Test() http.HandlerFunc
}

type UseCase interface {
	GetTest(ctx context.Context) (string, error)
}

type Repository interface {
	GetTest(ctx context.Context) (string, error)
}