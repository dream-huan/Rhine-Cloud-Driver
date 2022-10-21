package main

import (
	"io/ioutil"

	"github.com/dream-huan/Rhine-Cloud-Driver/Router"
	"github.com/dream-huan/Rhine-Cloud-Driver/config"
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var cf config.Config

func initConfig() {
	configFile, err := ioutil.ReadFile("./conf/Rhine-Cloud-Driver.yaml")
	if err != nil {
		logger.Error("Read yaml file error", zap.Error(err))
	}
	err = yaml.Unmarshal(configFile, &cf)
	if err != nil {
		logger.Error("Unmarshal yaml file error", zap.Error(err))
	}
}

func main() {
	initConfig()
	Router.InitRouter()
}
