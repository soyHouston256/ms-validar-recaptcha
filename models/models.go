package models

// RecaptchaRequest representa la estructura de la petición
// swagger:model RecaptchaRequest
type RecaptchaRequest struct {
	// Token de reCAPTCHA a validar
	Token string `json:"token"`
}

// RecaptchaResponse representa la respuesta de Google reCAPTCHA
// swagger:model RecaptchaResponse
type RecaptchaResponse struct {
	// Indica si la validación fue exitosa
	Success bool `json:"success"`
	// Score de reCAPTCHA v3
	Score float64 `json:"score"`
	// Acción asociada al token
	Action string `json:"action"`
	// Timestamp del desafío
	ChallengeTS string `json:"challenge_ts"`
	// Hostname donde se resolvió el reto
	Hostname string `json:"hostname"`
	// Códigos de error (si existen)
	ErrorCodes []string `json:"error-codes,omitempty"`
}

// APIResponse representa la respuesta estándar de nuestro API
// swagger:model APIResponse
type APIResponse struct {
	// Datos de la respuesta
	Data interface{} `json:"data,omitempty"`
	// Indica si la operación fue exitosa
	Success bool `json:"success"`
	// Mensaje de error (null si no hay error)
	ErrorMessage *string `json:"errorMessage"`
}
