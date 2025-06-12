package config

import (
	"flag"
	"os"
	"time"

	"github.com/numbergroup/cleanenv"
)

type Config struct {
	GRPCConfig `yaml:"grpc"`
	HTTPConfig `yaml:"http"`
}

type GRPCConfig struct {
	AuthConfig `yaml:"auth"`
	TaskConfig `yaml:"task"`
}

type AuthConfig struct {
	Port        int           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	Env         string        `yaml:"env"`
	StoragePath string        `yaml:"storage-path"`
	SecretKey   string        `yaml:"secret-key"`
	TokenTTL    time.Duration `yaml:"token-ttl"`
}

type TaskConfig struct {
	Port        int           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	Env         string        `yaml:"env"`
	StoragePath string        `yaml:"storage-path"`
}

type HTTPConfig struct {
	APIGatewayConfig `yaml:"gateway"`
}

type APIGatewayConfig struct {
	Port      int           `yaml:"port"`
	Timeout   time.Duration `yaml:"timeout"`
	Env       string        `yaml:"env"`
	SecretKey string        `yaml:"secret-key"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesnt exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG")
	}

	return res
}
