package appKit

import (
	"autohost-cli/assets"
	"autohost-cli/internal/domain"
	"autohost-cli/utils"
	"fmt"
	"os"
	"path/filepath"
)

func InstallApp(app domain.AppConfig) error {
	appDir := filepath.Join(utils.GetSubdir("apps"), app.Name)
	composePath := filepath.Join(appDir, "docker-compose.yml")

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return fmt.Errorf("error creando directorio de destino: %w", err)
	}

	data, err := assets.ReadCompose(app.Template)
	if err != nil {
		return fmt.Errorf("no se encontró plantilla embebida para %s: %w", app.Template, err)
	}
	fmt.Println("📦 Usando plantilla embebida para:", app.Template)

	values := setValues(app)

	fmt.Println(values)

	final := utils.ReplacePlaceholders(string(data), values)

	if err := os.WriteFile(composePath, []byte(final), 0o644); err != nil {
		return fmt.Errorf("error escribiendo archivo docker-compose: %w", err)
	}
	fmt.Println("✅ Archivo docker-compose creado en:", composePath)
	return nil
}

func setValues(app domain.AppConfig) map[string]string {
	values := map[string]string{}

	values["$service-name"] = app.Name // AppName
	values["$port"] = app.Port         // AppPort
	if app.MySQL != nil {
		values["$mysql-user"] = app.MySQL.User                  // UserName MySQL
		values["$mysql-password"] = app.MySQL.Password          // PasswordUser MySQL
		values["$mysql-root-password"] = app.MySQL.RootPassword // RootPassword MySQL
		values["$mysql-database"] = app.MySQL.Database + "-db"  // DatabaseName MySQL
		values["$mysql-port"] = app.MySQL.Port                  // Port MySQL
	}

	if app.Template == "bookstack" {
		values["$app-key"] = utils.GenerateRandomString(64) // AppKey BookStack
	}

	return values
}
