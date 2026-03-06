// Package implements an internal mechanism to communicate with an impress terminal.
package itrun

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

const (
	envVar      = "IMPRESS_TERMINAL_PATH"
	defaultName = "it"
)

var ErrITNotFound = errors.New("it executable not found")

type itRun struct {
	path string
	cmd  *exec.Cmd
}

func DefaultPath() (string, error) {
	var filenames []string
	if path := os.Getenv(envVar); path != "" {
		filenames = append(filenames, path)
	}
	if execName, err := os.Executable(); err == nil {
		filenames = append(filenames, path.Join(path.Dir(execName), defaultName))
	}
	filenames = append(filenames, fmt.Sprintf("./%s", defaultName))

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

	return "", ErrITNotFound
}

func New(path string) *itRun {
	return &itRun{
		path: path,
	}
}

func (r *itRun) Run(suffix string) error {
	r.cmd = exec.Command(r.path, suffix)
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr
	if err := r.cmd.Start(); err != nil {
		return fmt.Errorf("cmd.Start: %w", err)
	}
	return nil
}

func (r *itRun) Wait() error {
	if err := r.cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait: %w", err)
	}
	return nil
}
