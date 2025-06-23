package models

// RecaptchaRequest representa la estructura de la petición
type RecaptchaRequest struct {
	Token string `json:"token"`
}

// RecaptchaResponse representa la respuesta de Google reCAPTCHA
type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

// APIResponse representa la respuesta de nuestro API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
