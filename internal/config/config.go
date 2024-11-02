package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env     string        `yaml:"env" env-required:"true"`
	Graphql GraphqlConfig `yaml:"graphql"`
	Clients ClientsConfig `yaml:"clients"`
}

type GraphqlConfig struct {
	Port int `yaml:"port" env-required:"true"`
}

type Client struct {
	Address      string        `yaml:"address" env-required:"true"`
	Timeout      time.Duration `yaml:"timeout" env-required:"true"`
	RetriesCount int           `yaml:"retries_count" env-required:"true"`
	// Insecure     bool   `yaml:"insecure" env-required:"true"`
}

type ClientsConfig struct {
	Task Client `yaml:"task"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var response string

	// --config="path/to/config.yaml"
	flag.StringVar(&response, "config", "", "path to config file")
	flag.Parse()

	if response == "" {
		response = os.Getenv("CONFIG_PATH")
	}

	return response
}
