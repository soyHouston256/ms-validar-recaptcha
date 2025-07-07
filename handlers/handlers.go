package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"ms-validar-recaptcha/models"

	"github.com/labstack/echo/v4"
)

// createSuccessResponse crea una respuesta exitosa estándar
func createSuccessResponse(data interface{}) models.APIResponse {
	return models.APIResponse{
		Data:         data,
		Success:      true,
		ErrorMessage: nil,
	}
}

// createErrorResponse crea una respuesta de error estándar
func createErrorResponse(errorMessage string) models.APIResponse {
	return models.APIResponse{
		Data:         nil,
		Success:      false,
		ErrorMessage: &errorMessage,
	}
}

// validateRecaptcha valida el token de reCAPTCHA con Google
func validateRecaptcha(token string) (*models.RecaptchaResponse, error) {
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		return nil, fmt.Errorf("RECAPTCHA_SECRET_KEY no está configurada")
	}

	// Preparar los datos para enviar a Google
	data := url.Values{}
	data.Set("secret", secretKey)
	data.Set("response", token)

	// Hacer la petición a Google
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", data)
	if err != nil {
		return nil, fmt.Errorf("error al hacer petición a Google: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer respuesta: %v", err)
	}

	// Parsear la respuesta JSON
	var recaptchaResp models.RecaptchaResponse
	if err := json.Unmarshal(body, &recaptchaResp); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta JSON: %v", err)
	}

	return &recaptchaResp, nil
}

// ValidateRecaptchaHandler maneja las peticiones de validación de reCAPTCHA
// @Summary Validar token de reCAPTCHA
// @Description Valida un token de reCAPTCHA con Google y retorna el resultado de la validación
// @Tags reCAPTCHA
// @Accept json
// @Produce json
// @Param request body models.RecaptchaRequest true "Token de reCAPTCHA a validar"
// @Success 200 {object} models.APIResponse "Token válido"
// @Failure 400 {object} models.APIResponse "Token inválido o error en la petición"
// @Failure 405 {object} models.APIResponse "Método no permitido"
// @Failure 500 {object} models.APIResponse "Error interno del servidor"
// @Router /validate-recaptcha [post]
func ValidateRecaptchaHandler(c echo.Context) error {
	// Verificar que sea una petición POST
	if c.Request().Method != http.MethodPost {
		errorMsg := "Método no permitido. Use POST"
		return c.JSON(http.StatusMethodNotAllowed, createErrorResponse(errorMsg))
	}

	// Parsear el body de la petición
	var req models.RecaptchaRequest
	if err := c.Bind(&req); err != nil {
		errorMsg := "Error al parsear la petición: " + err.Error()
		return c.JSON(http.StatusBadRequest, createErrorResponse(errorMsg))
	}

	// Validar que el token no esté vacío
	if strings.TrimSpace(req.Token) == "" {
		errorMsg := "El token de reCAPTCHA es requerido"
		return c.JSON(http.StatusBadRequest, createErrorResponse(errorMsg))
	}

	// Validar el token con Google
	recaptchaResp, err := validateRecaptcha(req.Token)
	if err != nil {
		errorMsg := "Error al validar reCAPTCHA: " + err.Error()
		return c.JSON(http.StatusInternalServerError, createErrorResponse(errorMsg))
	}

	// Verificar si la validación fue exitosa
	if !recaptchaResp.Success {
		errorMsg := "Validación de reCAPTCHA fallida"
		if len(recaptchaResp.ErrorCodes) > 0 {
			errorMsg += ": " + strings.Join(recaptchaResp.ErrorCodes, ", ")
		}

		// Incluir los datos de la respuesta de Google en caso de error
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Data:         recaptchaResp,
			Success:      false,
			ErrorMessage: &errorMsg,
		})
	}

	// Para reCAPTCHA v3, también verificar el score (opcional)
	// Un score de 0.5 o mayor generalmente se considera válido
	if recaptchaResp.Score < 0.5 {
		errorMsg := "Score de reCAPTCHA muy bajo. Posible bot detectado"
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Data:         recaptchaResp,
			Success:      false,
			ErrorMessage: &errorMsg,
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, createSuccessResponse(recaptchaResp))
}

// HealthCheckHandler para verificar que el servicio esté funcionando
// @Summary Verificar estado del servicio
// @Description Verifica que el microservicio esté funcionando correctamente
// @Tags Health
// @Produce json
// @Success 200 {object} models.APIResponse "Servicio funcionando"
// @Router /health [get]
// @Router / [get]
func HealthCheckHandler(c echo.Context) error {
	healthData := map[string]interface{}{
		"status":    "OK",
		"service":   "reCAPTCHA Validation Service",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	return c.JSON(http.StatusOK, createSuccessResponse(healthData))
}
