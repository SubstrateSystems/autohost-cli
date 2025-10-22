package cloudflare

import (
	"fmt"
	"os"
	"os/exec"
)

func Tunnel() error {
	fmt.Println("🌐 Iniciando túnel de Cloudflare...")

	fmt.Printf("⚙️ Creando túnel para %s...\n", "autohost-tunnel")
	// Crear el túnel
	createCmd := exec.Command("cloudflared", "tunnel", "create", "autohost-tunnel")
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr

	if err := createCmd.Run(); err != nil {
		fmt.Println("❌ Error al crear túnel:", err)
		return err
	}

	return nil
}
