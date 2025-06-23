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
func ValidateRecaptchaHandler(c echo.Context) error {
	// Verificar que sea una petición POST
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, models.APIResponse{
			Success: false,
			Message: "Método no permitido. Use POST",
		})
	}

	// Parsear el body de la petición
	var req models.RecaptchaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Error al parsear la petición: " + err.Error(),
		})
	}

	// Validar que el token no esté vacío
	if strings.TrimSpace(req.Token) == "" {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "El token de reCAPTCHA es requerido",
		})
	}

	// Validar el token con Google
	recaptchaResp, err := validateRecaptcha(req.Token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error al validar reCAPTCHA: " + err.Error(),
		})
	}

	// Verificar si la validación fue exitosa
	if !recaptchaResp.Success {
		errorMsg := "Validación de reCAPTCHA fallida"
		if len(recaptchaResp.ErrorCodes) > 0 {
			errorMsg += ": " + strings.Join(recaptchaResp.ErrorCodes, ", ")
		}

		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: errorMsg,
			Data:    recaptchaResp,
		})
	}

	// Para reCAPTCHA v3, también verificar el score (opcional)
	// Un score de 0.5 o mayor generalmente se considera válido
	if recaptchaResp.Score < 0.5 {
		return c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Score de reCAPTCHA muy bajo. Posible bot detectado",
			Data:    recaptchaResp,
		})
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "reCAPTCHA validado exitosamente",
		Data:    recaptchaResp,
	})
}

// HealthCheckHandler para verificar que el servicio esté funcionando
func HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Servicio de validación de reCAPTCHA funcionando correctamente",
	})
}
