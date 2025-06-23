package middleware

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// LoggingMiddleware registra las peticiones HTTP
func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Printf("[LOG] %s %s\n", c.Request().Method, c.Request().URL.Path)

			// Log del body para peticiones POST y PUT
			if c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut {
				bodyBytes, _ := io.ReadAll(c.Request().Body)
				fmt.Printf("[LOG] Body: %s\n", string(bodyBytes))
				c.Request().Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			}

			return next(c)
		}
	}
}

// CORSMiddleware configura los headers CORS
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
}
