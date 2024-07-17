# Usar una imagen base oficial de Go
FROM golang:1.18-alpine

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el archivo go.mod y go.sum para gestionar dependencias
COPY go.mod go.sum ./

# Descargar las dependencias de Go. Esto se beneficia de la caché de las capas de Docker si los archivos go.mod y go.sum no cambian
RUN go mod download

# Copiar el código fuente del microservicio al contenedor
COPY . .

# Compilar la aplicación al binario ejecutable. Reemplaza 'cmd/main.go' con la ruta al punto de entrada de tu aplicación
RUN go build -o /users-service ./cmd/main.go

# Indicar el puerto que el contenedor debe exponer (cambia 8080 por el puerto que uses)
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["/users-service"]
