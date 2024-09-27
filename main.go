package main

import (
	"log"
	"os"
	"os/exec"
)

func runMain(dir string, port string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = dir

	// Obtener las variables de entorno actuales
	env := os.Environ()

	// Asegurar que GOPATH y GOMODCACHE estén configurados
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = "/Users/agustinglaiel/go" // Asegurar el valor por defecto de GOPATH
	}

	gomodcache := os.Getenv("GOMODCACHE")
	if gomodcache == "" {
		gomodcache = gopath + "/pkg/mod" // Asegurar el valor por defecto de GOMODCACHE
	}

	// Configurar las variables de entorno
	cmd.Env = append(env, "PORT="+port, "GOPATH="+gopath, "GOMODCACHE="+gomodcache)

	// Conectar la salida del proceso a la terminal
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	// Usar Start() para ejecutar en paralelo sin bloquear
	err := cmd.Start()
	if err != nil {
		return err
	}

	// No hacer cmd.Wait() aquí porque queremos que se ejecute en paralelo
	return nil
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
				log.Fatalf("Failed to run main.go in %s on port %s: %v", d, p, err)
			}
		}(dir, port)
	}

	// Mantener el programa corriendo
	select {}
}
