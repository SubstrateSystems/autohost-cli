package tailscale

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func tailscaleStatus() error {
	fmt.Println("📊 Estado de Tailscale:")
	statusCmd := exec.Command("sudo", "tailscale", "status")
	statusCmd.Stdout = os.Stdout
	statusCmd.Stderr = os.Stderr
	return statusCmd.Run()
}

// statusJSON es solo la parte que necesitamos del JSON de "tailscale status --json"
type statusJSON struct {
	Self struct {
		HostName string `json:"HostName"`
	} `json:"Self"`
}

// GetMachineName obtiene el nombre de la máquina en el tailnet (ej: "maza-server")
func GetMachineName() (string, error) {
	cmd := exec.Command("tailscale", "status", "--json")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error ejecutando tailscale status: %w", err)
	}

	var st statusJSON
	if err := json.Unmarshal(out, &st); err != nil {
		return "", fmt.Errorf("error parseando JSON: %w", err)
	}

	// El HostName que aparece en tu tailnet
	name := strings.TrimSpace(st.Self.HostName)
	if name == "" {
		return "", fmt.Errorf("no se encontró HostName en status")
	}
	return name, nil
}
