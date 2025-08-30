package config

type AppConfig struct {
	Name     string
	Template string
	Port     string
	MySQL    *MySQLConfig
}

type MySQLConfig struct {
	RootPassword string
	User         string
	Password     string
	Database     string
	Port         string
}

var ValidTemplates = map[string]bool{"bookstack": true, "nextcloud": true, "redis": true, "mysql": true}
var TemplatePorts = map[string]string{
	"bookstack": "6875",
	"nextcloud": "8081",
	"redis":     "6379",
	"mysql":     "3306",
}
