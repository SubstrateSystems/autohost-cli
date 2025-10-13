package app

import (
	"autohost-cli/assets"
	"autohost-cli/internal/domain"
	"autohost-cli/internal/platform/di"
	"autohost-cli/internal/ports"
	"autohost-cli/utils"
	"path/filepath"
	"strings"

	"bufio"
	"context"
	"fmt"
	"os"
)

type AppService struct {
	Docker ports.Docker
}

func (s *AppService) InstallApp(ctx context.Context, deps di.Deps) error {
	reader := bufio.NewReader(os.Stdin)

	ensureUnique := func(name string) error {
		exists, err := deps.Repos.Installed.IsInstalledApp(ctx, name)
		if err != nil {
			return fmt.Errorf("no se pudo validar el nombre: %w", err)
		}
		if exists {
			return fmt.Errorf("el nombre %q ya está en uso", name)
		}
		return nil
	}

	cfg := askAppConfig(reader, ensureUnique)

	if err := install(cfg); err != nil {
		return fmt.Errorf("error al instalar %s: %w", cfg.Name, err)
	}

	startApp := utils.AskInput(reader, fmt.Sprintf("¿Deseas iniciar %s ahora? [Y/N]: ", cfg.Name), "Y")

	appModel := domain.InstalledApp{
		Name:         cfg.Name,
		CatalogAppID: cfg.Template,
	}

	if err := deps.Repos.Installed.Add(ctx, appModel); err != nil {
		return fmt.Errorf("error al registrar la aplicación instalada: %w", err)
	}

	if strings.EqualFold(startApp, "Y") {
		if err := s.Docker.StartApp(cfg.Name); err != nil {
			return fmt.Errorf("error al iniciar %s: %w", cfg.Name, err)
		}
		fmt.Printf("🚀 La aplicación %s ha sido iniciada en http://localhost:%s\n", cfg.Name, cfg.Port)
	}

	return nil
}

func (s *AppService) StartApp(name string) error {
	if err := s.Docker.StartApp(name); err != nil {
		return fmt.Errorf("error al iniciar %s: %w", name, err)
	}
	return nil
}

func (s *AppService) StopApp(name string) error {
	if err := s.Docker.StopApp(name); err != nil {
		return fmt.Errorf("error al detener %s: %w", name, err)
	}
	return nil
}

func (s *AppService) RemoveApp(name string) error {
	if err := s.Docker.RemoveApp(name); err != nil {
		return fmt.Errorf("error al eliminar %s: %w", name, err)
	}
	return nil
}

func (s *AppService) GetAppStatus(name string) (string, error) {
	status, err := s.Docker.GetAppStatus(name)
	if err != nil {
		return "", fmt.Errorf("error obteniendo estado de %s: %w", name, err)
	}
	return status, nil
}

// ------------------------- utils -------------------------

func install(app domain.AppConfig) error {
	fmt.Printf("%+v\n", app)
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

func setValues(app domain.AppConfig) map[string]string {
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

	if app.Postgres != nil {
		values["$postgres-user"] = app.Postgres.User         // UserName Postgres
		values["$postgres-password"] = app.Postgres.Password // PasswordUser Postgres
		values["$postgres-database"] = app.Postgres.Database // DatabaseName Postgres
		values["$postgres-port"] = app.Postgres.Port         // Port Postgres
	}

	if app.Template == "bookstack" {
		values["$app-key"] = utils.GenerateRandomString(64) // AppKey BookStack
	}

	return values
}

func askAppConfig(reader *bufio.Reader, ensureUnique func(string) error) domain.AppConfig {
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

	defaultTemplate := "bookstack"

	template := utils.AskInput(reader, "📦 Tipo de template (bookstack, nextcloud, redis, mysql, postgres)", defaultTemplate)

	if template == "mysql" {
		mysqlCfg := askMySQLConfig(reader, name)
		return domain.AppConfig{
			Name:     name,
			Template: template,
			Port:     mysqlCfg.Port,
			MySQL:    mysqlCfg,
		}
	}

	if template == "postgres" {
		postgresCfg := askMyPostgresConfig(reader, name)
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
		mysqlCfg = askMySQLConfig(reader, name)
	}

	return domain.AppConfig{
		Name:     name,
		Template: template,
		Port:     port,
		MySQL:    mysqlCfg,
	}
}

func askMySQLConfig(reader *bufio.Reader, name string) *domain.MySQLConfig {
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

func askMyPostgresConfig(reader *bufio.Reader, name string) *domain.PostgresConfig {
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
