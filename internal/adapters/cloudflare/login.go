package cloudflare

import (
	"fmt"
	"os"
	"os/exec"
)

func LoginCloudflare() error {
	fmt.Println("🔑 Iniciando sesión en Cloudflare...")

	loginCmd := exec.Command("cloudflared", "tunnel", "login")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr

	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("error logging in to Cloudflare: %w", err)
	}
	return nil
}
