# Etapa 1: Construcción
FROM golang:1.23-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git

# Copiar el directorio db desde la raíz del proyecto a /db
COPY ./db /db

# Copiar los archivos de módulos y descargar dependencias
COPY groups/go.mod groups/go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY groups ./

# Compilar la aplicación
# El binario se genera desde cmd/main.go
RUN go build -o groups-service ./cmd/main.go

# Etapa 2: Imagen final
FROM alpine:latest

# Instalar certificados y dependencias mínimas
RUN apk add --no-cache ca-certificates tzdata

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/groups-service .
COPY --from=builder /app/wait-for-db.sh .
RUN chmod +x wait-for-db.sh

ENTRYPOINT ["./wait-for-db.sh"]

# Comando para ejecutar el microservicio
CMD ["./groups-service"]