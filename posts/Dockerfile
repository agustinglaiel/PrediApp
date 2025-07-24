# Etapa 1: Construcción
FROM golang:1.23-alpine AS builder

# Establecer el directorio de trabajo
WORKDIR /app

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git

# Copiar el directorio db desde la raíz del proyecto a /db
COPY ./db /db

# Copiar los archivos de módulos y descargar dependencias
COPY posts/go.mod posts/go.sum ./
RUN go mod download

# Copiar el código fuente del microservicio posts
COPY posts ./

# Compilar la aplicación
RUN go build -o posts-service ./cmd/main.go

# Etapa 2: Imagen final
FROM alpine:latest

# Instalar certificados y dependencias mínimas
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/posts-service .
COPY --from=builder /app/wait-for-db.sh .
RUN chmod +x wait-for-db.sh

ENTRYPOINT ["./wait-for-db.sh"]

# Comando para ejecutar el microservicio
CMD ["./posts-service"]