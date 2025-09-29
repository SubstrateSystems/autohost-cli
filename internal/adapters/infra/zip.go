package infra

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipSingleBinary(zipPath, wantName, destDir string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	var outPath string
	for _, f := range r.File {
		name := f.Name
		if wantName != "" && !strings.EqualFold(name, wantName) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		outPath = filepath.Join(destDir, filepath.Base(name))
		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			return "", err
		}
		out.Close()
		break
	}
	if outPath == "" {
		return "", fmt.Errorf("no se encontró el binario %q dentro del zip", wantName)
	}
	return outPath, nil
}
