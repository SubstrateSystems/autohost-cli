package app

import (
	"autohost-cli/assets"
	"autohost-cli/internal/domain"

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
	Docker    ports.Docker
	Installed ports.InstalledRepository
	Catalog   ports.CatalogRepository
}

func (s *AppService) InstallApp(ctx context.Context, appTemplate string) error {
	reader := bufio.NewReader(os.Stdin)

	ensureUnique := func(name string) error {
		exists, err := s.Installed.IsInstalled(ctx, domain.AppName(name))
		if err != nil {
			return fmt.Errorf("no se pudo validar el nombre: %w", err)
		}
		if exists {
			return fmt.Errorf("el nombre %q ya está en uso", name)
		}
		return nil
	}

	app, err := s.Catalog.FindByName(ctx, domain.AppName(appTemplate))
	if err != nil {
		return fmt.Errorf("error al buscar la aplicación en el catálogo: %w", err)
	}

	cfg := askAppConfig(reader, app, ensureUnique)

	if err := s.Installed.Install(ctx, cfg.AppSettings); err != nil {
		return fmt.Errorf("error al registrar la aplicación instalada: %w", err)
	}

	if err := install(cfg); err != nil {
		return fmt.Errorf("error al instalar %s: %w", cfg.AppSettings.Name, err)
	}

	startApp := utils.AskInput(reader, fmt.Sprintf("¿Deseas iniciar %s ahora? [Y/N]: ", cfg.AppSettings.Name), "Y")

	if strings.EqualFold(startApp, "Y") {
		if err := s.Docker.StartApp(cfg.AppSettings.Name); err != nil {
			return fmt.Errorf("error al iniciar %s: %w", cfg.AppSettings.Name, err)
		}
		fmt.Printf("🚀 La aplicación %s ha sido iniciada en http://localhost:%s\n", cfg.AppSettings.Name, cfg.AppSettings.Port)
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

func (s *AppService) RemoveApp(ctx context.Context, name domain.AppName) error {
	if err := s.Docker.RemoveApp(name); err != nil {

		return fmt.Errorf("error al eliminar %s: %w", name, err)
	}
	s.Installed.Remove(ctx, name)
	return nil
}

func (s *AppService) GetAppStatus(name string) (string, error) {
	status, err := s.Docker.GetAppStatus(name)
	if err != nil {
		return "", fmt.Errorf("error obteniendo estado de %s: %w", name, err)
	}
	return status, nil
}

func (s AppService) ListInstalled(ctx context.Context) ([]domain.InstalledApp, error) {
	return s.Installed.List(ctx)
}

func (s AppService) IsAppInstalled(ctx context.Context, name domain.AppName) (bool, error) {
	return s.Installed.IsInstalled(ctx, name)
}

type CatalogService struct {
	Catalog ports.CatalogRepository
}

func (s CatalogService) List(ctx context.Context) ([]domain.CatalogApp, error) {
	return s.Catalog.ListApps(ctx)
}

// ------------------------- utils -------------------------

func install(app domain.AppConfig) error {
	fmt.Printf("%+v\n", app)
	appDir := filepath.Join(utils.GetSubdir("apps"), app.AppSettings.Name)
	composePath := filepath.Join(appDir, "docker-compose.yml")

	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return fmt.Errorf("error creando directorio de destino: %w", err)
	}

	data, err := assets.ReadCompose(app.AppSettings.Template)
	if err != nil {
		return fmt.Errorf("no se encontró plantilla embebida para %s: %w", app.AppSettings.Template, err)
	}
	fmt.Println("📦 Usando plantilla embebida para:", app.AppSettings.Template)

	values := setValues(app)

	fmt.Println(values)

	final := utils.ReplacePlaceholders(string(data), values)

	if err := os.WriteFile(composePath, []byte(final), 0o644); err != nil {
		return fmt.Errorf("error escribiendo archivo docker-compose: %w", err)
	}
	fmt.Println("✅ Archivo docker-compose creado en:", composePath)
	return nil
}

func setValues(settings domain.AppConfig) map[string]string {
	values := map[string]string{}

	values["$service-name"] = settings.AppSettings.Name // AppName
	values["$port"] = settings.AppSettings.Port         // AppPort
	if settings.MySQL != nil {
		values["$mysql-user"] = settings.MySQL.User                  // UserName MySQL
		values["$mysql-password"] = settings.MySQL.Password          // PasswordUser MySQL
		values["$mysql-root-password"] = settings.MySQL.RootPassword // RootPassword MySQL
		values["$mysql-database"] = settings.MySQL.Database + "-db"  // DatabaseName MySQL
		values["$mysql-port"] = settings.MySQL.Port                  // Port MySQL
	}

	if settings.Postgres != nil {
		values["$postgres-user"] = settings.Postgres.User         // UserName Postgres
		values["$postgres-password"] = settings.Postgres.Password // PasswordUser Postgres
		values["$postgres-database"] = settings.Postgres.Database // DatabaseName Postgres
		values["$postgres-port"] = settings.Postgres.Port         // Port Postgres
	}

	if settings.AppSettings.Name == "bookstack" {
		values["$app-key"] = utils.GenerateRandomString(64) // AppKey BookStack
	}

	return values
}

func askAppConfig(reader *bufio.Reader, appTemplate domain.CatalogApp, ensureUnique func(string) error) domain.AppConfig {
	nameApp := "appdemo"
	appConfig := domain.AppConfig{}

	for {
		nameApp = utils.AskInput(reader, "📝 Nombre de la aplicación", nameApp)
		if err := ensureUnique(nameApp); err != nil {
			fmt.Printf("⚠️ %v\n", err)
			continue
		}
		break
	}
	port := utils.AskAppPort(reader, "🔌 Puerto del host a utilizar", appTemplate.DefaultPort)

	if appTemplate.ClientDB == "mysql" {
		mysqlCfg := askMySQLConfig(reader, nameApp)
		appConfig = domain.AppConfig{
			AppSettings: domain.InstalledApp{Name: nameApp, Port: port, PortDB: mysqlCfg.Port, Template: appTemplate.Name, CatalogAppID: int64(appTemplate.ID)},
			MySQL:       mysqlCfg,
		}
	}

	if appTemplate.Name == "postgres" {
		postgresCfg := askMyPostgresConfig(reader, nameApp)
		appConfig = domain.AppConfig{
			AppSettings: domain.InstalledApp{Name: nameApp, Port: port, PortDB: postgresCfg.Port, Template: appTemplate.Name, CatalogAppID: int64(appTemplate.ID)},
			Postgres:    postgresCfg,
		}
	}

	return appConfig
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
