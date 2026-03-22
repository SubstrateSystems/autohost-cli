package cc

import (
	"autohost-cli/internal/app"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// scriptNameRegexp validates script names: alphanumeric, hyphens, underscores.
var scriptNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func newCreateCmd(svc *app.CCService) *cobra.Command {
	var name, description string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create and register a custom command script",
		Long: `Creates a new bash script in /var/lib/autohost/commands/ and registers it
with the AutoHost API so it can be executed remotely from the dashboard.

The script name should be alphanumeric with hyphens or underscores (no extension needed).
Credentials are read from /etc/autohost/config.yaml.`,
		Example: `  autohost cc create --name backup-db --description "Backup the database"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			name = strings.TrimSpace(name)
			if !scriptNameRegexp.MatchString(name) {
				return fmt.Errorf("nombre inválido: solo se permiten letras, números, guiones y guiones bajos")
			}
			return svc.CreateCommand(name, description)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Nombre del comando (letras, números, guiones)")
	cmd.Flags().StringVar(&description, "description", "", "Descripción del comando personalizado")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newListCmd(svc *app.CCService) *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List custom command scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			paths, err := svc.ListCommands()
			if err != nil {
				return err
			}
			if len(paths) == 0 {
				fmt.Println("No hay comandos personalizados. Usa 'autohost cc create --name <nombre>' para crear uno.")
				return nil
			}
			fmt.Println("📋 Comandos personalizados:")
			fmt.Println()
			for _, p := range paths {
				name := strings.TrimSuffix(filepath.Base(p), ".sh")
				fmt.Printf("  • %s  (%s)\n", name, p)
			}
			return nil
		},
	}
}
