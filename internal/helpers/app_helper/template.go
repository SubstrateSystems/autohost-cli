package app_helper

import (
	"autohost-cli/internal/config"
	"autohost-cli/utils"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyTemplate copia los archivos de template al destino
func CopyTemplate(appName, destPath string) error {
	srcDir := filepath.Join("templates", appName)
	return copyDir(srcDir, destPath)
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copiar archivo
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}

func AskAppConfig(reader *bufio.Reader) config.AppConfig {
	defaultAppName := "appdemo"
	name := utils.AskInput(reader, "📝 Nombre de la aplicación", defaultAppName)

	defaultTemplate := "bookstack"

	template := utils.AskInput(reader, "📦 Tipo de template (bookstack, nextcloud, redis, mysql)", defaultTemplate)

	port := utils.AskAppPort(reader, "🔌 Puerto del host a utilizar", config.TemplatePorts[template])

	var mysqlCfg *config.MySQLConfig
	if template == "nextcloud" || template == "bookstack" {
		mysqlCfg = AskMySQLConfig(reader, name)
	}

	return config.AppConfig{
		Name:     name,
		Template: template,
		Port:     port,
		MySQL:    mysqlCfg,
	}
}

func AskMySQLConfig(reader *bufio.Reader, name string) *config.MySQLConfig {
	fmt.Println("\n⚙️  Configuración de MySQL:")
	user := utils.AskInput(reader, "MySQL usuario", "ah_user")
	pass := utils.AskInput(reader, "MySQL contraseña", "autohost")
	rootPass := utils.AskInput(reader, "MySQL contraseña root", "autohost")
	db := utils.AskInput(reader, "MySQL base", name)

	port := utils.AskAppPort(reader, "MySQL puerto", "3306")

	return &config.MySQLConfig{
		User:         user,
		Password:     pass,
		RootPassword: rootPass,
		Database:     db,
		Port:         port,
	}
}
