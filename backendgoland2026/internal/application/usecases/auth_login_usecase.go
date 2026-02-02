package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Paola199723/backendgoland2026/internal/infrastructure/persistence"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/services"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
)

type AuthLoginUseCase struct {
	userRepo    *persistence.UserRepository
	authService *services.AuthService
}

func NewAuthLoginUseCase(userRepo *persistence.UserRepository) *AuthLoginUseCase {
	return &AuthLoginUseCase{
		userRepo:    userRepo,
		authService: services.NewAuthService(),
	}
}

// hashPassword genera un hash SHA256 de la contraseña
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (uc *AuthLoginUseCase) Execute(email, password string) (*dto.AuthTokenResponse, error) {
	// Validar email
	if !uc.authService.ValidateEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Calcular hash de la contraseña
	passwordHash := hashPassword(password)

	// Buscar usuario por email
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Si el usuario no existe, crear uno nuevo
	if user == nil {
		// Crear usuario con hash de contraseña
		user, err = uc.userRepo.CreateUser(email, passwordHash)
		if err != nil {
			return nil, err
		}
	} else {
		// Usuario existe, verificar contraseña comparando hashes
		if user.PasswordEncrypted != passwordHash {
			return nil, fmt.Errorf("invalid credentials")
		}
	}

	// Generar token JWT
	token, expiresAt, err := uc.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Guardar token en la BD
	err = uc.userRepo.SaveAuthToken(user.ID, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return &dto.AuthTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		UserID:    user.ID,
		Email:     user.Email,
	}, nil
}
