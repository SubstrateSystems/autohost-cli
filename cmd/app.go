package cmd

import (
	"autohost-cli/internal/helpers/app"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Gestión de aplicaciones autohospedadas",
}

var appInstallCmd = &cobra.Command{
	Use:   "install [nombre]",
	Short: "Instala una aplicación (por ejemplo: nextcloud, bookstack, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		cfg := app.AskAppConfig(reader)

		if err := app.InstallApp(cfg); err != nil {
			fmt.Printf("❌ Error al instalar %s: %v\n", cfg.Name, err)
			return
		}

		startApp := utils.AskInput(reader, fmt.Sprintf("¿Deseas iniciar %s ahora? [Y/N]: ", cfg.Name), "Y")

		if startApp == "Y" {
			if err := app.StartApp(cfg.Name); err != nil {
				fmt.Printf("❌ Error al iniciar %s: %v\n", cfg.Name, err)
			} else {
				fmt.Printf("🚀 La aplicación %s ha sido iniciada en http://localhost:%s\n", cfg.Name, cfg.Port)
			}
		}
	},
}
var appStartCmd = &cobra.Command{
	Use:   "start [nombre]",
	Short: "Inicia una aplicación",
	Args:  cobra.ExactArgs(1),
	Run: utils.WithAppName(func(appName string) {
		err := app.StartApp(appName)
		if err != nil {
			fmt.Printf("❌ No se pudo iniciar %s: %v\n", appName, err)
		} else {
			fmt.Printf("🚀 %s iniciada correctamente.\n", appName)
		}
	}),
}

var appStopCmd = &cobra.Command{
	Use:   "stop [nombre]",
	Short: "Detiene una aplicación",
	Args:  cobra.ExactArgs(1),
	Run: utils.WithAppName(func(appName string) {
		err := app.StopApp(appName)

		if err != nil {
			fmt.Printf("❌ No se pudo detener %s: %v\n", appName, err)
		} else {
			fmt.Printf("🛑 %s detenida.\n", appName)
		}
	}),
}

var appRemoveCmd = &cobra.Command{
	Use:   "remove [nombre]",
	Short: "Elimina una aplicación",
	Args:  cobra.ExactArgs(1),
	Run: utils.WithAppName(func(appName string) {
		if utils.Confirm(fmt.Sprintf("¿Estás seguro que quieres eliminar %s? [y/N]: ", appName)) {
			err := app.RemoveApp(appName)
			if err != nil {
				fmt.Printf("❌ No se pudo eliminar %s: %v\n", appName, err)
			} else {
				fmt.Printf("🧹 %s eliminada correctamente.\n", appName)
			}
		}
	}),
}

var appStatusCmd = &cobra.Command{
	Use:   "status [nombre]",
	Short: "Muestra el estado de una aplicación",
	Args:  cobra.ExactArgs(1),
	Run: utils.WithAppName(func(appName string) {
		status, err := app.GetAppStatus(appName)
		if err != nil {
			fmt.Printf("❌ Error al obtener el estado de %s: %v\n", appName, err)
		} else {
			fmt.Printf("📊  Estado de %s: %s\n", appName, status)
		}
	}),
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appInstallCmd)
	appCmd.AddCommand(appStatusCmd)
	appCmd.AddCommand(appRemoveCmd)
	appCmd.AddCommand(appStopCmd)
	appCmd.AddCommand(appStartCmd)
}
