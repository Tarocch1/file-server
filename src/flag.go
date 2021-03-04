package main

import "flag"

var flagHost string

var flagPort string

var flagDir string

var flagAuth string

var flagHTTPSCert string

var flagHTTPSKey string

var flagVersion bool

func initFlag() {
	flag.StringVar(&flagHost, "h", "0.0.0.0", "Host to listen.")
	flag.StringVar(&flagPort, "p", "8080", "Port to listen.")
	flag.StringVar(&flagDir, "d", ".", "Dir path to serve.")
	flag.StringVar(&flagAuth, "a", "", "<username:password> Basic auth user.")
	flag.StringVar(&flagHTTPSCert, "https-cert", "", "Path to https cert.")
	flag.StringVar(&flagHTTPSKey, "https-key", "", "Path to https key.")
	flag.BoolVar(&flagVersion, "v", false, "Print version information.")
}
