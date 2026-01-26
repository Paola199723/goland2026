package persistence

import (
	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/db"
)

type UserRepositoryImpl struct{}

func NewUserRepository() *UserRepositoryImpl {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Create(user entities.User) error {
	_, err := db.DB.Exec(
		"INSERT INTO users (id, name, email) VALUES ($1,$2,$3)",
		user.ID, user.Name, user.Email,
	)
	return err
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*entities.User, error) {
	row := db.DB.QueryRow("SELECT id, name, email FROM users WHERE email=$1", email)

	var u entities.User
	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
