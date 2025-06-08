package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env:"ENV" env-default:"local"`
	ServerPort  string           `yaml:"server_port" env:"API_TRACKER_PORT" env-default:"8080"`
	MetricsPort string           `yaml:"metrics_port" env:"METRICS_PORT" env-default:"2112"`
	LogLevel    string           `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	ClickHouse  ClickHouseConfig `yaml:"clickhouse"`
	Prometheus  PrometheusConfig `yaml:"prometheus"`
	Grafana     GrafanaConfig    `yaml:"grafana"`
	HTTP        HTTPConfig       `yaml:"http"`
	Cache       CacheConfig      `yaml:"cache"`
}

type ClickHouseConfig struct {
	Addr            string        `yaml:"addr" env:"CLICKHOUSE_ADDR" env-default:"localhost"`
	PortHTTP        int           `yaml:"port_http" env:"CLICKHOUSE_PORT_HTTP" env-default:"8123"`
	PortNative      int           `yaml:"port_native" env:"CLICKHOUSE_PORT_NATIVE" env-default:"9000"`
	DB              string        `yaml:"db" env:"CLICKHOUSE_DB" env-default:"default"`
	User            string        `yaml:"user" env:"CLICKHOUSE_USER" env-default:"default"`
	Password        string        `yaml:"password" env:"CLICKHOUSE_PASSWORD" env-default:""`
	Secure          bool          `yaml:"secure" env:"CLICKHOUSE_SECURE" env-default:"false"`
	Debug           bool          `yaml:"debug" env:"CLICKHOUSE_DEBUG" env-default:"false"`
	Cluster         string        `yaml:"cluster" env:"CLICKHOUSE_CLUSTER" env-default:""`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"CLICKHOUSE_MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"CLICKHOUSE_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"CLICKHOUSE_CONN_MAX_LIFETIME" env-default:"5m"`
	MigrationsDir   string        `yaml:"migrations_dir" env:"CLICKHOUSE_MIGRATIONS_DIR" env-default:"./migrations"`
	MigrationsTable string        `yaml:"migrations_table" env:"CLICKHOUSE_MIGRATIONS_TABLE" env-default:"schema_migrations"`
}

type PrometheusConfig struct {
	Port        string `yaml:"port" env:"PROMETHEUS_PORT" env-default:"9090"`
	MetricsPath string `yaml:"metrics_path" env:"PROMETHEUS_METRICS_PATH" env-default:"/metrics"`
}

type GrafanaConfig struct {
	Port          string `yaml:"port" env:"GRAFANA_PORT" env-default:"3000"`
	User          string `yaml:"user" env:"GRAFANA_USER" env-default:"admin"`
	Password      string `yaml:"password" env:"GRAFANA_PASS" env-default:"admin"`
	DashboardsDir string `yaml:"dashboards_dir" env:"GRAFANA_DASHBOARDS_DIR" env-default:"./grafana/dashboards"`
}

type HTTPConfig struct {
	ServerTimeout int `yaml:"server_timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"30"`
	ClientTimeout int `yaml:"client_timeout" env:"HTTP_CLIENT_TIMEOUT" env-default:"10"`
}

type CacheConfig struct {
	Enabled bool `yaml:"enabled" env:"CACHE_ENABLED" env-default:"true"`
	TTL     int  `yaml:"ttl" env:"CACHE_TTL" env-default:"300"`
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
