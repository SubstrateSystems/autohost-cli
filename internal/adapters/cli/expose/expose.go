package expose

import (
	coredns "autohost-cli/internal/adapters/coreDNS"
	"autohost-cli/internal/adapters/tailscale"
	"autohost-cli/internal/adapters/terraform"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeAppCmd() *cobra.Command {
	var (
		exposeType string
		subdomain  string
		nameApp    string
	)

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Configura la exposición de servicios",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			exposeType = strings.ToLower(strings.TrimSpace(exposeType))
			switch exposeType {
			case "private", "public":
				// ok
			default:
				return fmt.Errorf("tipo inválido: %q (usa: private|public)", exposeType)
			}

			if subdomain == "" {
				return fmt.Errorf("subdominio no puede estar vacío")
			}
			if nameApp == "" {
				return fmt.Errorf("nombre de la app no puede estar vacío")
			}
			return nil

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			switch exposeType {
			case "public":

				fmt.Println("🌐 Exposición vía Cloudflare seleccionada (no implementado aún).")
			case "private":
				// create splitDns in Tailscale
				cfg := terraform.SplitDNSConfig{
					MagicDNS:    true,             // opcional pero útil
					SearchPaths: []string{"test"}, // opcional; permite resolver "maza-server" como "maza-server.test" o "maza-server.test2"
					SplitNameservers: map[string][]string{
						"test": {"100.112.92.90"},
					},
				}
				if err := terraform.ApplySplitDNS(ctx, "default", cfg); err != nil {
					fmt.Printf("⚠️  No se pudo aplicar Split DNS en Tailscale: %v\n", err)
				}

				// update CoreFile and restart
				tailscaleIP, err := tailscale.TailscaleIP()
				if err != nil {
					return fmt.Errorf("no se pudo obtener la IP de Tailscale: %w", err)
				}
				name, err := tailscale.GetMachineName()
				if err != nil {
					return fmt.Errorf("no se pudo obtener el nombre de la máquina en Tailscale: %w", err)
				}
				nameWithSubdomain := fmt.Sprintf("%s.%s", subdomain, name)
				fmt.Printf("🔍 La IP de Tailscale es %q y el nombre de la máquina es %q (usando %q)\n", tailscaleIP, name, nameWithSubdomain)
				coredns.UpdateCorefile(nameWithSubdomain, tailscaleIP)

				// update Caddyfile and restart

			}
			fmt.Printf("Exponiendo %q en %q a través de %q\n", subdomain, nameApp, exposeType)
			return nil
		},
	}
	cmd.Flags().StringVar(&exposeType, "type", "", "Tipo de exposición: private o public")
	cmd.Flags().StringVar(&subdomain, "subdomain", "", "Subdominio a exponer")
	cmd.Flags().StringVar(&nameApp, "app", "", "Nombre de la aplicación")

	return cmd
}
