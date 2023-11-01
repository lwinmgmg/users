package env

type ServerEnv struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

type DatabaseEnv struct {
	// ServerEnv
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

type RedisEnv struct {
	// ServerEnv
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
