// Package service contiene la lógica de negocio de la API.
package service

import (
	"context"
	"fmt"

	"github.com/talsory/coding-challenge/api-go/internal/client"
	"github.com/talsory/coding-challenge/api-go/internal/math"
	"github.com/talsory/coding-challenge/api-go/internal/model"
)

// QRService encapsula la lógica de negocio del endpoint QR.
type QRService struct {
	nodeClient *client.NodeClient
}

// NewQRService crea una nueva instancia del servicio.
func NewQRService(nodeClient *client.NodeClient) *QRService {
	return &QRService{
		nodeClient: nodeClient,
	}
}

// QRResponse agrupa todo lo que devolvemos al cliente final.
type QRResponse struct {
	Q          model.Matrix               `json:"q"`
	R          model.Matrix               `json:"r"`
	Statistics *client.StatisticsResponse `json:"statistics"`
}

// ProcessMatrix recibe una matriz, la valida, calcula su factorización QR,
// obtiene las estadísticas desde Node.js y devuelve el resultado consolidado.
func (s *QRService) ProcessMatrix(ctx context.Context, matrix model.Matrix, token string) (*QRResponse, error) {
	// 1. Validar la matriz
	if err := matrix.Validate(); err != nil {
		return nil, fmt.Errorf("matriz inválida: %w", err)
	}

	// 2. Calcular factorización QR
	qrResult, err := math.QR(matrix)
	if err != nil {
		return nil, fmt.Errorf("error en factorización QR: %w", err)
	}

	// 3. Llamar a Node.js para estadísticas (con token)
	stats, err := s.nodeClient.GetStatistics(ctx, qrResult.Q, qrResult.R, token)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo estadísticas de Node.js: %w", err)
	}

	// 4. Consolidar respuesta
	return &QRResponse{
		Q:          qrResult.Q,
		R:          qrResult.R,
		Statistics: stats,
	}, nil
}