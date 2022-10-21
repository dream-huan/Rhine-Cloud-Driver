package config

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
	Database string `yaml:"database"`
}

type RecaptchaPrivateKeyConfig struct {
	Key string `yaml:"key"`
}

type JwtConfig struct {
	Key string `yaml:"key"`
}
