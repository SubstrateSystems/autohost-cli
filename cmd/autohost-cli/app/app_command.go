package app

import (
	"autohost-cli/internal/platform/di"

	"github.com/spf13/cobra"
)

func AppCmd(deps di.Deps) *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Application management",
	}

	appCmd.AddCommand(appLsCmd(deps))
	appCmd.AddCommand(appRemoveCmd(deps))
	appCmd.AddCommand(appStartCmd())
	appCmd.AddCommand(appStatusCmd())
	appCmd.AddCommand(appStopCmd())

	return appCmd
}
