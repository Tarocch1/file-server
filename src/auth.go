package main

import (
	"strings"
)

var flagAuthUsername string
var flagAuthPassword string

func initAuth() {
	if flagAuth != "" {
		strArray := strings.Split(flagAuth, ":")
		flagAuthUsername = strArray[0]
		flagAuthPassword = strArray[1]
	}
}
