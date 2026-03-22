package install

import (
	"autohost-cli/internal/app"
	"autohost-cli/internal/domain"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func InstallCmd(svc *app.AppService) *cobra.Command {
	var listOnly bool

	cmd := &cobra.Command{
		Use:   "install [name]",
		Short: "Installs an application (e.g., nextcloud, bookstack, etc.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if listOnly {
				apps, err := svc.ListCatalog(ctx)
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
				apps, _ := svc.ListCatalog(ctx)
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

func printCatalogTable(apps []domain.CatalogApp) {
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
