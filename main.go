package main

import (
	"io/ioutil"

	"github.com/dream-huan/Rhine-Cloud-Driver/Router"
	"github.com/dream-huan/Rhine-Cloud-Driver/config"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Jwt"
	logger "github.com/dream-huan/Rhine-Cloud-Driver/middleware/Log"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Mysql"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Recaptcha"
	"github.com/dream-huan/Rhine-Cloud-Driver/middleware/Redis"
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

func initService() {
	Jwt.Init(cf.JwtKey.Key)
	Recaptcha.Init(cf.GoogleRecaptchaPrivateKey.Key)
	Redis.Init(cf)
	Mysql.Init(cf)
}

func main() {
	initConfig()
	initService()
	Router.InitRouter(cf)
}
