package config

import (
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

//go:embed *.toml
var embedded embed.FS

// cache en memoria para evitar reparseos
var (
	mu    sync.RWMutex
	files = map[string]map[string]any{} // nombre archivo -> mapa TOML (secciones -> claves)
)

// loadFile carga y cachea un archivo TOML embebido.
// file puede ser "urls.toml" (con o sin "config/").
func loadFile(file string) (map[string]any, error) {
	base := filepath.Base(file) // "urls.toml"
	path := filepath.Join("config", base)

	mu.RLock()
	if m, ok := files[path]; ok {
		mu.RUnlock()
		return m, nil
	}
	mu.RUnlock()

	data, err := embedded.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer %s embebido: %w", path, err)
	}

	var m map[string]any
	if err := toml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("no se pudo parsear %s: %w", path, err)
	}

	mu.Lock()
	files[path] = m
	mu.Unlock()
	return m, nil
}

// Get devuelve una interfaz (any) de section/key. Útil si quieres convertir tú.
func Get(file, section, key string) (any, error) {
	m, err := loadFile(file)
	if err != nil {
		return nil, err
	}
	secRaw, ok := m[section]
	if !ok {
		return nil, fmt.Errorf("sección %q no encontrada en %q", section, file)
	}
	sec, ok := secRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("sección %q en %q no es un objeto", section, file)
	}
	val, ok := sec[key]
	if !ok {
		return nil, fmt.Errorf("clave %q no encontrada en sección %q de %q", key, section, file)
	}
	return val, nil
}

// Helpers tipados
func GetString(file, section, key string) (string, error) {
	v, err := Get(file, section, key)
	if err != nil {
		return "", err
	}
	switch t := v.(type) {
	case string:
		return t, nil
	case fmt.Stringer:
		return t.String(), nil
	case int64: // números en TOML pueden venir como int64
		return strconv.FormatInt(t, 10), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case bool:
		if t {
			return "true", nil
		}
		return "false", nil
	default:
		return "", errors.New("valor no es string convertible")
	}
}

func GetInt(file, section, key string) (int, error) {
	v, err := Get(file, section, key)
	if err != nil {
		return 0, err
	}
	switch t := v.(type) {
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		i, convErr := strconv.Atoi(t)
		if convErr != nil {
			return 0, fmt.Errorf("no se pudo convertir a int: %w", convErr)
		}
		return i, nil
	default:
		return 0, errors.New("valor no es numérico")
	}
}

func GetBool(file, section, key string) (bool, error) {
	v, err := Get(file, section, key)
	if err != nil {
		return false, err
	}
	switch t := v.(type) {
	case bool:
		return t, nil
	case string:
		switch t {
		case "true", "TRUE", "1", "yes", "si", "sí":
			return true, nil
		case "false", "FALSE", "0", "no":
			return false, nil
		}
	}
	return false, errors.New("valor no es booleano")
}

// Must* variantes que paniquean si falla (útiles para defaults obligatorios)
func MustString(file, section, key string) string {
	s, err := GetString(file, section, key)
	if err != nil {
		panic(err)
	}
	return s
}
