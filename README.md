# Microservicio de Validación de reCAPTCHA

Este microservicio proporciona un endpoint REST para validar tokens de Google reCAPTCHA v2 y v3, construido con una arquitectura modular en Go.

## Características

- ✅ Validación de reCAPTCHA v2 y v3
- ✅ Arquitectura modular y escalable
- ✅ Manejo de errores robusto
- ✅ Respuestas JSON estructuradas
- ✅ CORS configurado para frontend
- ✅ Health check endpoint
- ✅ Configuración mediante archivo `.env`
- ✅ Middlewares reutilizables
- ✅ Separación clara de responsabilidades

## Requisitos

- Go 1.23.0 o superior
- Clave secreta de reCAPTCHA de Google

## Instalación

1. Clona el repositorio:
```bash
git clone <tu-repositorio>
cd ms-validar-recaptcha
```

2. Instala las dependencias:
```bash
go mod tidy
```

3. Configura las variables de entorno:
```bash
# Copia el archivo de ejemplo
cp env.example .env

# Edita el archivo .env con tus credenciales
nano .env
```

4. Ejecuta el servidor:
```bash
go run server.go
```

## Configuración

### Archivo .env

Crea un archivo `.env` en la raíz del proyecto con el siguiente contenido:

```env
# Clave secreta de reCAPTCHA (obligatoria)
# Obtén tu clave en: https://www.google.com/recaptcha/admin
RECAPTCHA_SECRET_KEY=tu_clave_secreta_aqui

# Puerto del servidor (opcional, default: 1323)
PORT=1323
```

### Variables de entorno

| Variable | Descripción | Requerida | Default |
|----------|-------------|-----------|---------|
| `RECAPTCHA_SECRET_KEY` | Clave secreta de reCAPTCHA | ✅ | - |
| `PORT` | Puerto del servidor | ❌ | 1323 |

## Arquitectura del Proyecto

El proyecto sigue una arquitectura modular con separación clara de responsabilidades:

```
ms-validar-recaptcha/
├── config/
│   └── config.go          # Configuración y variables de entorno
├── handlers/
│   └── handlers.go        # Lógica de manejo de peticiones HTTP
├── middleware/
│   └── middleware.go      # Middlewares (logging, CORS)
├── models/
│   └── models.go          # Estructuras de datos
├── server.go              # Punto de entrada principal
├── go.mod                 # Dependencias del módulo
├── go.sum                 # Checksums de dependencias
├── .env                   # Variables de entorno (crear desde env.example)
├── env.example           # Ejemplo de variables de entorno
├── Dockerfile            # Configuración de Docker
└── README.md             # Documentación
```

### Módulos

- **`config/`**: Maneja la configuración del servidor y carga de variables de entorno
- **`handlers/`**: Contiene la lógica de negocio para manejar las peticiones HTTP
- **`middleware/`**: Middlewares reutilizables (logging, CORS)
- **`models/`**: Estructuras de datos y tipos compartidos

## Endpoints

### 1. Health Check
```http
GET /health
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Servicio de validación de reCAPTCHA funcionando correctamente"
}
```

### 2. Validar reCAPTCHA
```http
POST /validate-recaptcha
Content-Type: application/json

{
  "token": "token_del_recaptcha_del_frontend"
}
```

**Respuesta exitosa:**
```json
{
  "success": true,
  "message": "reCAPTCHA validado exitosamente",
  "data": {
    "success": true,
    "score": 0.9,
    "action": "submit",
    "challenge_ts": "2024-01-01T12:00:00Z",
    "hostname": "example.com"
  }
}
```

**Respuesta de error:**
```json
{
  "success": false,
  "message": "Validación de reCAPTCHA fallida: invalid-input-response",
  "data": {
    "success": false,
    "error-codes": ["invalid-input-response"]
  }
}
```

## Códigos de estado HTTP

