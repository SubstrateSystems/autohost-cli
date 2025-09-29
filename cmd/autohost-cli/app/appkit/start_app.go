package appKit

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"path/filepath"
)

// StartApp ejecuta docker compose up -d para una app
func StartApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")

	// Validar si existe el archivo docker-compose.yml
	if _, err := os.Stat(ymlPath); os.IsNotExist(err) {
		return fmt.Errorf("el archivo de configuración no existe: %s", ymlPath)
	}

	fmt.Printf("🔄 Levantando aplicación '%s'...\n", app)

	// Usar Exec con working dir del compose
	return utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "up", "-d")
}
