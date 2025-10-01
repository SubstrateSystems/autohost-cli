package docker

import (
	"autohost-cli/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func StopApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")

	return utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "stop")
}

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

func RemoveApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")
	utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "down")
	return utils.Exec("rm", "-rf", filepath.Join(utils.GetSubdir("apps"), app))
}

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
