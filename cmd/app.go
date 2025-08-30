package cmd

import (
	"autohost-cli/internal/config"
	"autohost-cli/internal/helpers/app"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func askAppConfig(reader *bufio.Reader) config.AppConfig {
	defaultAppName := "appdemo"
	name := utils.AskInput(reader, "📝 Nombre de la aplicación", defaultAppName)

	defaultTemplate := "bookstack"
	var template string
	for {
		template = utils.AskInput(reader, "📦 Tipo de template (bookstack, nextcloud, redis, mysql)", defaultTemplate)
		if config.ValidTemplates[template] {
			break
		}
		fmt.Println("❌ Template no válido. Opciones: bookstack, nextcloud, redis, mysql.")
	}

	port := utils.AskInput(reader, "🔌 Puerto del host a utilizar", config.TemplatePorts[template])

	var mysqlCfg *config.MySQLConfig
	if template == "nextcloud" || template == "bookstack" {
		mysqlCfg = askMySQLConfig(reader, name)
	}

	return config.AppConfig{
		Name:     name,
		Template: template,
		Port:     port,
		MySQL:    mysqlCfg,
	}
}

func askMySQLConfig(reader *bufio.Reader, name string) *config.MySQLConfig {
	fmt.Println("\n⚙️  Configuración de MySQL:")
	user := utils.AskInput(reader, "MySQL usuario", "ah_user")
	pass := utils.AskInput(reader, "MySQL contraseña", "autohost")
	rootPass := utils.AskInput(reader, "MySQL contraseña root", "autohost")
	db := utils.AskInput(reader, "MySQL base", name)
	port := utils.AskInput(reader, "MySQL puerto", "3306")
	return &config.MySQLConfig{
		User:         user,
		Password:     pass,
		RootPassword: rootPass,
		Database:     db,
		Port:         port,
	}
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Gestión de aplicaciones autohospedadas",
}

var appInstallCmd = &cobra.Command{
	Use:   "install [nombre]",
	Short: "Instala una aplicación (por ejemplo: nextcloud, bookstack, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		cfg := askAppConfig(reader)

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
