package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var cmds []*exec.Cmd // Variable global para almacenar los comandos en ejecución

// Función para ejecutar main.go en cada servicio
func runMain(dir string, port string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = dir

	// Configurar variables de entorno
	env := os.Environ()
	cmd.Env = append(env, "PORT="+port)

	// Conectar la salida del proceso a la terminal
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	// Ejecutar el comando en segundo plano
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start service in %s on port %s: %v", dir, port, err)
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
	// Lista de directorios con sus respectivos puertos
	services := map[string]string{
		"./drivers/cmd":   "8051",
		"./prodes/cmd":    "8054",
		"./results/cmd":   "8055",
		"./sessions/cmd":  "8056",
		"./users/cmd":     "8057",
	}

	// Ejecutar cada main.go en paralelo
	for dir, port := range services {
		go func(d, p string) {
			err := runMain(d, p)
			if err != nil {
				log.Printf("Failed to run main.go in %s on port %s: %v", d, p, err)
			}
		}(dir, port)
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
