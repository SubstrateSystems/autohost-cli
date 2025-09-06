package appKit

import (
	"autohost-cli/internal/domain"
	"autohost-cli/utils"
	"bufio"
	"fmt"
)

func AskAppConfig(reader *bufio.Reader, ensureUnique func(string) error) domain.AppConfig {
	defaultAppName := "appdemo"
	var name string
	for {
		name = utils.AskInput(reader, "📝 Nombre de la aplicación", defaultAppName)
		if err := ensureUnique(name); err != nil {
			fmt.Printf("⚠️ %v\n", err)
			continue
		}
		break
	}
	// name := utils.AskInput(reader, "📝 Nombre de la aplicación", defaultAppName)

	defaultTemplate := "bookstack"

	template := utils.AskInput(reader, "📦 Tipo de template (bookstack, nextcloud, redis, mysql, postgres)", defaultTemplate)

	if template == "mysql" {
		mysqlCfg := AskMySQLConfig(reader, name)
		return domain.AppConfig{
			Name:     name,
			Template: template,
			Port:     mysqlCfg.Port,
			MySQL:    mysqlCfg,
		}
	}

	if template == "postgres" {
		postgresCfg := AskMyPostgresConfig(reader, name)
		return domain.AppConfig{
			Name:     name,
			Template: template,
			Port:     postgresCfg.Port,
			Postgres: postgresCfg,
		}
	}

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

func AskMyPostgresConfig(reader *bufio.Reader, name string) *domain.PostgresConfig {
	fmt.Println("\n⚙️  Configuración de Postgres:")
	user := utils.AskInput(reader, "Postgres usuario", "ah_user")
	pass := utils.AskInput(reader, "Postgres contraseña", "autohost")
	db := utils.AskInput(reader, "Postgres base", name)

	port := utils.AskAppPort(reader, "Postgres puerto", "5432")

	return &domain.PostgresConfig{
		User:     user,
		Password: pass,
		Database: db,
		Port:     port,
	}
}
