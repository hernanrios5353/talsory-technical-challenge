package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/talsory/coding-challenge/api-go/internal/model"
	"github.com/talsory/coding-challenge/api-go/internal/service"
)

// AuthHandler maneja las peticiones de autenticación.
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler crea una nueva instancia del handler.
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

// Login procesa POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "JSON inválido: " + err.Error(),
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "username y password son obligatorios",
		})
	}

	resp, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
				Error: "credenciales inválidas",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "error interno generando el token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}