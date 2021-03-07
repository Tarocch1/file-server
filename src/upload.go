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
	targetPath, _ := getTargetPath(body["path"].(string))
	if isDir, _ := pathIsDir(targetPath); isDir {
		return endWithError(c, 3)
	}
	os.RemoveAll(targetPath)
	os.MkdirAll(targetPath, 0777)
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}

func uploadChunkHandler(c echo.Context) error {
	path := c.QueryParam("path")
	id := c.QueryParam("id")
	targetPath, _ := getTargetPath(path)
	bytes, _ := io.ReadAll(c.Request().Body)
	filePath := filepath.Join(targetPath, filepath.Base(targetPath)+"."+id)
	os.WriteFile(filePath, bytes, 0777)
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}

func uploadEndHandler(c echo.Context) error {
	var body map[string]interface{}
	c.Bind(&body)
	targetPath, _ := getTargetPath(body["path"].(string))
	tempTargetPath := targetPath + ".UPLOADING"
	file, _ := os.Create(tempTargetPath)
	defer file.Close()
	var id int64 = 0
	for {
		filePath := filepath.Join(targetPath, filepath.Base(targetPath)+"."+strconv.FormatInt(id, 10))
		if pathNotExist(filePath) {
			break
		}
		bytes, _ := os.ReadFile(filePath)
		file.Write(bytes)
		id++
	}
	os.RemoveAll(targetPath)
	os.Rename(tempTargetPath, targetPath)
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}
