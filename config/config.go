package config

import (
	"os"
)

type Config struct {
	Server                    ServerConfig              `yaml:"server"`
	RedisManager              RedisConfig               `yaml:"redis"`
	MysqlManager              MysqlConfig               `yaml:"mysql"`
	GoogleRecaptchaPrivateKey RecaptchaPrivateKeyConfig `yaml:"googlerecaptchaprivatekey"`
	JwtKey                    JwtConfig                 `yaml:"jwtkey"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
}

type RedisConfig struct {
	Address  string `yaml:"addr"`
	Password string `yaml:"pwd"`
}

type MysqlConfig struct {
	Address  string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"pwd"`
}

type RecaptchaPrivateKeyConfig struct {
	Key string `yaml:"key"`
}

type JwtConfig struct {
	Key string `yaml:"key"`
}

var pwd, _ = os.Getwd()
var privatekey = ""
var originstorage int64
var jwtkey = ""

func GetPrivateKey() string {
	return privatekey
}

func GetJwtKey() string {
	return jwtkey
}
