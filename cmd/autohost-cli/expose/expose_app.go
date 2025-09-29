package expose

import (
	"autohost-cli/internal/adapters/caddy"
	coredns "autohost-cli/internal/adapters/coreDNS"
	"autohost-cli/internal/adapters/tailscale"
	"autohost-cli/internal/app"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeAppCmd() *cobra.Command {
	var (
		exposeType string
		subdomain  string
		nameApp    string
		port       int
	)

	var svc = &app.ExposeService{
		Caddy:     caddy.New(),
		CoreDNS:   coredns.New(),
		Tailscale: tailscale.New(),
	}

	cmd := &cobra.Command{
		Use:   "app",
		Short: "Configura la exposición de servicios",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			exposeType = strings.ToLower(strings.TrimSpace(exposeType))
			switch exposeType {
			case "private", "public":

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
				if err := svc.ExposeApp(ctx, subdomain, nameApp, port); err != nil {
					return fmt.Errorf("error exponiendo app: %w", err)
				}
			}
			fmt.Printf("Exponiendo %q en %q a través de %q\n", subdomain, nameApp, exposeType)
			return nil
		},
	}
	cmd.Flags().StringVar(&exposeType, "type", "", "Tipo de exposición: private o public")
	cmd.Flags().StringVar(&subdomain, "subdomain", "", "Subdominio a exponer")
	cmd.Flags().StringVar(&nameApp, "app", "", "Nombre de la aplicación")
	cmd.Flags().IntVar(&port, "port", 8080, "Puerto de la aplicación")

	return cmd
}
