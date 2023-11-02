package env

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	EnvFileName = "env.yaml"
)

type Settings struct {
	Host      string      `yaml:"host"`
	Port      int         `yaml:"port"`
	JwtKey    string      `yaml:"jwt_key"`
	GRPC_PORT int         `yaml:"grpc_port"`
	Redis     RedisEnv    `yaml:"redis"`
	Postgres  DatabaseEnv `yaml:"postgres"`
	Mail      ServerEnv   `yaml:"mail"`
}

type Env struct {
	Settings Settings `yaml:"settings"`
}

func GetEnv() Env {
	data, err := os.ReadFile(EnvFileName)
	if err != nil {
		panic(err)
	}
	var env Env
	if err := yaml.Unmarshal(data, &env); err != nil {
		panic(err)
	}
	fmt.Println(env)
	return env
}
