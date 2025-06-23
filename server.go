package main

import (
	"fmt"
	"os"

	"ms-validar-recaptcha/config"
	"ms-validar-recaptcha/handlers"
	"ms-validar-recaptcha/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	// Cargar configuración
	if err := config.LoadConfig(); err != nil {
		fmt.Println("Error al cargar la configuración:", err)
		os.Exit(1)
	}

	e := echo.New()

	// Aplicar middlewares
	e.Use(middleware.LoggingMiddleware())
	e.Use(middleware.CORSMiddleware())

	// Rutas
	e.GET("/", handlers.HealthCheckHandler)
	e.GET("/health", handlers.HealthCheckHandler)
	e.POST("/validate-recaptcha", handlers.ValidateRecaptchaHandler)

	// Obtener puerto y mostrar información
	port := config.GetPort()
	config.PrintServerInfo(port)

	// Iniciar el servidor
	e.Logger.Fatal(e.Start(":" + port))
}
