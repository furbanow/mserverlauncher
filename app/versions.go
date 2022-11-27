package app

import (
	"os"
	"path"
	"strings"
)

func LoadVersions(config Config) []string {
	var versionFiles []string
	files, _ := os.ReadDir(path.Join(config.RootPath, config.VersionsPath))
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".jar") {
			versionFiles = append(versionFiles, file.Name())
		}
	}
	return versionFiles
}
