// Package client contiene los clientes HTTP para comunicarse
// con otros servicios (en este caso, la API de Node.js).
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/talsory/coding-challenge/api-go/internal/model"
)

// NodeClient encapsula la comunicación con la API de Node.js.
type NodeClient struct {
	baseURL    string
	httpClient *http.Client
}

// StatisticsRequest es el payload que enviamos a Node.js.
type StatisticsRequest struct {
	Q model.Matrix `json:"q"`
	R model.Matrix `json:"r"`
}

// StatisticsResponse es la respuesta que esperamos de Node.js.
// Los campos coinciden con el contrato JSON definido.
type StatisticsResponse struct {
	Q      MatrixStats `json:"q"`
	R      MatrixStats `json:"r"`
	Global MatrixStats `json:"global"`
}

// MatrixStats representa las estadísticas calculadas por Node.js.
type MatrixStats struct {
	Max        float64 `json:"max"`
	Min        float64 `json:"min"`
	Average    float64 `json:"average"`
	Sum        float64 `json:"sum"`
	IsDiagonal bool    `json:"isDiagonal"`
}

// NewNodeClient crea una nueva instancia del cliente.
// Configura un timeout para evitar bloqueos si el servicio no responde.
func NewNodeClient(baseURL string) *NodeClient {
	return &NodeClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Timeout crítico para la entrevista
		},
	}
}

// GetStatistics envía las matrices Q y R a Node.js y recibe las estadísticas.
// Incluye el token JWT para autenticarse contra la API de Node.
func (c *NodeClient) GetStatistics(ctx context.Context, q, r model.Matrix, token string) (*StatisticsResponse, error) {
	payload := StatisticsRequest{
		Q: q,
		R: r,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/statistics", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // <-- Enviamos el token

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en la llamada HTTP a Node.js: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Node.js respondió con status %d: %s", resp.StatusCode, string(body))
	}

	var statsResp StatisticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&statsResp); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de Node.js: %w", err)
	}

	return &statsResp, nil
}