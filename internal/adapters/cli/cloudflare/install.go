package cloudflare

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func cloudflareInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Instala Cloudflare Tunnel (cloudflared)",
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.IsInitialized() {
				fmt.Println("⚠️ AutoHost no está inicializado. Ejecuta `autohost init` primero.")
				return
			}

			fmt.Println("🌐 Instalando Cloudflare Tunnel (cloudflared)...")
			installCmd := exec.Command("sh", "-c", `
			curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o cloudflared &&
			chmod +x cloudflared &&
			sudo mv cloudflared /usr/local/bin/
		`)
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr

			err := installCmd.Run()
			if err != nil {
				fmt.Println("❌ Error al instalar cloudflared:", err)
			} else {
				fmt.Println("✅ Cloudflare Tunnel instalado con éxito.")
			}
		},
	}
}
