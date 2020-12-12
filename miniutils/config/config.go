package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/tutamuniz/yasrp/minihttp/reverseproxy/configtypes"
)

// Config information from json file
type Config struct {
	BindIP      string                 `json:"bind_ip"`
	BindPort    uint16                 `json:"bind_port"`
	EnableCache bool                   `json:"enable_cache"`
	Locations   []configtypes.Location `json:"locations"`
	CacheEngine string                 `json:"cache_engine"`
}

// ParseConfig parses the config data
func ParseConfig(reader io.Reader) (*Config, error) {
	content, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf("ParseConfig(): %s", err.Error())
	}

	config := &Config{}

	err = json.Unmarshal(content, config)

	if err != nil {
		return nil, fmt.Errorf("ParseConfig(): %s", err.Error())
	}

	return config, configValidations()
}

// ParseConfigFromFile parse the config data using a file
func ParseConfigFromFile(filename string) (*Config, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file.:%s", err.Error())
	}

	config, err := ParseConfig(f)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config file.:%s", err.Error())
	}
	return config, nil
}

func configValidations() error {
	// Check valid CacheEngine, BindIP, etc.
	return nil
}
