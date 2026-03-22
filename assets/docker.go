package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed docker/**/*
var dockerFS embed.FS

// ReadCompose reads the compose file for an app template.
// It tries "compose.yml" first, then falls back to "docker-compose.yml".
func ReadCompose(app string) ([]byte, error) {
	for _, name := range []string{"compose.yml", "docker-compose.yml"} {
		data, err := fs.ReadFile(dockerFS, filepath.Join("docker", app, name))
		if err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("no compose file found for %q", app)
}

func ReadEnvExample(app string) ([]byte, error) {
	return fs.ReadFile(dockerFS, filepath.Join("docker", app, ".env.example"))
}

// FS returns the embedded docker assets filesystem.
func FS() fs.FS { return dockerFS }
