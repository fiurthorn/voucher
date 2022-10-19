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

func homeConfigFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	base, err := baseName()
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(
		home,
		"."+base,
		"config.yaml",
	)

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		content, err := yaml.Marshal(&Config)
		if err != nil {
			return "", err
		}
		os.WriteFile(configFile, content, 0644)
	}

	return configFile, nil
}

func marshalConfigFile(configFile string, config any) (err error) {
	configDir := filepath.Dir(configFile)
	if _, err = os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		log.Println("create ", configDir)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return
		}
	}

	log.Println("marshal config")
	content, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("error marshal file: %v", err)
		return
	}

	log.Println("store config to", configFile)
	err = os.WriteFile(configFile, content, 0644)
	if err != nil {
		log.Printf("error read file: %v", err)
		return
	}

	return
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

func configFile() (configFile string, err error) {
	configFile, err = executableConfigSibling()
	if err != nil {
		return
	}

	if _, err = os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		log.Println("not found", configFile)
		configFile, err = homeConfigFile()
		log.Println("try", configFile)
		if err != nil {
			return
		}
	}

	return
}

func StoreConfig() (file string, err error) {
	file, err = configFile()
	if err != nil {
		return "", err
	}

	err = marshalConfigFile(file, &Config)
	return
}

func LoadConfig() (err error) {
	configFile, err := configFile()
	if err != nil {
		return err
	}

	return unmarshalConfigFile(configFile, &Config)
}

type config struct {
	ApiKey   string   `yaml:"ApiKey"`
	Products []string `yaml:"Products"`
}

var Config = config{
	ApiKey:   "xxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	Products: []string{"4711", "0815"},
}
