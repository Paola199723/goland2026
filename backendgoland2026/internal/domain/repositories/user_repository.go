package repositories

import "github.com/Paola199723/backendgoland2026/internal/domain/entities"

type UserRepository interface {
	Create(user entities.User) error
	FindByEmail(email string) (*entities.User, error)
}
