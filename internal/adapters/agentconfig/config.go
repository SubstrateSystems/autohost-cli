// Package agentconfig manages the agent configuration file at /etc/autohost/config.yaml.
package agentconfig

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/yaml.v3"
)

// AgentConfig holds the agent's runtime configuration.
type AgentConfig struct {
	ApiToken string
	ApiURL   string
}

const configPath = "/etc/autohost/config.yaml"

// Load reads the agent config from /etc/autohost/config.yaml.
func Load() (*AgentConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer %s: %w", configPath, err)
	}
	var raw struct {
		APIURL     string `yaml:"api_url"`
		AgentToken string `yaml:"agent_token"`
	}
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return nil, fmt.Errorf("config inválido: %w", err)
	}
	return &AgentConfig{ApiURL: raw.APIURL, ApiToken: raw.AgentToken}, nil
}

// Save updates api_url and agent_token in /etc/autohost/config.yaml, preserving
// all other content. Uses sudo when the process is not running as root.
func Save(cfg AgentConfig) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("archivo de configuración no existe: %s. Por favor, ejecuta 'autohost agent install' primero", configPath)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error leyendo configuración: %w", err)
	}

	updated := string(content)
	updated = regexp.MustCompile(`(?m)^agent_token:.*$`).
		ReplaceAllString(updated, fmt.Sprintf(`agent_token: "%s"`, cfg.ApiToken))
	if cfg.ApiURL != "" {
		updated = regexp.MustCompile(`(?m)^api_url:.*$`).
			ReplaceAllString(updated, fmt.Sprintf(`api_url: "%s"`, cfg.ApiURL))
	}

	if os.Geteuid() == 0 {
		return os.WriteFile(configPath, []byte(updated), 0600)
	}

	// Write to a temp file first, then copy with sudo.
	tmp, err := os.CreateTemp("", "autohost-config-*.yaml")
	if err != nil {
		return fmt.Errorf("error creando archivo temporal: %w", err)
	}
	defer os.Remove(tmp.Name())

	if err := os.WriteFile(tmp.Name(), []byte(updated), 0600); err != nil {
		return fmt.Errorf("error escribiendo archivo temporal: %w", err)
	}
	tmp.Close()

	cmd := exec.Command("sudo", "cp", tmp.Name(), configPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error copiando archivo con sudo: %w", err)
	}
	return nil
}
