// cmd/app/ls.go
package app

import (
	"autohost-cli/internal/di"
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func appLsCmd(deps di.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "Lista las apps instaladas",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			apps, err := deps.Services.App.ListInstalled(ctx)
			if err != nil {
				return fmt.Errorf("no se pudo obtener lista de apps: %w", err)
			}

			if len(apps) == 0 {
				fmt.Println("No hay aplicaciones instaladas aún.")
				return nil
			}

			// salida tabulada
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tINSTALADA")
			for _, a := range apps {
				fmt.Fprintf(
					w,
					"%d\t%s\t%s\n",
					a.ID, a.Name, a.CreatedAt.Format(time.RFC3339),
				)
			}
			_ = w.Flush()
			return nil
		},
	}
}
