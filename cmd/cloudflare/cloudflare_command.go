package cloudflare

import "github.com/spf13/cobra"

func CloudflareCmd() *cobra.Command {
	cloudflareCmd := &cobra.Command{
		Use:   "cloudflare",
		Short: "Comandos para instalar y configurar Cloudflare Tunnel",
	}

	cloudflareCmd.AddCommand(cloudflareInstallCmd())
	cloudflareCmd.AddCommand(cloudflareLoginCmd())
	cloudflareCmd.AddCommand(cloudflareTunnelCmd())

	return cloudflareCmd

}
