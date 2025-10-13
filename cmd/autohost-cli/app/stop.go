package app

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func appStopCmd() *cobra.Command {

	var svc = &app.AppService{
		Docker: docker.New(),
	}

	cmd := &cobra.Command{
		Use:   "stop [nombre]",
		Short: "Detiene una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			err := svc.StopApp(appName)

			if err != nil {
				fmt.Printf("❌ No se pudo detener %s: %v\n", appName, err)
			} else {
				fmt.Printf("🛑 %s detenida.\n", appName)
			}
			return nil
		},
	}
	return cmd

}
