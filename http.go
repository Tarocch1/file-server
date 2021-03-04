package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"reflect"
)

//go:embed html
var contextFS embed.FS

func initHTTP(host string) {
	htmlFS, _ := fs.Sub(contextFS, "html")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(htmlFS)))
	mux.Handle("/api/list", http.HandlerFunc(basicAuth(getListHandler)))

	if flagHTTPSCert != "" && flagHTTPSKey != "" {
		log.Fatal(http.ListenAndServeTLS(host, flagHTTPSCert, flagHTTPSKey, mux))
	} else {
		log.Fatal(http.ListenAndServe(host, mux))
	}
}

func end(w http.ResponseWriter, status int, code int, message string, data interface{}) {
	if message == "" {
		message = http.StatusText(status)
	}
	if code == -1 {
		code = status
	}
	dataBytes, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(dataBytes)
}

func getListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		end(w, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, "", nil)
		return
	}
	bodyBytes, _ := io.ReadAll(r.Body)
	var body map[string]interface{}
	err := json.Unmarshal(bodyBytes, &body)
	if err != nil {
		end(w, http.StatusInternalServerError, -1, err.Error(), nil)
		return
	}
	if reflect.TypeOf(body["path"]).Kind() != reflect.String {
		end(w, http.StatusBadRequest, -1, "", nil)
		return
	}
	targetPath, err := getTargetPath(body["path"].(string))
	if err != nil {
		end(w, http.StatusInternalServerError, -1, err.Error(), nil)
		return
	}
	if pathNotExist(targetPath) {
		end(w, http.StatusNotFound, -1, "", nil)
		return
	}
	isDir, err := pathIsDir(targetPath)
	if err != nil {
		end(w, http.StatusInternalServerError, -1, err.Error(), nil)
		return
	}
	if !isDir {
		end(w, http.StatusOK, 1, "指定路径不是一个文件夹", nil)
		return
	}
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		end(w, http.StatusInternalServerError, -1, err.Error(), nil)
		return
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
	end(w, http.StatusOK, 0, "", data)
}
