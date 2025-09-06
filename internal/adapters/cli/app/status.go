package app

import (
	appKit "autohost-cli/internal/adapters/cli/app/appkit"
	"autohost-cli/utils"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func appStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [nombre]",
		Short: "Muestra el estado de una aplicación",
		Args:  cobra.ExactArgs(1),
		Run: utils.WithAppName(func(ctx context.Context, appName string) {
			status, err := appKit.GetAppStatus(appName)
			if err != nil {
				fmt.Printf("❌ Error al obtener el estado de %s: %v\n", appName, err)
			} else {
				fmt.Printf("📊  Estado de %s: %s\n", appName, status)
			}
		}),
	}
}
