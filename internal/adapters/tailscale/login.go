package tailscale

import (
	"fmt"
	"os"
	"os/exec"
)

func LoginTailscale() error {
	fmt.Println("🔐 Autenticando con Tailscale...")

	loginCmd := exec.Command("sudo", "tailscale", "up")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr

	if err := loginCmd.Run(); err != nil {
		fmt.Println("❌ Error al conectar con Tailscale:", err)
		return err
	}

	fmt.Println("✅ Conectado a Tailscale.")
	return nil
}
