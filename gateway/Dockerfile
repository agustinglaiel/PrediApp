# Etapa 1: Construcción
FROM golang:1.23-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git

# Copiar los archivos de módulos y descargar dependencias desde el directorio gateway/
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download

# Copiar el resto del código fuente desde el directorio gateway/
COPY gateway/ .

# Compilar la aplicación
RUN go build -o gateway-service ./main.go

# Etapa 2: Imagen final
FROM alpine:latest

# Instalar certificados y dependencias mínimas
RUN apk add --no-cache ca-certificates tzdata

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/gateway-service .

# Copiar los archivos .env desde la raíz del proyecto al directorio padre (/)
COPY .env* /

# Exponer el puerto 8080 (fijo según el código)
EXPOSE 8080

# Comando para ejecutar el microservicio
CMD ["./gateway-service"]