package expose

import (
	"fmt"
	"strings"

	"autohost-cli/internal/app"

	acaddy "autohost-cli/internal/adapters/caddy"
	acloud "autohost-cli/internal/adapters/cloudflare"
	ats "autohost-cli/internal/adapters/tailscale"

	"github.com/spf13/cobra"
)

func exposeSetupCmd() *cobra.Command {
	var mode, domain string
	var yes bool

	// Composition root: construir adapters reales
	var svc = &app.ExposeService{
		Caddy:      acaddy.New(),
		Tailscale:  ats.New(),
		Cloudflare: acloud.New(),
	}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configura la exposición de servicios",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			mode = strings.ToLower(strings.TrimSpace(mode))
			switch mode {
			case "private", "public":
			default:
				return fmt.Errorf("modo inválido: %q (usa: private|public)", mode)
			}
			if mode == "public" && strings.TrimSpace(domain) == "" {
				return fmt.Errorf("--domain es requerido en --mode=public")
			}
			_ = yes // si quieres suprimir prompts luego
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch mode {
			case "private":
				fmt.Println("🔒 Modo PRIVATE: tailnet + DNS interno")
				return svc.SetupPrivate(domain)
			case "public":
				fmt.Printf("🌍 Modo PUBLIC con dominio: %s\n", domain)
				return svc.SetupPublic(domain)
			default:
				return fmt.Errorf("modo inválido")
			}
		},
	}

	cmd.Flags().StringVarP(&mode, "mode", "m", "private", "private | public")
	cmd.Flags().StringVar(&domain, "domain", "", "Dominio base para exposición pública")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Confirmar automáticamente")

	_ = cmd.RegisterFlagCompletionFunc("mode", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"private", "public"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
