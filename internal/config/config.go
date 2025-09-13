package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type BankConfig struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type TLSConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type Config struct {
	Banks  []BankConfig `yaml:"banks"`
	Server ServerConfig `yaml:"server"`
	TLS    TLSConfig    `yaml:"tls"`
}

var Global *Config

func LoadConfig(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)

	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)

	}

	Global = &cfg

	log.Printf("Loaded %d banks from config", len(cfg.Banks))

}
