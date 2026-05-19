package postgres

import (
	"context"
	"main/internal/domain"
	"main/internal/ports"
)

type Repo struct {
}

func (r *Repo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	// DB query logic...
	return &domain.User{
		ID:   id,
		Name: "John Doe",
		Age:  30,
	}, nil
}

func (r *Repo) Save(ctx context.Context, u *domain.User) error {
	// INSERT/UPDATE logic...
	return nil
}

var _ ports.UserRepository = (*Repo)(nil)
