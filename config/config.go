package config

import (
	"log/slog"
	"os"

	"golang.org/x/net/publicsuffix"
	"gopkg.in/yaml.v3"
)

type Domain struct {
	Name string
	TLD  string
}

func (d *Domain) UnmarshalYAML(value *yaml.Node) error {
	var name string
	if err := value.Decode(&name); err != nil {
		return err
	}
	d.Name = name
	d.TLD, _ = publicsuffix.PublicSuffix(name)
	return nil
}

type Config struct {
	Domains []Domain `yaml:"domains"`
}

func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("finished reading config", "path", configPath, "domain_count", len(cfg.Domains))
	return &cfg, nil
}
