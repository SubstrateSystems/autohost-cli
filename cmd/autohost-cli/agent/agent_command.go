package agent

import (
	"autohost-cli/internal/app"

	"github.com/spf13/cobra"
)

// AgentCmd returns the root `autohost agent` command.
func AgentCmd(svc *app.AgentService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Gestión del agente de AutoHost",
	}
	cmd.AddCommand(newInstallCmd(svc))
	cmd.AddCommand(newUpdateCmd(svc))
	return cmd
}

func newInstallCmd(svc *app.AgentService) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Instala el agente (binario, config base y servicio systemd)",
		Long: `Descarga e instala el binario del agente, una configuración base vacía
y el archivo de servicio systemd. Después de instalar, ejecuta 'autohost up'
para conectar el nodo al cloud.`,
		Example: "  autohost agent install",
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.Install()
		},
	}
}

func newUpdateCmd(svc *app.AgentService) *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Actualiza el agente a la última versión disponible",
		Long: `Descarga e instala la última versión del agente, reemplazando el binario actual,
y reinicia el servicio systemd. La nueva versión se reporta automáticamente al cloud.

Usa la variable de entorno VERSION para fijar una versión específica:
  VERSION=v0.3.0 autohost agent update`,
		Example: "  autohost agent update\n  VERSION=v0.3.0 autohost agent update",
		RunE: func(cmd *cobra.Command, args []string) error {
			return svc.Update()
		},
	}
}
