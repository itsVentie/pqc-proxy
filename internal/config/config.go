package config

import (
	"flag"
	"os"
)

type Config struct {
	Mode        string
	ListenAddr  string
	TargetAddr  string
	MetricsAddr string
	Debug       bool
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Mode, "mode", getEnv("PQC_MODE", ""), "Running mode: client or server")
	flag.StringVar(&cfg.ListenAddr, "listen", getEnv("PQC_LISTEN", ""), "Address to listen on")
	flag.StringVar(&cfg.TargetAddr, "target", getEnv("PQC_TARGET", ""), "Target remote address")
	flag.StringVar(&cfg.MetricsAddr, "metrics", getEnv("PQC_METRICS", ":2112"), "Address for Prometheus metrics")
	flag.BoolVar(&cfg.Debug, "debug", false, "Enable debug logging")

	flag.Parse()

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
