package main

import (
	"crypto/subtle"
	"embed"
	"io/fs"
	"net/http"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed html
var contextFS embed.FS

func initHTTP(host string) {
	htmlFS, _ := fs.Sub(contextFS, "html")

	e := echo.New()

	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	if flagAuth != "" {
		e.Use(middleware.BasicAuth(basicAuthValidator))
	}

	htmlHandler := http.FileServer(http.FS(htmlFS))

	e.GET("/*", echo.WrapHandler(htmlHandler))
	e.POST("/api/list", getListHandler)

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

func getListHandler(c echo.Context) error {
	var body map[string]interface{}
	if err := c.Bind(&body); err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	if reflect.TypeOf(body["path"]).Kind() != reflect.String {
		return end(c, http.StatusBadRequest, -1, "路径格式不正确", nil)
	}
	targetPath, err := getTargetPath(body["path"].(string))
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	if pathNotExist(targetPath) {
		return end(c, http.StatusNotFound, -1, "指定路径不存在", nil)
	}
	isDir, err := pathIsDir(targetPath)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	if !isDir {
		return end(c, http.StatusOK, 1, "指定路径不是一个文件夹", nil)
	}
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return end(c, http.StatusInternalServerError, -1, err.Error(), nil)
	}
	var data []interface{}
	var dirs []interface{}
	var files []interface{}
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
