package app

import (
	"autohost-cli/internal/adapters/agentconfig"
	"autohost-cli/internal/adapters/enrollapi"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CCCommandsDir is the directory where custom command scripts are stored.
const CCCommandsDir = "/var/lib/autohost/commands"

// CCService manages custom command scripts and their registration with the AutoHost API.
type CCService struct{}

// CreateCommand creates a bash script template in CCCommandsDir and registers it
// with the AutoHost API so it can be triggered remotely from the dashboard.
func (s *CCService) CreateCommand(name, description string) error {
	scriptPath := filepath.Join(CCCommandsDir, name+".sh")

	if err := s.ensureCommandsDir(); err != nil {
		return fmt.Errorf("error creando directorio de comandos: %w", err)
	}

	if _, err := os.Stat(scriptPath); err == nil {
		return fmt.Errorf("el script '%s.sh' ya existe en %s", name, CCCommandsDir)
	}

	tmpl := fmt.Sprintf(`#!/usr/bin/env bash
# Custom command: %s
# Description: %s
# Created: %s
#
# Edit this script with your custom logic.
# It will be executed by the AutoHost agent when triggered from the dashboard.

set -euo pipefail

echo "Running custom command: %s"

# TODO: Add your commands here
`, name, description, time.Now().Format("2006-01-02 15:04:05"), name)

	if err := s.writeFilePrivileged(scriptPath, []byte(tmpl), 0755); err != nil {
		return fmt.Errorf("error creando script: %w", err)
	}
	fmt.Printf("📄 Script creado: %s\n", scriptPath)

	cfg, err := agentconfig.Load()
	if err != nil {
		return fmt.Errorf("no se pudo leer la configuración: %w", err)
	}
	if cfg.ApiURL == "" {
		return fmt.Errorf("api_url no encontrada en /etc/autohost/config.yaml")
	}
	if cfg.ApiToken == "" {
		return fmt.Errorf("agent_token no encontrado en /etc/autohost/config.yaml")
	}

	fmt.Println("📡 Registrando comando en la API...")
	client := enrollapi.NewAgentClient(cfg.ApiURL, cfg.ApiToken)
	payload := map[string]string{
		"name":        name,
		"type":        "custom",
		"description": description,
		"script_path": scriptPath,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := client.PostJSON(ctx, "/v1/node-commands", payload, nil)
	if err != nil {
		return fmt.Errorf("error registrando comando en la API: %w", err)
	}
	if status >= 300 {
		return fmt.Errorf("la API respondió con código %d", status)
	}

	fmt.Println("✅ Comando registrado exitosamente")
	fmt.Println()
	fmt.Printf("📝 Edita el script con tu lógica:\n   sudo nano %s\n", scriptPath)
	fmt.Println()
	fmt.Println("El comando estará disponible en el dashboard para ejecutarse remotamente.")
	return nil
}

// ListCommands returns the full paths of all .sh scripts in CCCommandsDir.
func (s *CCService) ListCommands() ([]string, error) {
	entries, err := os.ReadDir(CCCommandsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error leyendo directorio: %w", err)
	}
	var paths []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sh") {
			paths = append(paths, filepath.Join(CCCommandsDir, e.Name()))
		}
	}
	return paths, nil
}

func (s *CCService) ensureCommandsDir() error {
	if _, err := os.Stat(CCCommandsDir); err == nil {
		return nil
	}
	if os.Geteuid() != 0 {
		if _, err := exec.LookPath("sudo"); err != nil {
			return fmt.Errorf("se requiere sudo para crear %s. Ejecuta como root o instala sudo", CCCommandsDir)
		}
		cmd := exec.Command("sudo", "mkdir", "-p", CCCommandsDir)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}
	return os.MkdirAll(CCCommandsDir, 0755)
}

func (s *CCService) writeFilePrivileged(path string, data []byte, perm os.FileMode) error {
	if os.Geteuid() == 0 {
		return os.WriteFile(path, data, perm)
	}
	tmp, err := os.CreateTemp("", "autohost-cc-*.sh")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := os.WriteFile(tmp.Name(), data, perm); err != nil {
		return err
	}
	tmp.Close()

	cmd := exec.Command("sudo", "cp", tmp.Name(), path)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	chmod := exec.Command("sudo", "chmod", fmt.Sprintf("%o", perm), path)
	chmod.Stdout, chmod.Stderr = os.Stdout, os.Stderr
	return chmod.Run()
}
