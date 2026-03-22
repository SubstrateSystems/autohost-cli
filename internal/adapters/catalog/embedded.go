// Package catalog provides a CatalogRepository backed by the embedded docker
// templates (assets/docker/*). No database is required.
package catalog

import (
	"autohost-cli/assets"
	"autohost-cli/internal/domain"
	"context"
	"fmt"
	"io/fs"
)

// entry holds static metadata for each embedded app template.
type entry struct {
	description   string
	defaultPort   string
	defaultPortDB string
	clientDB      string // "mysql" | "postgres" | ""
}

// registry maps template folder names to metadata.
var registry = map[string]entry{
	"bookstack": {
		description:   "BookStack — wiki y base de conocimiento",
		defaultPort:   "6875",
		defaultPortDB: "3306",
		clientDB:      "mysql",
	},
	"excalidraw": {
		description: "Excalidraw — pizarras colaborativas",
		defaultPort: "3000",
	},
	"joplin": {
		description:   "Joplin Server — notas sincronizadas",
		defaultPort:   "22300",
		defaultPortDB: "5432",
		clientDB:      "postgres",
	},
	"minio": {
		description:   "MinIO — almacenamiento de objetos compatible con S3",
		defaultPort:   "9000",
		defaultPortDB: "9001",
	},
	"mysql": {
		description: "MySQL — base de datos relacional",
		defaultPort: "3306",
	},
	"nextcloud": {
		description:   "Nextcloud — nube personal de archivos",
		defaultPort:   "8080",
		defaultPortDB: "3306",
		clientDB:      "mysql",
	},
	"postgres": {
		description: "PostgreSQL — base de datos relacional avanzada",
		defaultPort: "5432",
	},
	"redis": {
		description: "Redis — caché y broker de mensajes",
		defaultPort: "6379",
	},
}

// Adapter implements ports.CatalogRepository using the embedded FS.
type Adapter struct{}

func New() *Adapter { return &Adapter{} }

// ListApps returns all apps present in the embedded assets.
func (a *Adapter) ListApps(_ context.Context) ([]domain.CatalogApp, error) {
	entries, err := fs.ReadDir(assets.FS(), "docker")
	if err != nil {
		return nil, fmt.Errorf("catalog: read embedded docker dir: %w", err)
	}

	apps := make([]domain.CatalogApp, 0, len(entries))
	for i, e := range entries {
		if !e.IsDir() {
			continue
		}
		meta := registry[e.Name()]
		apps = append(apps, domain.CatalogApp{
			ID:            int64(i + 1),
			Name:          e.Name(),
			Description:   meta.description,
			DefaultPort:   meta.defaultPort,
			DefaultPortDB: meta.defaultPortDB,
			ClientDB:      meta.clientDB,
		})
	}
	return apps, nil
}

// FindByName returns the catalog entry for the named app.
func (a *Adapter) FindByName(ctx context.Context, name domain.AppName) (domain.CatalogApp, error) {
	apps, err := a.ListApps(ctx)
	if err != nil {
		return domain.CatalogApp{}, err
	}
	for _, app := range apps {
		if app.Name == string(name) {
			return app, nil
		}
	}
	return domain.CatalogApp{}, fmt.Errorf("catalog: app %q not found", name)
}
