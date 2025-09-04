package caddy

import "github.com/spf13/cobra"

func CaddyCmd() *cobra.Command {
	caddyCmd := &cobra.Command{
		Use:   "caddy",
		Short: "Comandos para instalar y administrar el servidor Caddy",
	}

	caddyCmd.AddCommand(caddyInstallCmd())
	caddyCmd.AddCommand(caddyStartCmd())
	caddyCmd.AddCommand(caddyAddServiceCmd())

	return caddyCmd
}
