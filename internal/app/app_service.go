package app

import (
	"autohost-cli/assets"
	"autohost-cli/internal/domain"
	"autohost-cli/internal/ports"
	"autohost-cli/utils"
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func (s AppService) ListCatalog(ctx context.Context) ([]domain.CatalogApp, error) {
	return s.Catalog.ListApps(ctx)
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

	if settings.Minio != nil {
		values["$minio-user"] = settings.Minio.User
		values["$minio-password"] = settings.Minio.Password
		values["$console-port"] = settings.Minio.ConsolePort
		values["$minio-data-path"] = settings.Minio.DataPath
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

	if appTemplate.Name == "minio" {
		minioCfg := askMinioConfig(reader, port)
		appConfig = domain.AppConfig{
			AppSettings: domain.InstalledApp{Name: nameApp, Port: port, PortDB: minioCfg.ConsolePort, Template: appTemplate.Name, CatalogAppID: int64(appTemplate.ID)},
			Minio:       minioCfg,
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

func askMinioConfig(reader *bufio.Reader, apiPort string) *domain.MinioConfig {
	fmt.Println("\n⚙️  Configuración de MinIO:")

	user := utils.AskInput(reader, "Usuario admin", "minioadmin")
	pass := utils.AskInput(reader, "Contraseña admin", utils.GenerateRandomString(20))
	consolePort := utils.AskAppPort(reader, "Puerto de la consola web", "9001")

	dataPath := askDiskDataPath(reader)

	return &domain.MinioConfig{
		User:        user,
		Password:    pass,
		ConsolePort: consolePort,
		DataPath:    dataPath,
	}
}

// askDiskDataPath lists available block devices and asks the user where MinIO
// should store its data.  The user can pick an external disk (the CLI will
// ensure it is mounted) or type a custom path.
func askDiskDataPath(reader *bufio.Reader) string {
	disks := listExternalDisks()

	fmt.Println()
	if len(disks) == 0 {
		fmt.Println("⚠️  No se detectaron discos externos.")
	} else {
		fmt.Println("💾 Discos externos detectados:")
		for i, d := range disks {
			fmt.Printf("  [%d] %s  %s  %s\n", i+1, d.Device, d.Size, d.Label)
		}
	}

	fmt.Println("  [0] Introducir ruta manualmente")
	fmt.Println()

	for {
		raw := utils.AskInput(reader, "Elige un disco [número] o escribe una ruta", "0")
		raw = strings.TrimSpace(raw)

		// Numeric selection
		if idx := parseChoice(raw, len(disks)); idx > 0 {
			disk := disks[idx-1]
			mountPoint := ensureDiskMounted(disk)
			dataPath := filepath.Join(mountPoint, "minio-data")
			fmt.Printf("📂 Los datos de MinIO se guardarán en: %s\n", dataPath)
			return dataPath
		}

		// Manual path or "0"
		if raw == "0" || raw == "" {
			path := utils.AskInput(reader, "Ruta del directorio de datos", "/opt/minio/data")
			path = strings.TrimSpace(path)
			if path != "" {
				fmt.Printf("📂 Los datos de MinIO se guardarán en: %s\n", path)
				return path
			}
			fmt.Println("❌ La ruta no puede estar vacía.")
			continue
		}

		// Treat as a direct path input
		if strings.HasPrefix(raw, "/") {
			fmt.Printf("📂 Los datos de MinIO se guardarán en: %s\n", raw)
			return raw
		}

		fmt.Println("❌ Opción no válida. Elige un número de la lista o escribe una ruta absoluta.")
	}
}

// diskInfo holds basic info about a block device gathered from lsblk.
type diskInfo struct {
	Device string
	Size   string
	Label  string
}

// listExternalDisks returns non-system block devices using lsblk.
// virtualPrefixes lists device name prefixes that represent virtual or
// in-memory block devices that should never appear as storage options.
var virtualPrefixes = []string{"zram", "loop", "ram", "dm-"}

func listExternalDisks() []diskInfo {
	// Include MOUNTPOINT so we can detect the system/boot disk.
	out, err := exec.Command("lsblk", "-o", "NAME,SIZE,LABEL,TYPE,MOUNTPOINT", "--json", "--nodeps").Output()
	if err != nil {
		// lsblk not available or failed — return empty
		return nil
	}
	disks := parseLsblkJSON(out)

	// Determine which block device backs the root filesystem so we can exclude it.
	sysDisk := systemDiskName()

	var result []diskInfo
	for _, d := range disks {
		name := filepath.Base(d.Device)
		// Skip virtual/memory devices.
		if isVirtualDevice(name) {
			continue
		}
		// Skip the disk that contains the root/boot filesystem.
		if sysDisk != "" && (name == sysDisk || strings.HasPrefix(name, sysDisk)) {
			continue
		}
		result = append(result, d)
	}
	return result
}

func isVirtualDevice(name string) bool {
	for _, p := range virtualPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// systemDiskName returns the base device name (e.g. "sda", "mmcblk0") that
// contains the root filesystem, so it can be excluded from the disk menu.
func systemDiskName() string {
	// lsblk -no PKNAME / — prints the parent kernel name of the root partition.
	out, err := exec.Command("lsblk", "-no", "PKNAME", "/").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// parseLsblkJSON parses lsblk --json output and returns disk entries.
// Uses a lightweight line-scan so no extra import is needed.
func parseLsblkJSON(data []byte) []diskInfo {
	lines := strings.Split(string(data), "\n")
	var result []diskInfo
	var cur diskInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if v, ok := jsonStringField(line, "name"); ok {
			cur.Device = "/dev/" + v
		}
		if v, ok := jsonStringField(line, "size"); ok {
			cur.Size = v
		}
		if v, ok := jsonStringField(line, "label"); ok {
			cur.Label = v
		}
		if v, ok := jsonStringField(line, "type"); ok && v == "disk" {
			// Flush the current entry when we hit a "type": "disk" field.
			result = append(result, cur)
			cur = diskInfo{}
		}
	}
	return result
}

func jsonStringField(line, key string) (string, bool) {
	prefix := `"` + key + `": "`
	idx := strings.Index(line, prefix)
	if idx == -1 {
		return "", false
	}
	rest := line[idx+len(prefix):]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return "", false
	}
	return rest[:end], true
}

// diskPartition holds info about a single partition on a disk.
type diskPartition struct {
	device     string // e.g. /dev/sda1
	mountPoint string // non-empty if already mounted
	fsType     string // e.g. ext4, xfs, vfat
	size       string
}

// listPartitions returns the partitions of diskDevice (e.g. /dev/sda).
func listPartitions(diskDevice string) []diskPartition {
	out, err := exec.Command("lsblk", "-no", "NAME,FSTYPE,SIZE,MOUNTPOINT", diskDevice).Output()
	if err != nil {
		return nil
	}
	var parts []diskPartition
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	// First line is the disk itself; the rest are partitions/children.
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		p := diskPartition{device: "/dev/" + strings.TrimLeft(fields[0], "└─├─")}
		if len(fields) >= 2 {
			p.fsType = fields[1]
		}
		if len(fields) >= 3 {
			p.size = fields[2]
		}
		if len(fields) >= 4 {
			p.mountPoint = fields[3]
		}
		parts = append(parts, p)
	}
	return parts
}

// resolveTarget picks the block device (partition or whole disk) to mount.
// If a partition is already mounted it returns that mount point directly.
func resolveTarget(disk diskInfo) (device string, alreadyMounted string) {
	parts := listPartitions(disk.Device)

	// If any partition is already mounted, use it immediately.
	for _, p := range parts {
		if p.mountPoint != "" {
			return p.device, p.mountPoint
		}
	}

	// Prefer the first partition with a known filesystem; otherwise use the
	// whole disk (unpartitioned/GPT-raw scenario).
	for _, p := range parts {
		if p.fsType != "" {
			return p.device, ""
		}
	}
	if len(parts) > 0 {
		return parts[0].device, ""
	}
	return disk.Device, ""
}

// ensureDiskMounted returns the mount point for the given disk, resolving
// the correct partition and using sudo when needed.
func ensureDiskMounted(disk diskInfo) string {
	target, existingMount := resolveTarget(disk)

	if existingMount != "" {
		fmt.Printf("✅ El disco %s ya está montado en %s\n", target, existingMount)
		return existingMount
	}

	devName := filepath.Base(disk.Device)
	mountPoint := filepath.Join("/mnt", devName)

	fmt.Printf("🔧 Montando %s en %s ...\n", target, mountPoint)

	// Create mount directory — try directly, fall back to sudo.
	if err := os.MkdirAll(mountPoint, 0o755); err != nil {
		if sudoErr := exec.Command("sudo", "mkdir", "-p", mountPoint).Run(); sudoErr != nil {
			fmt.Printf("⚠️  No se pudo crear %s (ni con sudo): %v\n", mountPoint, sudoErr)
			return mountPoint
		}
	}

	// Mount — always use sudo since mounting requires root.
	sudoMount := exec.Command("sudo", "mount", target, mountPoint)
	sudoMount.Stdout = os.Stdout
	sudoMount.Stderr = os.Stderr
	if err := sudoMount.Run(); err != nil {
		fmt.Printf("⚠️  mount falló. Monta el disco manualmente en %s y vuelve a ejecutar.\n", mountPoint)
	} else {
		fmt.Printf("✅ Disco montado en %s\n", mountPoint)
	}
	return mountPoint
}

// parseChoice converts a 1-based numeric string to an int; returns 0 when out-of-range or not numeric.
func parseChoice(s string, max int) int {
	if max == 0 {
		return 0
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	if n < 1 || n > max {
		return 0
	}
	return n
}
