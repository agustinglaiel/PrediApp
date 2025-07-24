# Etapa 1: Construcción
FROM golang:1.23-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git

# Copiar el directorio db
COPY db ./db

# Copiar el archivo go.mod y go.sum de la raíz
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY main.go ./

# Compilar la aplicación
RUN go build -o migrator main.go

# Etapa 2: Imagen final
FROM alpine:latest

# Instalar certificados y dependencias mínimas
RUN apk add --no-cache ca-certificates tzdata

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/migrator .

# Copiar el directorio db para las migraciones
COPY --from=builder /app/db ./db

# # Copiar los archivos .env
# COPY .env ./
# COPY .env.stage ./

# Comando para ejecutar las migraciones
CMD ["./migrator", "migrate"]