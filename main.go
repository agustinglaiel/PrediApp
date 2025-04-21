package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var cmds []*exec.Cmd // Variable global para almacenar los comandos en ejecución

// Función para ejecutar main.go en cada servicio
func runMain(dir string, envVars map[string]string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = dir

	// Configurar variables de entorno
	env := os.Environ()
	for key, value := range envVars {
		env = append(env, key+"="+value)
	}
	cmd.Env = env

	// Conectar la salida del proceso a la terminal
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	// Ejecutar el comando en segundo plano
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start service in %s: %v", dir, err)
		return err
	}

	cmds = append(cmds, cmd) // Agregar el comando en ejecución a la lista
	return nil
}

// Función para matar todos los procesos en ejecución
func cleanup() {
	for _, cmd := range cmds {
		if cmd != nil && cmd.Process != nil {
			log.Printf("Killing process for %s (PID: %d)", cmd.Path, cmd.Process.Pid)
			cmd.Process.Kill() // Matar el proceso
		}
	}
	log.Println("All services stopped.")
}

func main() {
	// Construir la ruta al archivo .env
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
		os.Exit(1)
	}
	envPath := filepath.Join(currentDir, ".env")

	// Cargar el archivo .env para obtener el valor de ENV
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Obtener el valor de ENV desde el archivo .env (o desde la variable de entorno si está definida)
	env := os.Getenv("ENV")
	if env == "" {
		log.Println("ENV not set in .env file, defaulting to 'stage'")
		env = "stage"
	}

	// Construir la ruta al archivo .env.stage o .env.prod según el valor de ENV
	envFile := fmt.Sprintf(".env.%s", env)
	envSpecificPath := filepath.Join(currentDir, envFile)

	// Cargar las variables de entorno específicas del ambiente
	err = godotenv.Load(envSpecificPath)
	if err != nil {
		log.Printf("Error loading %s file: %v", envFile, err)
		log.Println("Continuing with variables from .env or system environment")
	}

	// Lista de servicios con sus respectivas variables de entorno para el puerto, en orden
	services := []struct {
		dir        string
		serviceURL string
	}{
		{"./users/cmd", os.Getenv("USERS_SERVICE_URL")},
		{"./sessions/cmd", os.Getenv("SESSIONS_SERVICE_URL")},
		{"./drivers/cmd", os.Getenv("DRIVERS_SERVICE_URL")},
		{"./results/cmd", os.Getenv("RESULTS_SERVICE_URL")},
		{"./prodes/cmd", os.Getenv("PRODES_SERVICE_URL")},
		{"./groups/cmd", os.Getenv("GROUPS_SERVICE_URL")},
	}

	// Obtener todas las variables de entorno actualizadas
	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	// Ejecutar cada main.go secuencialmente con un retraso de 2 segundos
	for i, service := range services {
		// Extraer el puerto de la URL del microservicio
		parts := strings.Split(service.serviceURL, ":")
		if len(parts) < 3 {
			log.Printf("Invalid service URL for %s: %s", service.dir, service.serviceURL)
			continue
		}
		port := parts[2]

		// Configurar las variables de entorno para este microservicio
		serviceEnvVars := make(map[string]string)
		for k, v := range envVars {
			serviceEnvVars[k] = v
		}
		serviceEnvVars["PORT"] = port

		// Iniciar el microservicio
		log.Printf("Starting service in %s on port %s...", service.dir, port)
		err := runMain(service.dir, serviceEnvVars)
		if err != nil {
			log.Printf("Failed to run main.go in %s on port %s: %v", service.dir, port, err)
		}

		// Esperar 2 segundos antes de iniciar el siguiente servicio (excepto el último)
		if i < len(services)-1 {
			log.Printf("Waiting 2 seconds before starting the next service...")
			time.Sleep(1 * time.Second)
		}
	}

	// Capturar señales de interrupción (Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Esperar la señal de interrupción
	<-sigs
	log.Println("Interrupt signal received. Stopping services...")

	// Limpiar procesos
	cleanup()
}