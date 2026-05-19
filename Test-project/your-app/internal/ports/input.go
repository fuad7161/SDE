package ports

import (
	"context"
	"main/internal/domain"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, name string, age int) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
}
