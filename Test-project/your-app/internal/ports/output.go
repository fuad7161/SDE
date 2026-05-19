package ports

import (
	"context"
	"main/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) error
}
