package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func uploadStartHandler(c echo.Context) error {
	var body map[string]interface{}
	c.Bind(&body)
	targetPath, err := getTargetPath(body["path"].(string))
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	if isDir, err := pathIsDir(targetPath); err != nil && !os.IsNotExist(err) {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	} else if isDir {
		return endWithError(c, 3)
	}
	err = os.RemoveAll(targetPath)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	err = os.MkdirAll(targetPath, 0777)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}

func uploadChunkHandler(c echo.Context) error {
	path := c.QueryParam("path")
	id := c.QueryParam("id")
	targetPath, err := getTargetPath(path)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	bytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	filePath := filepath.Join(targetPath, filepath.Base(targetPath)+"."+id)
	err = os.WriteFile(filePath, bytes, 0777)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}

func uploadEndHandler(c echo.Context) error {
	var body map[string]interface{}
	c.Bind(&body)
	targetPath, err := getTargetPath(body["path"].(string))
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	tempTargetPath := targetPath + ".UPLOADING"
	file, err := os.Create(tempTargetPath)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	defer file.Close()
	var id int64 = 0
	for {
		filePath := filepath.Join(targetPath, filepath.Base(targetPath)+"."+strconv.FormatInt(id, 10))
		if pathNotExist(filePath) {
			break
		}
		bytes, err := os.ReadFile(filePath)
		if err != nil {
			return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
		}
		file.Write(bytes)
		id++
	}
	err = os.RemoveAll(targetPath)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	err = os.Rename(tempTargetPath, targetPath)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}
