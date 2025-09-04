package app_helper

import (
	"autohost-cli/assets"
	"autohost-cli/internal/config"
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func InstallApp(app config.AppConfig) error {
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

// StopApp ejecuta docker compose stop para una app
func StopApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")

	return utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "stop")
}

// RemoveApp ejecuta docker compose down para una app
func RemoveApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")
	utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "down")
	return utils.Exec("rm", "-rf", filepath.Join(utils.GetSubdir("apps"), app))

}

// GetAppStatus devuelve si los contenedores están "running", "exited", etc.
func GetAppStatus(app string) (string, error) {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")

	cmd := exec.Command("docker", "compose", "-f", ymlPath, "ps", "--status=running")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if strings.Contains(string(out), "Up") {
		return "en ejecución", nil
	}
	return "detenida", nil
}

func setValues(app config.AppConfig) map[string]string {
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
