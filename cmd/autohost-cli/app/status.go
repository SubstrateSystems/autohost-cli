package app

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func appStatusCmd() *cobra.Command {
	var svc = &app.AppService{
		Docker: docker.New(),
	}

	cmd := &cobra.Command{
		Use:   "status [nombre]",
		Short: "Muestra el estado de una aplicación",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			status, err := svc.GetAppStatus(appName)
			if err != nil {
				fmt.Printf("❌ Error al obtener el estado de %s: %v\n", appName, err)
			} else {
				fmt.Printf("📊  Estado de %s: %s\n", appName, status)
			}
			return nil
		},
	}
	return cmd
}
