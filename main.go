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

    // Cargar variables de entorno desde .env
    err = godotenv.Load(envPath)
    if err != nil {
        fmt.Println("Error al cargar el archivo .env")
        log.Fatalf("Error loading .env file: %v", err)
        os.Exit(1)
    }

    // Lista de directorios con sus respectivos puertos
    services := map[string]string{
        "./drivers/cmd":   "",
        "./prodes/cmd":    "",
        "./results/cmd":   "",
        "./sessions/cmd":  "",
        "./users/cmd":     "",
        "./groups/cmd":    "",
        // "./gateway":       "8080",
    }

    // Obtener todas las variables de entorno
    envVars := make(map[string]string)
    for _, env := range os.Environ() {
        parts := strings.SplitN(env, "=", 2)
        if len(parts) == 2 {
            envVars[parts[0]] = parts[1]
        }
    }

    // Ejecutar cada main.go en paralelo
    for dir, port := range services {
        go func(d string, port string, envVars map[string]string) {
            if port != "" {
                err := runMain(d, envVars)
                if err != nil {
                    log.Printf("Failed to run main.go in %s on port %s: %v", d, port, err)
                }
            } else {
                err := runMain(d, envVars)
                if err != nil {
                    log.Printf("Failed to run main.go in %s: %v", d, err)
                }
            }
        }(dir, port, envVars)
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