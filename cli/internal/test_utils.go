package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"syscall"
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

func GetProcDir() string {
	if uid := syscall.Getuid(); uid == 0 {
		return fmt.Sprintf("/var/run/digkala-api/%d", os.Getpid())
	} else {
		return fmt.Sprintf("/run/user/%d/digikala-api/%d", uid, os.Getpid())
	}
}

func AssertDir(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return errors.New("path is not a directory: " + path)
		}

		return nil
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.New("failed to create directory: " + path)
		}
		return nil
	}

	return err
}
