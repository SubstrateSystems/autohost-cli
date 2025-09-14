package cloudflare

import (
	"autohost-cli/internal/platform/config"
	"fmt"
	"os"
	"os/exec"
)

func InstallCloudflare() {
	fmt.Println("🌐 Instalando Cloudflare Tunnel (cloudflared)...")

	downloadURL := config.MustString("url.toml", "cloudflared", "download_url")
	installPath := config.MustString("url.toml", "cloudflared", "install_path")

	cmd := exec.Command("sh", "-c", fmt.Sprintf(`
		curl -L %s -o cloudflared &&
		chmod +x cloudflared &&
		sudo mv cloudflared %s
	`, downloadURL, installPath))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Error al instalar cloudflared:", err)
	} else {
		fmt.Println("✅ Cloudflare Tunnel instalado con éxito.")
	}
}
