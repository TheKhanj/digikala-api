package internal

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		if dir == filepath.Dir(dir) {
			log.Fatalf("go.mod not found in any parent directories")
		}

		dir = filepath.Dir(dir)
	}
}

func GetAbsPath(name string) string {
	env := os.Getenv("env")

	if !(env == "test" || env == "dev") {
		return name
	}

	return path.Join(GetProjectRoot(), name)
}

func AssertAtLeastOneFile(dir string) error {
	entries, err := os.ReadDir(GetAbsPath(dir))
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return nil
		}
	}

	return fmt.Errorf(
		"no files found in directory: %s", dir,
	)
}
