package appkit

import (
	"autohost-cli/utils"
	"os/exec"
	"path/filepath"
	"strings"
)

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
