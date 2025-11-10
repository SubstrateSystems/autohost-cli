package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AgentConfig struct {
	ID        string `json:"id"`
	TokenUser string `json:"token"`
}

func Save(cfg AgentConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".autohost", "config", "agent.json")
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	// 0600 para proteger secretos
	return os.WriteFile(path, b, 0o600)
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
