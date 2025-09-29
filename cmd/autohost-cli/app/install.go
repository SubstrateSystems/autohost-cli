package app

import (
	appKit "autohost-cli/cmd/autohost-cli/app/appkit"
	"autohost-cli/internal/domain"
	"autohost-cli/internal/platform/di"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func appInstallCmd(deps di.Deps) *cobra.Command {
	var listOnly bool

	cmd := &cobra.Command{
		Use:   "install [nombre]",
		Short: "Instala una aplicación (por ejemplo: nextcloud, bookstack, etc.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if listOnly {
				apps, err := deps.Repos.Catalog.ListApps(ctx)
				if err != nil {
					return err
				}
				printCatalogTable(apps)
				return nil
			}

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
				return fmt.Errorf("error al instalar %s: %w", cfg.Name, err)
			}

			startApp := utils.AskInput(reader, fmt.Sprintf("¿Deseas iniciar %s ahora? [Y/N]: ", cfg.Name), "Y")

			appModel := domain.InstalledApp{
				Name:         cfg.Name,
				CatalogAppID: cfg.Template,
			}

			if err := deps.Repos.Installed.Add(ctx, appModel); err != nil {
				return fmt.Errorf("error al registrar la aplicación instalada: %w", err)
			}

			if strings.EqualFold(startApp, "Y") {
				if err := appKit.StartApp(cfg.Name); err != nil {
					return fmt.Errorf("error al iniciar %s: %w", cfg.Name, err)
				}
				fmt.Printf("🚀 La aplicación %s ha sido iniciada en http://localhost:%s\n", cfg.Name, cfg.Port)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&listOnly, "list", "l", false, "Mostrar catálogo e ignorar instalación")
	return cmd
}

func printCatalogTable(apps []domain.CatalogApp) {
	if len(apps) == 0 {
		fmt.Println("No hay apps disponibles en el catálogo.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tDESCRIPTION")

	for _, a := range apps {
		desc := strings.ReplaceAll(a.Description, "\n", " ")
		fmt.Fprintf(w, "%s\t%s\n", a.Name, desc)
	}

	_ = w.Flush()
}
