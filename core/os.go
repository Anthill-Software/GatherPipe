package core

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
)

type SystemInfo struct {
	OS   string
	Arch string
}

// GetSystemInfo retourne les infos pour le catalogue
func GetSystemInfo() SystemInfo {
	return SystemInfo{
		OS:   runtime.GOOS,   // "linux", "windows", "darwin"...
		Arch: runtime.GOARCH, // "amd64", "arm64", "arm"...
	}
}

// resolveURL remplace les variables {{.OS}} etc par les valeurs système
func resolveURL(template string, version string, pluginDir string, pluginName string) (string, string) {
	pluginName = filepath.Base(pluginName)
	osExt := ""
	if runtime.GOOS == "windows" {
		osExt = ".exe"
	}
	r := strings.NewReplacer(
		"{{.Version}}", version,
		"{{.OS}}", runtime.GOOS,
		"{{.Arch}}", runtime.GOARCH,
		"{{.Ext}}", osExt,
	)
	destPath := filepath.Join(pluginDir, fmt.Sprintf("%s_%s%s", pluginName, version, osExt))
	url := r.Replace(template)
	slog.Debug(fmt.Sprintf("[URL] %s: %s", pluginName, url))
	return url, destPath
}
