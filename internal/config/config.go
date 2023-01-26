package config

import (
	"errors"
	"os"

	"github.com/pelletier/go-toml"
)

var (
	noKeyEnvironmentVariables = errors.New("no key in environment variables")
)

const (
	BindAddr = "server.bind_addr"

	DatabaseURL   = "store.database_url"
	PathMigration = "store.path_migration"
)

func GetValue(key string) (interface{}, error) {
	configPath, ok := os.LookupEnv("PATH_CONFIG")
	if !ok {
		return nil, noKeyEnvironmentVariables
	}

	config, err := toml.LoadFile(configPath)
	if err != nil {
		return nil, err
	}

	return config.Get(key), nil
}
