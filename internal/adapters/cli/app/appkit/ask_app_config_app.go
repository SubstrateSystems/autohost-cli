package appKit

import (
	"autohost-cli/internal/domain"
	"autohost-cli/utils"
	"bufio"
	"fmt"
)

func AskAppConfig(reader *bufio.Reader) domain.AppConfig {
	defaultAppName := "appdemo"
	name := utils.AskInput(reader, "📝 Nombre de la aplicación", defaultAppName)

	defaultTemplate := "bookstack"

	template := utils.AskInput(reader, "📦 Tipo de template (bookstack, nextcloud, redis, mysql)", defaultTemplate)

	port := utils.AskAppPort(reader, "🔌 Puerto del host a utilizar", domain.TemplatePorts[template])

	var mysqlCfg *domain.MySQLConfig
	if template == "nextcloud" || template == "bookstack" {
		mysqlCfg = AskMySQLConfig(reader, name)
	}

	return domain.AppConfig{
		Name:     name,
		Template: template,
		Port:     port,
		MySQL:    mysqlCfg,
	}
}

func AskMySQLConfig(reader *bufio.Reader, name string) *domain.MySQLConfig {
	fmt.Println("\n⚙️  Configuración de MySQL:")
	user := utils.AskInput(reader, "MySQL usuario", "ah_user")
	pass := utils.AskInput(reader, "MySQL contraseña", "autohost")
	rootPass := utils.AskInput(reader, "MySQL contraseña root", "autohost")
	db := utils.AskInput(reader, "MySQL base", name)

	port := utils.AskAppPort(reader, "MySQL puerto", "3306")

	return &domain.MySQLConfig{
		User:         user,
		Password:     pass,
		RootPassword: rootPass,
		Database:     db,
		Port:         port,
	}
}
