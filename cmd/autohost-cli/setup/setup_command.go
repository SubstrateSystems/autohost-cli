package setup

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

func SetupCmd() *cobra.Command {

	svc := &app.SetupService{
		Docker: docker.New(),
	}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configura tu servidor para autohospedar servicios",
		Long: `Este comando instala Docker, Caddy, configura dominios,
		y prepara túneles seguros para desplegar tus apps autohospedadas.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n🔧 Iniciando configuración del servidor...")

			if err := svc.Setup(); err != nil {
				return err
			}

			fmt.Println("\n✅ Configuración inicial completa.")
			return nil
		},
	}

	return cmd
}
