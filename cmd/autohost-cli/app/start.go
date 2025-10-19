package app

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func appStartCmd() *cobra.Command {

	var svc = &app.AppService{
		Docker: docker.New(),
	}

	cmd := &cobra.Command{
		Use:   "start [nombre]",
		Short: "Inicia una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]

			if err := svc.StartApp(appName); err != nil {
				fmt.Printf("❌ No se pudo iniciar %s: %v\n", appName, err)
			} else {
				fmt.Printf("🚀 %s iniciada correctamente.\n", appName)
			}
			return nil
		},
	}

	return cmd
}
