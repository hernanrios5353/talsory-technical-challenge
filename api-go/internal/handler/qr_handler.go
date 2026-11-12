// Package handler contiene los manejadores HTTP de la API.
package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/talsory/coding-challenge/api-go/internal/model"
	"github.com/talsory/coding-challenge/api-go/internal/service"
)

// QRHandler maneja las peticiones relacionadas con la factorización QR.
type QRHandler struct {
	service *service.QRService
}

// NewQRHandler crea una nueva instancia del handler.
func NewQRHandler(svc *service.QRService) *QRHandler {
	return &QRHandler{service: svc}
}

// QRRequest define el contrato JSON de entrada.
type QRRequest struct {
	Matrix model.Matrix `json:"matrix"`
}

// ErrorResponse define el contrato JSON de error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Handle procesa la petición POST /api/v1/qr.
func (h *QRHandler) Handle(c *fiber.Ctx) error {
	// 1. Parsear el body
	var req QRRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "JSON inválido: " + err.Error(),
		})
	}

	// 2. Validar que la matriz no sea nil
	if req.Matrix == nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "el campo 'matrix' es obligatorio",
		})
	}

	// 3. Crear un contexto con timeout para la operación completa
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	// Extraer el token del header (ya validado por el middleware)
	authHeader := c.Get("Authorization")
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// 4. Llamar al servicio
	response, err := h.service.ProcessMatrix(ctx, req.Matrix, token)
	if err != nil {
		return h.handleError(c, err) // <-- AQUÍ ESTABA EL ERROR, FALTABA ESTO
	}

	// 5. Devolver respuesta exitosa
	return c.Status(fiber.StatusOK).JSON(response)
}

// handleError clasifica los errores y devuelve el status HTTP adecuado.
func (h *QRHandler) handleError(c *fiber.Ctx, err error) error {
	// Errores de validación de la matriz → 400 Bad Request
	if errors.Is(err, model.ErrEmptyMatrix) ||
		errors.Is(err, model.ErrIrregularMatrix) ||
		errors.Is(err, model.ErrInvalidDimensions) ||
		errors.Is(err, model.ErrMoreColumnsThanRows) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	// Errores de timeout o cancelación → 504 Gateway Timeout
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return c.Status(fiber.StatusGatewayTimeout).JSON(ErrorResponse{
			Error: "la operación excedió el tiempo límite",
		})
	}

	// Si el error viene de Node.js → 502 Bad Gateway
	if contains(err.Error(), "Node.js") {
		return c.Status(fiber.StatusBadGateway).JSON(ErrorResponse{
			Error: "error comunicándose con el servicio de estadísticas",
		})
	}

	// Cualquier otro error → 500 Internal Server Error
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Error: "error interno del servidor",
	})
}

// contains es un helper simple para verificar subcadenas.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}