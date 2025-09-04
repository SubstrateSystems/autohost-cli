package utils

import (
	"os"
	"path/filepath"
)

func GetAutohostDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// fallback razonable si falla
		return "/tmp/autohost"
	}
	return filepath.Join(home, ".autohost")
}

func GetSubdir(subdir string) string {
	return filepath.Join(GetAutohostDir(), subdir)
}

func IsInitialized() bool {
	// _, err := os.Stat(GetAutohostDir())
	return true
}
