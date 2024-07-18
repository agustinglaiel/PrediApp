# Use the official Golang image as the base image
FROM golang:1.19-alpine

# Set the working directory
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git curl

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Expose the port the app runs on
EXPOSE 8080

# Start the application
CMD ["./main"]
