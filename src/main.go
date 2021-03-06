package main

import (
	"flag"
	"fmt"
	"log"
)

var workDir string
var (
	version   = ""
	goVersion = ""
	buildTime = ""
	commitID  = ""
)

func init() {
	initFlag()
}

func main() {
	flag.Parse()

	if flagVersion {
		fmt.Println("Version:", version)
		fmt.Println("Go Version:", goVersion)
		fmt.Println("Build Time:", buildTime)
		fmt.Println("Git Commit ID:", commitID)
		return
	}

	initAuth()

	var err error
	workDir, err = getWorkDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	host := flagHost + ":" + flagPort

	initHTTP(host)
}
