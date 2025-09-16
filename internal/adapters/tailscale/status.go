package tailscale

import (
	"fmt"
	"os"
	"os/exec"
)

func tailscaleStatus() error {
	fmt.Println("📊 Estado de Tailscale:")
	statusCmd := exec.Command("sudo", "tailscale", "status")
	statusCmd.Stdout = os.Stdout
	statusCmd.Stderr = os.Stderr
	return statusCmd.Run()
}
