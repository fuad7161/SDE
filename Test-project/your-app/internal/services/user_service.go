package usecases

import (
	"context"
	"main/internal/domain"
	"main/internal/ports"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, name string, age int) (*domain.User, error) {
	u := &domain.User{Name: name, Age: age}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	u.ID = "gen-" + name
	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

var _ ports.UserUseCase = (*userService)(nil)
