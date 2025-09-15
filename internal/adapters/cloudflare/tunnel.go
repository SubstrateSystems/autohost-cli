package cloudflare

import (
	"fmt"
	"os"
	"os/exec"
)

func TunnelCloudflare(domain string) {
	fmt.Println("🌐 Iniciando túnel de Cloudflare...")

	fmt.Printf("⚙️ Creando túnel para %s...\n", domain)
	// Crear el túnel
	createCmd := exec.Command("cloudflared", "tunnel", "create", "autohost-tunnel")
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	err := createCmd.Run()
	if err != nil {
		fmt.Println("❌ Error al crear túnel:", err)
		return
	}

}