| Código | Descripción |
|--------|-------------|
| 200 | Validación exitosa |
| 400 | Token inválido o faltante |
| 405 | Método HTTP no permitido |
| 500 | Error interno del servidor |

## Ejemplos de uso

### Frontend (JavaScript)

```javascript
// Ejemplo con reCAPTCHA v2
grecaptcha.ready(function() {
  grecaptcha.execute('SITE_KEY', {action: 'submit'}).then(function(token) {
    fetch('http://localhost:1323/validate-recaptcha', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ token: token })
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        console.log('reCAPTCHA validado:', data.data);
        // Continuar con el formulario
      } else {
        console.error('Error:', data.message);
      }
    });
  });
});
```

### cURL

```bash
curl -X POST http://localhost:1323/validate-recaptcha \
  -H "Content-Type: application/json" \
  -d '{"token": "tu_token_aqui"}'
```

### Python

```python
import requests

response = requests.post(
    'http://localhost:1323/validate-recaptcha',
    json={'token': 'tu_token_aqui'}
)

if response.status_code == 200:
    data = response.json()
    if data['success']:
        print('reCAPTCHA válido')
    else:
        print('Error:', data['message'])
else:
    print('Error HTTP:', response.status_code)
```

## Docker

### Construir imagen
```bash
docker build -t ms-validar-recaptcha .
```

### Ejecutar contenedor
```bash
docker run -p 1323:1323 \
  -e RECAPTCHA_SECRET_KEY="tu_clave_secreta" \
  ms-validar-recaptcha
```

## Desarrollo

### Estructura del proyecto
```
ms-validar-recaptcha/
├── config/
│   └── config.go          # Configuración y variables de entorno
├── handlers/
│   └── handlers.go        # Lógica de manejo de peticiones HTTP
├── middleware/
│   └── middleware.go      # Middlewares (logging, CORS)
├── models/
│   └── models.go          # Estructuras de datos
├── server.go              # Punto de entrada principal
├── go.mod                 # Dependencias de Go
├── go.sum                 # Checksums de dependencias
├── .env                   # Variables de entorno
├── env.example           # Ejemplo de variables de entorno
├── Dockerfile            # Configuración de Docker
└── README.md             # Documentación
```

### Ejecutar tests
```bash
go test ./...
```

### Formatear código
```bash
go fmt ./...
```

### Agregar nuevos middlewares

Para agregar un nuevo middleware, crea una función en `middleware/middleware.go`:

```go
func NewMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Tu lógica aquí
            return next(c)
        }
    }
}
```

### Agregar nuevos handlers

Para agregar un nuevo handler, crea una función en `handlers/handlers.go`:

```go
func NewHandler(c echo.Context) error {
    // Tu lógica aquí
    return c.JSON(http.StatusOK, models.APIResponse{
        Success: true,
        Message: "Nuevo endpoint funcionando",
    })
}
```

## Seguridad

- ✅ Nunca expongas la clave secreta en el frontend
- ✅ Usa HTTPS en producción
- ✅ Valida el hostname en la respuesta de Google
- ✅ Considera implementar rate limiting
- ✅ Monitorea los scores de reCAPTCHA v3
- ✅ El archivo `.env` está en `.gitignore` para proteger credenciales

## Troubleshooting

### Error: "RECAPTCHA_SECRET_KEY no está configurada"
Solución: 
1. Verifica que el archivo `.env` existe
2. Asegúrate de que `RECAPTCHA_SECRET_KEY` esté definida en el archivo `.env`
3. Reinicia el servidor después de modificar el archivo `.env`

### Error: "No se pudo cargar el archivo .env"
Solución: 
1. Copia `env.example` a `.env`: `cp env.example .env`
2. Edita el archivo `.env` con tus credenciales reales

### Error: "invalid-input-response"
Solución: Verifica que el token sea válido y no haya expirado

### Error: "Score de reCAPTCHA muy bajo"
Solución: Para reCAPTCHA v3, ajusta el umbral de score según tus necesidades

## Licencia

MIT License
