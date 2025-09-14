package cloudflare

import (
	"fmt"
	"os"
	"os/exec"
)

func LoginCloudflare() {
	fmt.Println("🔑 Iniciando sesión en Cloudflare...")

	loginCmd := exec.Command("cloudflared", "tunnel", "login")
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	err := loginCmd.Run()
	if err != nil {
		fmt.Println("❌ Error al iniciar sesión:", err)
	} else {
		fmt.Println("✅ Sesión iniciada correctamente.")
	}
}
