package model

// LoginRequest define el contrato JSON para el login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse define el contrato JSON de respuesta del login.
type LoginResponse struct {
	Token string `json:"token"`
}