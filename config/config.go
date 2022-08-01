package config

import "os"

var pwd, _ = os.Getwd()
var baseURL = pwd + "/upload/"

func getBaseURL() string {
	return baseURL
}
