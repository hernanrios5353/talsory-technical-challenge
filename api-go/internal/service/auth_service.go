package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/talsory/coding-challenge/api-go/internal/model"
)

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrJWTSecretMissing   = errors.New("JWT_SECRET no está configurado")
)

// AuthService maneja la lógica de autenticación.
type AuthService struct {
	jwtSecret string
}

// NewAuthService crea una nueva instancia del servicio.
func NewAuthService() *AuthService {
	return &AuthService{
		jwtSecret: os.Getenv("JWT_SECRET"),
	}
}

// Login valida las credenciales y genera un token JWT.
func (s *AuthService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	if s.jwtSecret == "" {
		return nil, ErrJWTSecretMissing
	}

	expectedUser := os.Getenv("AUTH_USER")
	expectedPassword := os.Getenv("AUTH_PASSWORD")

	// Validación simple de credenciales (hardcodeadas en env para el reto)
	if req.Username != expectedUser || req.Password != expectedPassword {
		return nil, ErrInvalidCredentials
	}

	// Generar token JWT con expiración de 1 hora
	claims := jwt.MapClaims{
		"sub": req.Username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{Token: signedToken}, nil
}