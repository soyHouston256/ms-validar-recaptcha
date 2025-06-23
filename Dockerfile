# Usar la imagen oficial de Go
FROM golang:1.23-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos de dependencias
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Imagen final más pequeña
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/main .

# Exponer el puerto
EXPOSE 1323

# Comando para ejecutar la aplicación
CMD ["./main"] 