package config

import (
	"flag"
	"os"
	"time"

	"github.com/numbergroup/cleanenv"
)

type Config struct {
	Env             string        `yaml:"env" env-default:"local"`
	AuthStoragePath string        `yaml:"auth-storage-path" env-required:"true"`
	TokenTTL        time.Duration `yaml:"token-ttl" env-reuired:"true"`
	GRPC            GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
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
