package cc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	enrollcfg "autohost-cli/internal/plugins/enroll/config"
	enrollhttp "autohost-cli/internal/plugins/enroll/http"

	"github.com/spf13/cobra"
)

const (
	// commandsDir is the directory where custom command scripts are stored.
	commandsDir = "/var/lib/autohost/commands"
)

// scriptNameRegexp validates script names: alphanumeric, hyphens, underscores.
var scriptNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func newCreateCmd() *cobra.Command {
	var name, description string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create and register a custom command script",
		Long: `Creates a new bash script in /var/lib/autohost/commands/ and registers it
with the AutoHost API so it can be executed remotely from the dashboard.

The script name should be alphanumeric with hyphens or underscores (no extension needed).
A .sh file will be created with a basic template that you can edit.

Credentials are read from /etc/autohost/config.yaml.`,
		Example: `  autohost cc create --name backup-db --description "Backup the database"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name es obligatorio")
			}
			name = strings.TrimSpace(name)
			if !scriptNameRegexp.MatchString(name) {
				return fmt.Errorf("nombre inválido: solo se permiten letras, números, guiones y guiones bajos")
			}

			cfg, err := enrollcfg.Load()
			if err != nil {
				return fmt.Errorf("no se pudo leer la configuración: %w", err)
			}
			if cfg.ApiURL == "" {
				return fmt.Errorf("api_url no encontrada en /etc/autohost/config.yaml")
			}
			if cfg.ApiToken == "" {
				return fmt.Errorf("agent_token no encontrado en /etc/autohost/config.yaml")
			}

			return createCustomCommand(name, description, cfg.ApiURL, cfg.ApiToken)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Nombre del comando (letras, números, guiones)")
	cmd.Flags().StringVar(&description, "description", "", "Descripción del comando personalizado")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func createCustomCommand(name, description, apiURL, token string) error {
	filename := name + ".sh"
	scriptPath := filepath.Join(commandsDir, filename)

	// 1. Ensure the commands directory exists.
	if err := ensureCommandsDir(); err != nil {
		return fmt.Errorf("error creando directorio de comandos: %w", err)
	}

	// 2. Check if the script already exists.
	if _, err := os.Stat(scriptPath); err == nil {
		return fmt.Errorf("el script '%s' ya existe en %s", filename, commandsDir)
	}

	// 3. Create the script file with a template.
	template := fmt.Sprintf(`#!/usr/bin/env bash
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

	if err := writeFileWithPrivileges(scriptPath, []byte(template), 0755); err != nil {
		return fmt.Errorf("error creando script: %w", err)
	}

	fmt.Printf("📄 Script creado: %s\n", scriptPath)

	// 4. Register the command with the API.
	fmt.Println("📡 Registrando comando en la API...")

	client := enrollhttp.NewAgentClient(apiURL, token)
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
	fmt.Printf("📝 Edita el script con tu lógica:\n")
	fmt.Printf("   sudo nano %s\n", scriptPath)
	fmt.Println()
	fmt.Println("El comando estará disponible en el dashboard para ejecutarse remotamente.")

	return nil
}

// ensureCommandsDir creates the commands directory if it does not exist.
func ensureCommandsDir() error {
	if _, err := os.Stat(commandsDir); err == nil {
		return nil
	}

	needsSudo := os.Geteuid() != 0
	if needsSudo {
		if _, err := exec.LookPath("sudo"); err != nil {
			return fmt.Errorf("se requiere sudo para crear %s. Ejecuta como root o instala sudo", commandsDir)
		}
		cmd := exec.Command("sudo", "mkdir", "-p", commandsDir)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}

	return os.MkdirAll(commandsDir, 0755)
}

// writeFileWithPrivileges writes a file, using sudo if necessary.
func writeFileWithPrivileges(path string, data []byte, perm os.FileMode) error {
	needsSudo := os.Geteuid() != 0

	if needsSudo {
		tmpFile, err := os.CreateTemp("", "autohost-cc-*.sh")
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile.Name())

		if err := os.WriteFile(tmpFile.Name(), data, perm); err != nil {
			return err
		}

		cmd := exec.Command("sudo", "cp", tmpFile.Name(), path)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}

		chmodCmd := exec.Command("sudo", "chmod", fmt.Sprintf("%o", perm), path)
		chmodCmd.Stdout, chmodCmd.Stderr = os.Stdout, os.Stderr
		return chmodCmd.Run()
	}

	return os.WriteFile(path, data, perm)
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List custom command scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := os.ReadDir(commandsDir)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("No hay comandos personalizados. Usa 'autohost-cli cc create <nombre>' para crear uno.")
					return nil
				}
				return fmt.Errorf("error leyendo directorio: %w", err)
			}

			if len(entries) == 0 {
				fmt.Println("No hay comandos personalizados registrados.")
				return nil
			}

			fmt.Println("📋 Comandos personalizados:")
			fmt.Println()
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".sh") {
					name := strings.TrimSuffix(e.Name(), ".sh")
					fmt.Printf("  • %s  (%s)\n", name, filepath.Join(commandsDir, e.Name()))
				}
			}
			return nil
		},
	}
}
