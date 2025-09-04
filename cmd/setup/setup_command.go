package setup

import (
	"autohost-cli/internal/helpers/caddy_helper"
	"autohost-cli/internal/helpers/cloudflared_helper"
	"autohost-cli/internal/helpers/docker_helper"
	"autohost-cli/internal/helpers/initializer_helper"
	"autohost-cli/internal/helpers/tailscale_helper"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func SetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Configura tu servidor para autohospedar servicios",
		Long: `Este comando instala Docker, Caddy, configura dominios,
		y prepara túneles seguros para desplegar tus apps autohospedadas.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\n🔧 Iniciando configuración del servidor...")

			initializer_helper.EnsureAutohostDirs()

			if !docker_helper.DockerInstalled() {
				if utils.Confirm("⚠️ Docker no está instalado. ¿Deseas instalarlo automáticamente? [y/N]: ") {
					docker_helper.InstallDocker()
					docker_helper.CreateDockerNetwork()
					fmt.Println("✅ Docker instalado correctamente.")
					fmt.Println("✅ Red Docker 'autohost_net' creada.")
					fmt.Println("🔄 Reiniciando sesión para aplicar cambios de grupo...")
				} else {
					fmt.Println("🚫 Instalación cancelada. Instala Docker manualmente y vuelve a ejecutar el setup.")
					return
				}
			} else {
				fmt.Println("✅ Docker ya está instalado.")
			}

			if utils.Confirm("¿Deseas agregar tu usuario al grupo 'docker' para usar Docker sin sudo? [y/N]: ") {
				docker_helper.AddUserToDockerGroup()
			}

			if utils.Confirm("¿Deseas instalar y configurar Caddy como reverse proxy? [y/N]: ") {
				caddy_helper.InstallCaddy()
				caddy_helper.CreateCaddyfile()
			}

			option := utils.AskOption("🔒 ¿Qué tipo de acceso quieres configurar?", []string{"Tailscale (privado)", "Cloudflare Tunnel (público con dominio)"})
			switch option {
			case "Tailscale (privado)":
				tailscale_helper.InstallTailscale()
			case "Cloudflare Tunnel (público con dominio)":
				cloudflared_helper.InstallCloudflared()
				fmt.Print("Introduce el subdominio para el túnel (ej: blog.misitio.com): ")
				reader := bufio.NewReader(os.Stdin)
				domain, _ := reader.ReadString('\n')
				domain = strings.TrimSpace(domain)
				cloudflared_helper.ConfigureCloudflareTunnel(domain)
			}

			fmt.Println("\n✅ Configuración inicial completa.")
		},
	}
}
