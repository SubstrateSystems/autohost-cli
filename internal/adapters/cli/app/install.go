package app

import (
	appKit "autohost-cli/internal/adapters/cli/app/appkit"
	"autohost-cli/internal/domain"
	"autohost-cli/internal/platform/di"
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
			ctx := cmd.Context()
			reader := bufio.NewReader(os.Stdin)

			ensureUnique := func(name string) error {
				exists, err := deps.Repos.Installed.IsInstalledApp(ctx, name)
				if err != nil {
					return fmt.Errorf("no se pudo validar el nombre: %w", err)
				}
				if exists {
					return fmt.Errorf("el nombre %q ya está en uso", name)
				}
				return nil
			}

			cfg := appKit.AskAppConfig(reader, ensureUnique)

			if err := appKit.InstallApp(cfg); err != nil {
				fmt.Printf("❌ Error al instalar %s: %v\n", cfg.Name, err)
				return
			}

			startApp := utils.AskInput(reader, fmt.Sprintf("¿Deseas iniciar %s ahora? [Y/N]: ", cfg.Name), "Y")

			appModel := domain.InstalledApp{
				Name:         cfg.Name,
				CatalogAppID: cfg.Template,
			}

			err := deps.Repos.Installed.Add(ctx, appModel)
			if err != nil {
				fmt.Println(appModel)
				fmt.Printf("❌ Error al registrar la aplicación instalada: %v\n", err)
				return
			}

			if startApp == "Y" {
				if err := appKit.StartApp(cfg.Name); err != nil {
					fmt.Printf("❌ Error al iniciar %s: %v\n", cfg.Name, err)

				} else {
					fmt.Printf("🚀 La aplicación %s ha sido iniciada en http://localhost:%s\n", cfg.Name, cfg.Port)

				}
			}
		},
	}
}
