package app

import (
	appKit "autohost-cli/cmd/autohost-cli/app/appkit"
	"autohost-cli/utils"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func appStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [nombre]",
		Short: "Inicia una aplicación",
		Args:  cobra.ExactArgs(1),
		Run: utils.WithAppName(func(ctx context.Context, appName string) {
			err := appKit.StartApp(appName)
			if err != nil {
				fmt.Printf("❌ No se pudo iniciar %s: %v\n", appName, err)
			} else {
				fmt.Printf("🚀 %s iniciada correctamente.\n", appName)
			}
		}),
	}

}
