package env

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Settings struct {
	Host         string      `yaml:"host"`
	Port         int         `yaml:"port"`
	JwtKey       string      `yaml:"jwt_key"`
	GRPC_PORT    int         `yaml:"grpc_port"`
	Redis        RedisEnv    `yaml:"redis"`
	Postgres     DatabaseEnv `yaml:"postgres"`
	IsMailEnable bool        `yaml:"is_email_enable"`
	Mail         ServerEnv   `yaml:"mail"`
}

type Env struct {
	Settings Settings `yaml:"settings"`
}

func GetEnv() Env {
	path, ok := os.LookupEnv("SETTING_PATH")
	if !ok {
		path = "env.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var env Env
	if err := yaml.Unmarshal(data, &env); err != nil {
		panic(err)
	}
	return env
}
