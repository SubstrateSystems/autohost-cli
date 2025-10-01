package app

import (
	"autohost-cli/internal/adapters/docker"
	"autohost-cli/internal/app"
	"autohost-cli/internal/domain"
	"autohost-cli/internal/platform/di"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func appInstallCmd(deps di.Deps) *cobra.Command {
	var listOnly bool

	var svc = &app.AppService{
		Docker: docker.New(),
	}

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
			if err := svc.InstallApp(ctx, deps); err != nil {
				return fmt.Errorf("error instalando app: %w", err)
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
