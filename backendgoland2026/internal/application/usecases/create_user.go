package usecases

import (
	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/Paola199723/backendgoland2026/internal/domain/repositories"
)

type CreateUserUseCase struct {
	repo repositories.UserRepository
}

func NewCreateUserUseCase(r repositories.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: r}
}

func (uc *CreateUserUseCase) Execute(user entities.User) error {
	return uc.repo.Create(user)
}
