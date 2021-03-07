package main

import (
	"crypto/subtle"
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _errors = map[int]string{
	1: "指定路径不存在",
	2: "指定路径不是一个文件夹",
}

//go:embed html
var contextFS embed.FS

func initHTTP(host string) {
	htmlFS, _ := fs.Sub(contextFS, "html")

	e := echo.New()

	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	if flagAuth != "" {
		e.Use(middleware.BasicAuth(basicAuthValidator))
	}

	htmlHandler := http.FileServer(http.FS(htmlFS))

	e.GET("/*", echo.WrapHandler(htmlHandler))
	e.POST("/api/list", getListHandler)
	e.POST("/api/remove", removeHandler)

	if flagHTTPSCert != "" && flagHTTPSKey != "" {
		e.Logger.Fatal(e.StartTLS(host, flagHTTPSCert, flagHTTPSKey))
	} else {
		e.Logger.Fatal(e.Start(host))
	}
}

func basicAuthValidator(username, password string, c echo.Context) (bool, error) {
	if subtle.ConstantTimeCompare([]byte(username), []byte(flagAuthUsername)) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte(flagAuthPassword)) == 1 {
		return true, nil
	}
	return false, nil
}

func end(c echo.Context, status, code int, message string, data interface{}) error {
	if message == "" {
		message = http.StatusText(status)
	}
	if code == -1 {
		code = status
	}
	return c.JSON(status, map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func endWithError(c echo.Context, code int) error {
	return end(c, http.StatusOK, code, _errors[code], nil)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	status := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		status = he.Code
	}
	if err := end(c, status, -1, err.Error(), nil); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func getListHandler(c echo.Context) error {
	var body map[string]interface{}
	c.Bind(&body)
	targetPath, _ := getTargetPath(body["path"].(string))
	if pathNotExist(targetPath) {
		return endWithError(c, 1)
	}
	if isDir, _ := pathIsDir(targetPath); !isDir {
		return endWithError(c, 2)
	}
	entries, _ := os.ReadDir(targetPath)
	data := []interface{}{}
	dirs := []interface{}{}
	files := []interface{}{}
	for _, entry := range entries {
		fileInfo, _ := entry.Info()
		var item = map[string]interface{}{
			"name":  fileInfo.Name(),
			"isDir": fileInfo.IsDir(),
			"time":  fileInfo.ModTime().Unix(),
			"size":  formatSize(fileInfo.Size()),
		}
		if fileInfo.IsDir() {
			dirs = append(dirs, item)
		} else {
			files = append(files, item)
		}
	}
	data = append(data, dirs...)
	data = append(data, files...)
	return end(c, http.StatusOK, 0, "SUCCESS", data)
}

func removeHandler(c echo.Context) error {
	var body map[string]interface{}
	c.Bind(&body)
	targetPath, _ := getTargetPath(body["path"].(string))
	if pathNotExist(targetPath) {
		return endWithError(c, 1)
	}
	os.RemoveAll(targetPath)
	return end(c, http.StatusOK, 0, "SUCCESS", nil)
}
