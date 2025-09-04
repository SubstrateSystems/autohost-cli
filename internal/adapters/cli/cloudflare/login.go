package cloudflare

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func cloudflareLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Inicia sesión con tu cuenta de Cloudflare",
		Run: func(cmd *cobra.Command, args []string) {
			if !utils.IsInitialized() {
				fmt.Println("⚠️ Ejecuta `autohost init` primero.")
				return
			}

			fmt.Println("🔐 Ejecutando 'cloudflared tunnel login'...")
			loginCmd := exec.Command("cloudflared", "tunnel", "login")
			loginCmd.Stdout = os.Stdout
			loginCmd.Stderr = os.Stderr
			err := loginCmd.Run()
			if err != nil {
				fmt.Println("❌ Error al iniciar sesión:", err)
			} else {
				fmt.Println("✅ Sesión iniciada correctamente.")
			}
		},
	}
}
