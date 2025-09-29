package cloudflare

import (
	"autohost-cli/internal/platform/config"
	"fmt"
	"os"
	"os/exec"
)

func InstallCloudflare() error {
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
		return fmt.Errorf("error installing cloudflared: %w", err)
	}
	return nil
}
