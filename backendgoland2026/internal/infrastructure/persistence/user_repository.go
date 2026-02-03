package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/db"
)

type User struct {
	ID                int
	Email             string
	PasswordEncrypted string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: db.GetConnection()}
}

// Implementar interfaz: Create
func (r *UserRepository) Create(user entities.User) error {
	query := `
		INSERT INTO users (email, password_encrypted, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`
	_, err := r.db.Exec(query, user.Email, user.Name) // Note: usando Name como password por compatibilidad
	return err
}

// Implementar interfaz: FindByEmail
func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id, email FROM users WHERE email = $1
	`
	var user entities.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser crea un nuevo usuario en la base de datos
func (r *UserRepository) CreateUser(email, encryptedPassword string) (*User, error) {
	query := `
		INSERT INTO users (email, password_encrypted, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, email, password_encrypted, created_at, updated_at
	`

	var user User
	err := r.db.QueryRow(query, email, encryptedPassword).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordEncrypted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("error creating user: %w", err)
		}
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail obtiene un usuario por su email
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, password_encrypted, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordEncrypted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuario no existe
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

// GetUserByID obtiene un usuario por su ID
func (r *UserRepository) GetUserByID(userID int) (*User, error) {
	query := `
		SELECT id, email, password_encrypted, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordEncrypted,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Usuario no existe
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

// SaveAuthToken guarda un token de autenticación en la base de datos
func (r *UserRepository) SaveAuthToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO auth_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
	`

	_, err := r.db.Exec(query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("error saving auth token: %w", err)
	}

	return nil
}

// GetAuthToken obtiene un token de autenticación
func (r *UserRepository) GetAuthToken(token string) (int, error) {
	query := `
		SELECT user_id FROM auth_tokens
		WHERE token = $1 AND expires_at > NOW()
	`

	var userID int
	err := r.db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("token not found or expired")
		}
		return 0, fmt.Errorf("error getting auth token: %w", err)
	}

	return userID, nil
}

// DeleteExpiredTokens elimina tokens expirados
func (r *UserRepository) DeleteExpiredTokens() error {
	query := `DELETE FROM auth_tokens WHERE expires_at <= NOW()`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error deleting expired tokens: %w", err)
	}
	return nil
}
