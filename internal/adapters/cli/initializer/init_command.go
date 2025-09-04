package initializer

import (
	initializerkit "autohost-cli/internal/adapters/cli/initializer/initializerKit"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func InitCommand() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Inicializa el entorno de AutoHost en ~/.autohost",
		Run: func(cmd *cobra.Command, args []string) {
			err := initializerkit.EnsureAutohostDirs()
			if err != nil {
				fmt.Println("❌ Error al crear estructura de carpetas:", err)
				os.Exit(1)
			}
			fmt.Println("✅ Entorno de AutoHost creado")
		},
	}

	return initCmd
}
