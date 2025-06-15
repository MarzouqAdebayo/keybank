package utils

import (
	"errors"
	"os"
)

func FileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func DirExists(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if err == nil {
		return info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0o700)
}
