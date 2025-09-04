package app

import (
	"autohost-cli/internal/di"

	"github.com/spf13/cobra"
)

func AppCmd(deps di.Deps) *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Gestión de aplicaciones autohospedadas",
	}

	// agrega subcomandos construidos en este mismo paquete
	appCmd.AddCommand(appInstallCmd(deps))
	appCmd.AddCommand(appLsCmd(deps))
	appCmd.AddCommand(appRemoveCmd())
	appCmd.AddCommand(appStartCmd())
	appCmd.AddCommand(appStatusCmd())
	appCmd.AddCommand(appStopCmd())

	return appCmd
}
