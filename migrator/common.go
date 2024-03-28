package migrator

import (
	"log"
	"os"
	"path"
	"strings"
)

const (
	sqlFileSuffix = ".sql"
)

func getMigrationFiles(migrationPath string) ([]string, error) {
	dir, err := os.ReadDir(migrationPath)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		if !strings.HasSuffix(entryName, sqlFileSuffix) {
			log.Printf("[WARN] File '%s' does not have the suffix '%s'", entryName, sqlFileSuffix)
			continue
		}
		migrationFile := path.Join(migrationPath, entryName)
		migrationFiles = append(migrationFiles, migrationFile)
	}

	return migrationFiles, err
}
