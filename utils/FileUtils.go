package utils

import (
	"errors"
	"path/filepath"
	"runtime"
)

// Filename is the __filename equivalent
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

// Dirname is the __dirname equivalent
func Dirname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}

// Rootname get
func Rootname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(filename), "../"), nil
}
