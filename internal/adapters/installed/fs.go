// Package installed provides an InstalledRepository backed by a JSON file at
// ~/.autohost/installed.json. No database is required.
package installed

import (
	"autohost-cli/internal/domain"
	"autohost-cli/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const storeFile = "installed.json"

// Adapter implements ports.InstalledRepository using a JSON file.
type Adapter struct {
	path string // full path to installed.json
}

func New() *Adapter {
	return &Adapter{path: filepath.Join(utils.GetAutohostDir(), storeFile)}
}

// -- helpers --

func (a *Adapter) load() ([]domain.InstalledApp, error) {
	data, err := os.ReadFile(a.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("installed: read store: %w", err)
	}
	var apps []domain.InstalledApp
	if err := json.Unmarshal(data, &apps); err != nil {
		return nil, fmt.Errorf("installed: parse store: %w", err)
	}
	return apps, nil
}

func (a *Adapter) save(apps []domain.InstalledApp) error {
	if err := os.MkdirAll(filepath.Dir(a.path), 0o700); err != nil {
		return fmt.Errorf("installed: create dir: %w", err)
	}
	data, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return fmt.Errorf("installed: marshal: %w", err)
	}
	return os.WriteFile(a.path, data, 0o600)
}

// -- interface --

func (a *Adapter) List(_ context.Context) ([]domain.InstalledApp, error) {
	return a.load()
}

func (a *Adapter) Install(_ context.Context, app domain.InstalledApp) error {
	apps, err := a.load()
	if err != nil {
		return err
	}
	app.CreatedAt = time.Now()
	apps = append(apps, app)
	return a.save(apps)
}

func (a *Adapter) Remove(_ context.Context, name domain.AppName) error {
	apps, err := a.load()
	if err != nil {
		return err
	}
	filtered := apps[:0]
	for _, a := range apps {
		if domain.AppName(a.Name) != name {
			filtered = append(filtered, a)
		}
	}
	return a.save(filtered)
}

func (a *Adapter) IsInstalled(_ context.Context, name domain.AppName) (bool, error) {
	apps, err := a.load()
	if err != nil {
		return false, err
	}
	for _, a := range apps {
		if domain.AppName(a.Name) == name {
			return true, nil
		}
	}
	return false, nil
}
