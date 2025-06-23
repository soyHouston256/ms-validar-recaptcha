package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig carga las variables de entorno desde el archivo .env
func LoadConfig() error {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Advertencia: No se pudo cargar el archivo .env:", err)
		fmt.Println("Asegúrate de que existe un archivo .env con las variables necesarias")
	}

	// Verificar que la variable de entorno esté cargada
	secretKey := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secretKey == "" {
		fmt.Println("ERROR: RECAPTCHA_SECRET_KEY no está configurada en el archivo .env")
		return fmt.Errorf("RECAPTCHA_SECRET_KEY no está configurada")
	} else {
		fmt.Println("✓ RECAPTCHA_SECRET_KEY cargada correctamente")
	}

	return nil
}

// GetPort obtiene el puerto del servidor desde las variables de entorno
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	return port
}

// PrintServerInfo imprime la información del servidor
func PrintServerInfo(port string) {
	fmt.Printf("Servidor iniciado en el puerto %s\n", port)
	fmt.Println("Endpoints disponibles:")
	fmt.Println("- GET  /health - Verificar estado del servicio")
	fmt.Println("- POST /validate-recaptcha - Validar token de reCAPTCHA")
	fmt.Println("\nVariables de entorno requeridas:")
	fmt.Println("- RECAPTCHA_SECRET_KEY: Clave secreta de reCAPTCHA")
	fmt.Println("- PORT: Puerto del servidor (opcional, default: 1323)")
}
