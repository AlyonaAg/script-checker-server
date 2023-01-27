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
	BindAddr            = "server.bind_addr"
	DatabaseURL         = "store.database_url"
	OriginalScriptTopic = "kafka.original_script_topic"
	DeobfScriptTopic    = "kafka.deobf_script_topic"
	Brokers             = "kafka.brokers"
	RetryMax            = "kafka.retry_max"
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
