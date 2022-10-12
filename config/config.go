package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func baseName() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	filename := filepath.Base(executable)
	return strings.TrimSuffix(filename, filepath.Ext(filename)), nil
}

func baseDir() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(executable), nil
}

func absolutExecutableBaseName() (string, error) {
	basename, err := baseName()
	if err != nil {
		return "", err
	}

	basedir, err := baseDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		basedir,
		basename,
	), nil
}

func executableConfigSibling() (string, error) {
	absBaseName, err := absolutExecutableBaseName()
	if err != nil {
		return "", err
	}
	return absBaseName + ".yaml", nil
}

func unmarshalConfigFile(configFile string, config any) (err error) {
	log.Println("try", configFile)
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		log.Printf("error config from %v", err)
		return nil
	}

	log.Println("load config from", configFile)
	bytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("error read file: %v", err)
		return
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		log.Printf("error unmarshal file: %v", err)
		return
	}
	return
}

func LoadConfig() (err error) {
	file, err := executableConfigSibling()
	if err != nil {
		return
	}
	return unmarshalConfigFile(file, &Config)
}

type config struct {
	ApiKey   string   `yaml:"ApiKey"`
	Products []string `yaml:"Products"`
}

var Config config
