package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	configNamespace = "kp2"
	configFile      = "config.yml"
)

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	Databases []DatabaseConfig `yaml:"databases"`
}

func configFileLocation() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, configNamespace, configFile), nil
}

func loadConfig(path string) (Config, error) {
	// Initialize defaults
	conf := Config{}

	// Ignore file read errors and fall back to defaults
	data, err := os.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal([]byte(data), &conf)
	return conf, err
}

func LoadConfig() (Config, error) {
	configPath, err := configFileLocation()
	if err != nil {
		return Config{}, err
	}

	return loadConfig(configPath)
}

func WriteConfig(conf Config) error {
	configPath, err := configFileLocation()
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer configFile.Close()

	yaml, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	_, err = configFile.Write(yaml)

	return err
}
