package install

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

func InstallCmd(deps di.Deps) *cobra.Command {
	var listOnly bool

	var svc = &app.AppService{
		Docker:    docker.New(),
		Installed: deps.Repos.Installed,
	}

	cmd := &cobra.Command{
		Use:   "install [name]",
		Short: "Installs an application (e.g., nextcloud, bookstack, etc.)",
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
			var name string
			if len(args) > 0 {
				name = args[0]
			}

			if name == "" {
				fmt.Println("Selecciona una aplicación para instalar:")
				apps, _ := deps.Repos.Catalog.ListApps(ctx)
				for i, a := range apps {
					fmt.Printf("[%d] %s\n", i+1, a.Name)
				}
				fmt.Print("enter a number: ")
				var choice int
				fmt.Scanln(&choice)
				if choice < 1 || choice > len(apps) {
					return fmt.Errorf("invalid selection")
				}
				name = apps[choice-1].Name
			}

			if err := svc.InstallApp(ctx, name); err != nil {
				return fmt.Errorf("error installing app: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&listOnly, "list", "l", false, "Show catalog and skip installation")
	return cmd
}

func printCatalogTable(apps []domain.CatalogItem) {
	if len(apps) == 0 {
		fmt.Println("No apps available in the catalog.")
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
