package internal

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
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

func RandomListeningAddress() string {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	return "127.0.0.1:" + strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}
