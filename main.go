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

	"prediapp.local/db"
)

var cmds []*exec.Cmd

func runMain(dir string, envVars map[string]string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = dir

	env := os.Environ()
	for k, v := range envVars {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting %s: %w", dir, err)
	}
	cmds = append(cmds, cmd)
	return nil
}

func cleanup() {
	for _, c := range cmds {
		if c != nil && c.Process != nil {
			c.Process.Kill()
		}
	}
	log.Println("All services stopped.")
}

func main() {
	// Verificar si se pasó el argumento "migrate"
	args := os.Args
	migrateOnly := len(args) > 1 && args[1] == "migrate"

	// 1) Carga de .env
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	godotenv.Load(filepath.Join(cwd, ".env"))
	env := os.Getenv("ENV")
	if env == "" {
		env = "stage"
	}
	godotenv.Load(filepath.Join(cwd, ".env."+env))

	// 2) Conectar a la base de datos
	if err := db.Init(); err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	log.Println("DB connected ✔")

	// 3) Ejecutar migraciones
	if err := db.Migrate(); err != nil {
		log.Fatalf("DB migration failed: %v", err)
	}
	log.Println("DB migrations applied ✔")

	// Si es solo migración, salir después de migrar
	if migrateOnly {
		if err := db.DisconnectDB(); err != nil {
			log.Fatalf("DB disconnect failed: %v", err)
		}
		log.Println("DB disconnected ✔")
		return
	}

	// 4) Preparar URLs de servicios
	services := []struct{ dir, url string }{
		{"./users/cmd", os.Getenv("USERS_SERVICE_URL")},
		{"./sessions/cmd", os.Getenv("SESSIONS_SERVICE_URL")},
		{"./drivers/cmd", os.Getenv("DRIVERS_SERVICE_URL")},
		{"./results/cmd", os.Getenv("RESULTS_SERVICE_URL")},
		{"./prodes/cmd", os.Getenv("PRODES_SERVICE_URL")},
		{"./groups/cmd", os.Getenv("GROUPS_SERVICE_URL")},
		{"./posts/cmd", os.Getenv("POSTS_SERVICE_URL")},
	}

	// 5) Recopilar env vars base
	baseEnv := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			baseEnv[parts[0]] = parts[1]
		}
	}

	// 6) Arrancar los microservicios
	for i, svc := range services {
		parts := strings.Split(svc.url, ":")
		if len(parts) < 3 {
			log.Printf("invalid URL for %s: %s", svc.dir, svc.url)
			continue
		}
		svcEnv := map[string]string{}
		for k, v := range baseEnv {
			svcEnv[k] = v
		}
		svcEnv["PORT"] = parts[2]

		log.Printf("Starting %s on %s...", svc.dir, parts[2])
		if err := runMain(svc.dir, svcEnv); err != nil {
			log.Printf("error: %v", err)
		}
		if i < len(services)-1 {
			time.Sleep(time.Second)
		}
	}

	// 7) Esperar Ctrl+C y limpiar
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Stopping services…")
	cleanup()

	// 8) Desconectar DB
	if err := db.DisconnectDB(); err != nil {
		log.Fatalf("DB disconnect failed: %v", err)
	}
	log.Println("DB disconnected ✔")
}
