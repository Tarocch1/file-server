package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

func getWorkDir() (string, error) {
	if filepath.IsAbs(flagDir) {
		return flagDir, nil
	}
	dirPath, err := filepath.Abs(flagDir)
	if err != nil {
		return "", err
	}
	isDir, err := pathIsDir(dirPath)
	if err != nil {
		return "", err
	}
	if isDir {
		return dirPath, nil
	}
	return "", errors.New(dirPath + " is not a dir.")
}

func getTargetPath(path string) (string, error) {
	return filepath.Abs(filepath.Join(workDir, path))
}

func pathNotExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}

func pathIsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func formatSize(intSize int64) string {
	size := float64(intSize)
	if size < 1024 {
		return strconv.FormatFloat(size, 'f', 2, 64) + " B"
	}
	if size /= 1024; size < 1024 {
		return strconv.FormatFloat(size, 'f', 2, 64) + " KB"
	}
	if size /= 1024; size < 1024 {
		return strconv.FormatFloat(size, 'f', 2, 64) + " MB"
	}
	if size /= 1024; size < 1024 {
		return strconv.FormatFloat(size, 'f', 2, 64) + " GB"
	}
	if size /= 1024; size < 1024 {
		return strconv.FormatFloat(size, 'f', 2, 64) + " TB"
	}
	size /= 1024
	return strconv.FormatFloat(size, 'f', 2, 64) + " PB"
}
