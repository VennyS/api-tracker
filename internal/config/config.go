package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env:"ENV" env-default:"local"`
	ServerPort  string           `yaml:"server_port" env:"SERVER_PORT" env-default:"8080"`
	MetricsPort string           `yaml:"metrics_port" env:"METRICS_PORT" env-default:"9090"`
	ClickHouse  ClickHouseConfig `yaml:"clickhouse"`
	Prometheus  PrometheusConfig `yaml:"prometheus"`
}

type ClickHouseConfig struct {
	Addr string `yaml:"addr" env:"CLICKHOUSE_ADDR" env-default:"localhost:9000"`
	DB   string `yaml:"db" env:"CLICKHOUSE_DB" env-default:"default"`
}

type PrometheusConfig struct {
	Port string `yaml:"port" env:"PROMETHEUS_PORT" env-default:"9091"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
