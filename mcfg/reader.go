package mcfg

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/caarlos0/env/v6"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

const DefaultFilePath = "config.yml"

func NewConfig[T any]() (T, error) {
	i := new(T)
	return *i, ReadConfig(i)
}

func ReadConfig(i any) error {
	var filePath string

	flag.StringVar(&filePath, "cfg", DefaultFilePath, "Path to config file")
	flag.Parse()

	if _, err := os.Stat(filePath); err == nil {
		if err = LoadFromFile(i, filePath); err != nil {
			return fmt.Errorf("load from file: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("get file stat: %w", err)
	}

	if err := LoadFromEnv(i); err != nil {
		return fmt.Errorf("load from env: %w", err)
	}

	return nil
}

func LoadFromEnv(i any) error {
	defaults.SetDefaults(i)

	if err := env.Parse(i); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	return nil
}

func LoadFromFile(i any, filePath string) error {
	defaults.SetDefaults(i)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	ext := path.Ext(filePath)
	switch {
	case ext == ".yaml" || ext == ".yml":
		err = yaml.Unmarshal(data, i)
		if err != nil {
			return fmt.Errorf("unmarshal yaml: %w", err)
		}
	case ext == ".json":
		err = json.Unmarshal(data, i)
		if err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
	default:
		return fmt.Errorf("unknown config format")
	}

	return nil
}
