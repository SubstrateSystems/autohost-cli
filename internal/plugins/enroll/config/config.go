package config

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

type AgentConfig struct {
	ApiToken string
	ApiURL   string
}

const configPath = "/etc/autohost/config.yaml"

// Save actualiza el agent_token en /etc/autohost/config.yaml
func Save(cfg AgentConfig) error {
	// Verificar si el archivo existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("archivo de configuración no existe: %s. Por favor, ejecuta 'autohost agent install' primero", configPath)
	}

	// Leer archivo existente
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error leyendo configuración: %w", err)
	}

	// Reemplazar agent_token preservando el resto del archivo
	re := regexp.MustCompile(`(?m)^agent_token:.*$`)
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf(`agent_token: "%s"`, cfg.ApiToken))

	// Si también se proporciona ApiURL, actualizarlo
	if cfg.ApiURL != "" {
		reURL := regexp.MustCompile(`(?m)^api_url:.*$`)
		newContent = reURL.ReplaceAllString(newContent, fmt.Sprintf(`api_url: "%s"`, cfg.ApiURL))
	}

	// Detectar si necesitamos sudo
	needsSudo := os.Geteuid() != 0

	if needsSudo {
		// Escribir a archivo temporal primero
		tmpFile, err := os.CreateTemp("", "autohost-config-*.yaml")
		if err != nil {
			return fmt.Errorf("error creando archivo temporal: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		if err := os.WriteFile(tmpFile.Name(), []byte(newContent), 0600); err != nil {
			return fmt.Errorf("error escribiendo archivo temporal: %w", err)
		}
		tmpFile.Close()

		// Copiar con sudo
		cmd := exec.Command("sudo", "cp", tmpFile.Name(), configPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error copiando archivo con sudo: %w", err)
		}

		// Establecer permisos
		cmd = exec.Command("sudo", "chmod", "600", configPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error estableciendo permisos: %w", err)
		}
	} else {
		// Escribir directamente si somos root
		if err := os.WriteFile(configPath, []byte(newContent), 0600); err != nil {
			return fmt.Errorf("error escribiendo configuración: %w", err)
		}
	}

	return nil
}

// func Load() (*AgentConfig, error) {
// 	path, err := ConfigPath()
// 	if err != nil {
// 		return nil, err
// 	}
// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var cfg AgentConfig
// 	if err := json.Unmarshal(b, &cfg); err != nil {
// 		return nil, err
// 	}
// 	return &cfg, nil
// }
