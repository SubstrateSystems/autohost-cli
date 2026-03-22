package app

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

func AppCmd(svc *app.AppService) *cobra.Command {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Application management",
	}

	appCmd.AddCommand(appLsCmd(svc))
	appCmd.AddCommand(appRemoveCmd(svc))
	appCmd.AddCommand(appStartCmd(svc))
	appCmd.AddCommand(appStatusCmd(svc))
	appCmd.AddCommand(appStopCmd(svc))

	return appCmd
}
