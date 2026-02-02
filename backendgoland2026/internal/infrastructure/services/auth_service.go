package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret string
	jwtExpiry time.Duration
}

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService() *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-key-change-this-in-production"
	}

	expiryStr := os.Getenv("JWT_EXPIRATION")
	expiryHours := 24
	if h, err := strconv.Atoi(expiryStr); err == nil {
		expiryHours = h
	}

	return &AuthService{
		jwtSecret: secret,
		jwtExpiry: time.Duration(expiryHours) * time.Hour,
	}
}

// GenerateRSAKeyPair genera un par de claves RSA (privada y pública)
func (s *AuthService) GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating RSA key: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// EncryptPassword cifra una contraseña usando la clave pública RSA
func (s *AuthService) EncryptPassword(password string, publicKey *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(password),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error encrypting password: %w", err)
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// DecryptPassword descifra una contraseña usando la clave privada RSA
func (s *AuthService) DecryptPassword(encryptedPassword string, privateKey *rsa.PrivateKey) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("error decoding encrypted password: %w", err)
	}

	decryptedBytes, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		encryptedBytes,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error decrypting password: %w", err)
	}

	return string(decryptedBytes), nil
}

// GenerateToken genera un JWT token con expiración
func (s *AuthService) GenerateToken(userID int, email string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.jwtExpiry)

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// VerifyToken verifica y parsea un JWT token
func (s *AuthService) VerifyToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Verificar que el token no ha expirado
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// ValidateEmail valida que el email sea válido (contiene @ y .)
func (s *AuthService) ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	// Verificar que contenga @ y .
	hasAt := false
	hasDot := false
	atIndex := -1

	for i, char := range email {
		if char == '@' {
			if hasAt {
				return false // No puede haber dos @
			}
			hasAt = true
			atIndex = i
		} else if char == '.' {
			hasDot = true
		}
	}

	// Verificar estructura válida
	if !hasAt || !hasDot || atIndex == -1 || atIndex == 0 {
		return false
	}

	// Verificar que hay caracteres después del @
	if atIndex >= len(email)-1 {
		return false
	}

	// Verificar que hay un . después del @
	afterAt := email[atIndex+1:]
	if !contains(afterAt, '.') {
		return false
	}

	return true
}

func contains(s string, char rune) bool {
	for _, c := range s {
		if c == char {
			return true
		}
	}
	return false
}
