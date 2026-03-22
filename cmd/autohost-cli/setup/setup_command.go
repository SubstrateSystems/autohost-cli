package setup

import (
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func SetupCmd(svc *app.SetupService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup <provider>",
		Short: "Instala y configura providers (docker, tailscale)",
		Long: `Instala y configura los providers necesarios para el self-hosting.

Providers disponibles:
  docker     Instala Docker y crea la red compartida
  tailscale  Instala Tailscale y conecta al tailnet`,
		Example: `  autohost setup docker
  autohost setup tailscale`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(dockerSubCmd(svc))
	cmd.AddCommand(tailscaleSubCmd(svc))

	return cmd
}

func dockerSubCmd(svc *app.SetupService) *cobra.Command {
	return &cobra.Command{
		Use:     "docker",
		Short:   "Instala Docker y configura la red compartida",
		Example: `  autohost setup docker`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🐳 Configurando Docker...")
			if err := svc.SetupDocker(); err != nil {
				return err
			}
			return nil
		},
	}
}

func tailscaleSubCmd(svc *app.SetupService) *cobra.Command {
	return &cobra.Command{
		Use:     "tailscale",
		Short:   "Instala Tailscale y conecta al tailnet",
		Example: `  autohost setup tailscale`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🔐 Configurando Tailscale...")
			if err := svc.SetupTailscale(); err != nil {
				return err
			}
			return nil
		},
	}
}
