package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/vinted/domain-exporter/config"
	"github.com/vinted/domain-exporter/internal/collector"
)

func main() {
	var configPath, httpListenAddress string

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	flag.StringVar(&configPath, "config_path", "/etc/domain-exporter/config.yaml", "Path to config file.")
	flag.StringVar(&httpListenAddress, "http_listen_address", "0.0.0.0:9553", "Address to bind to.")
	flag.Parse()

	config, err := config.Load(configPath)
	if err != nil {
		slog.Error("failed to read configuration", "error", err, "path", configPath)
		os.Exit(1)
	}

	if err := collector.Start(httpListenAddress, config.Domains); err != nil {
		slog.Error("failed to start collector", "error", err)
		os.Exit(1)
	}
}
