package expose

import (
	"autohost-cli/internal/adapters/caddy"
	"autohost-cli/internal/adapters/infra"
	tailscale "autohost-cli/internal/adapters/tilscale"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeSetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configura la exposición de servicios",
		Run: func(cmd *cobra.Command, args []string) {
			caddy.InstallCaddy()
			caddy.CreateCaddyfile()

			// 1) IP tailscale local (este host será el nameserver)
			tailIP, err := tailscale.TailscaleIP()
			if err != nil || tailIP == "" {
				fmt.Println("❌ No se pudo obtener IP de tailscale (¿logueado?):", err)
				return
			}

			fmt.Printf("🛰️  IP tailnet local: %s\n", tailIP)

			// 2) dividir host y apex (zona)
			host, zone := splitHostZone("test-server")
			if zone == "" || host == "" {
				fmt.Printf("subdomain inválido: %s (esperado: host.zona, p.ej. app.maza-server)\n", "test-server")
				return
			}

			fmt.Printf("🌐 Zona: %s | Host: %s\n", zone, host)

			// 3) CoreDNS (Docker): asegurar contenedor y Corefile base
			corefilePath, err := infra.InstallAndRunCoreDNSWithDocker(zone, "test-server", tailIP)
			if err != nil {
				fmt.Println("❌ No se pudo instalar CoreDNS con Docker:", err)
				return
			}

			fmt.Println("🧩 CoreDNS (Docker) listo. Corefile:", corefilePath)
		},
	}
	return cmd
}

func splitHostZone(fqdn string) (host, zone string) {
	s := strings.TrimSpace(fqdn)
	if s == "" {
		return "", ""
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return "", ""
	}
	host = parts[0]
	zone = strings.Join(parts[1:], ".")
	return
}
