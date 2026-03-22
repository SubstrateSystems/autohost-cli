package expose

import (
	"autohost-cli/internal/app"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func exposeSetupCmd(svc *app.ExposeService) *cobra.Command {
	var mode, domain string
	var yes bool

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configura la exposición de servicios",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			mode = strings.ToLower(strings.TrimSpace(mode))
			fmt.Println("Modo seleccionado:", mode)
			switch mode {
			case "private", "public":
			default:
				return fmt.Errorf("modo inválido: %q (usa: private|public)", mode)
			}
			_ = yes
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			switch mode {
			case "private":
				fmt.Println("🔒 Modo PRIVATE: tailnet + DNS interno")
				return svc.SetupPrivate(ctx)
			case "public":
				fmt.Printf("🌍 Modo PUBLIC con dominio: %s\n", domain)
				return svc.SetupPublic()
			default:
				return fmt.Errorf("modo inválido")
			}
		},
	}

	cmd.Flags().StringVarP(&mode, "mode", "m", "", "private | public")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Confirmar automáticamente")
	_ = cmd.RegisterFlagCompletionFunc("mode", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"private", "public"}, cobra.ShellCompDirectiveNoFileComp
	})
	return cmd
}
