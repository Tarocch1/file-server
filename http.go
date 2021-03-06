package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	targetPath, err := getTargetPath(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	switch r.Method {
	case http.MethodHead:
		headHandler(w, r, targetPath)
		break
	case http.MethodGet:
		getHandler(w, r, targetPath)
		break
	case http.MethodPost:
		postHandler(w, r, targetPath)
		break
	case http.MethodPut:
		putHandler(w, r, targetPath)
		break
	case http.MethodDelete:
		deleteHandler(w, r, targetPath)
		break
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func headHandler(w http.ResponseWriter, r *http.Request, targetPath string) {
	if pathNotExist(targetPath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, targetPath)
}

func getHandler(w http.ResponseWriter, r *http.Request, targetPath string) {
	if pathNotExist(targetPath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	isDir, err := pathIsDir(targetPath)
	if err != nil {
		errorHandler(w, err)
		return
	}
	if isDir {
		allFiles, err := ioutil.ReadDir(targetPath)
		if err != nil {
			errorHandler(w, err)
			return
		}
		var data listTemplateDataStruct
		data.Title = targetPath
		data.Files = []listTemplateDataFileStruct{}
		var dirs, files []listTemplateDataFileStruct
		for _, fileInfo := range allFiles {
			if fileInfo.IsDir() {
				dirs = append(dirs, listTemplateDataFileStruct{
					Name:  fileInfo.Name(),
					IsDir: fileInfo.IsDir(),
					Time:  fileInfo.ModTime().Unix(),
					Size:  fileInfo.Size(),
				})
			} else {
				files = append(files, listTemplateDataFileStruct{
					Name:  fileInfo.Name(),
					IsDir: fileInfo.IsDir(),
					Time:  fileInfo.ModTime().Unix(),
					Size:  fileInfo.Size(),
				})
			}
		}
		data.Files = append(data.Files, dirs...)
		data.Files = append(data.Files, files...)
		t, err := template.New("listTemplate").Funcs(template.FuncMap{"formatSize": formatSize}).Parse(listTemplate)
		if err != nil {
			errorHandler(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		t.Execute(w, data)
		return
	}
	http.ServeFile(w, r, targetPath)
}

func postHandler(w http.ResponseWriter, r *http.Request, targetPath string) {
	if !pathNotExist(targetPath) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if len(r.Header["X-File-Server-Mkdir"]) > 0 {
		err := os.MkdirAll(targetPath, 0666)
		if err != nil {
			errorHandler(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	fileBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, err)
		return
	}
	err = ioutil.WriteFile(targetPath, fileBytes, 0666)
	if err != nil {
		errorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func putHandler(w http.ResponseWriter, r *http.Request, targetPath string) {
	notExist := pathNotExist(targetPath)
	fileBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, err)
		return
	}
	if len(r.Header["X-File-Server-Rename"]) > 0 {
		if notExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = os.Rename(targetPath, filepath.Join(filepath.Dir(targetPath), filepath.Base(string(fileBytes))))
		if err != nil {
			errorHandler(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	err = ioutil.WriteFile(targetPath, fileBytes, 0666)
	if err != nil {
		errorHandler(w, err)
		return
	}
	if notExist {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request, targetPath string) {
	if pathNotExist(targetPath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := os.RemoveAll(targetPath)
	if err != nil {
		errorHandler(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func errorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
