package appKit

import (
	"autohost-cli/utils"
	"path/filepath"
)

func RemoveApp(app string) error {
	ymlPath := filepath.Join(utils.GetSubdir("apps"), app, "docker-compose.yml")
	utils.ExecWithDir(filepath.Dir(ymlPath), "docker", "compose", "-f", ymlPath, "down")
	return utils.Exec("rm", "-rf", filepath.Join(utils.GetSubdir("apps"), app))
}
