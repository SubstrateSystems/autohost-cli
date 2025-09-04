package appkit

import (
	"autohost-cli/utils"
	"path/filepath"
)

func StopApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")

	return utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "stop")
}
