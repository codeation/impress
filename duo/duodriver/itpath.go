package duodriver

import (
	"errors"
	"fmt"
	"os"
	"path"
)

var errITNotFound = errors.New("it executable not found")

func itPath() (string, error) {
	var filenames []string
	if path := os.Getenv("IMPRESS_TERMINAL_PATH"); path != "" {
		filenames = append(filenames, path)
	}
	if execName, err := os.Executable(); err == nil {
		filenames = append(filenames, path.Join(path.Dir(execName), "it"))
	}
	filenames = append(filenames, "./it")

	for _, filename := range filenames {
		stat, err := os.Stat(filename)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("os.Stat: %w", err)
		}
		if stat.Mode().Type().IsRegular() {
			return filename, nil
		}
	}

	return "", errITNotFound
}
