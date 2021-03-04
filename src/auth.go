package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func initAuth() {
	if flagAuth != "" {
		flagAuth = base64.StdEncoding.EncodeToString([]byte(flagAuth))
	}
}

func basicAuth(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if flagAuth == "" {
			f(w, r)
			return
		}
		basicAuthPrefix := "Basic "
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			token := auth[len(basicAuthPrefix):]
			if token == flagAuth {
				f(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
	}
}
