package domain

var ValidTemplates = map[string]bool{"bookstack": true, "nextcloud": true, "redis": true, "mysql": true}
var TemplatePorts = map[string]string{
	"bookstack": "6875",
	"nextcloud": "8081",
	"redis":     "6379",
	"mysql":     "3306",
	"postgres":  "5432",
}
