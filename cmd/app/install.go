package app

import (
	"autohost-cli/internal/di"
	"autohost-cli/internal/helpers/app"
	"autohost-cli/internal/models"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func appInstallCmd(deps di.Deps) *cobra.Command {
	return &cobra.Command{
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

			appModel := models.InstalledApp{
				Name:         cfg.Name,
				CatalogAppID: cfg.Template,
			}
			err := deps.Repos.Installed.Add(cmd.Context(), appModel)
			if err != nil {
				fmt.Println(appModel)
				fmt.Printf("❌ Error al registrar la aplicación instalada: %v\n", err)
				return
			}

			if startApp == "Y" {
				if err := app.StartApp(cfg.Name); err != nil {
					fmt.Printf("❌ Error al iniciar %s: %v\n", cfg.Name, err)

				} else {
					fmt.Printf("🚀 La aplicación %s ha sido iniciada en http://localhost:%s\n", cfg.Name, cfg.Port)

				}
			}
		},
	}
}
