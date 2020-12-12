package config

import (
	"os"
	"strings"
	"testing"

	"github.com/tutamuniz/yasrp/minihttp/reverseproxy/configtypes"
)

func TestParseConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		content  string
		expected Config
	}{
		{
			desc: "",
			content: `
			{				
				"bind_ip":"127.0.0.1",
				"bind_port":8080,
				"enable_cache":false,
				"cache_engine": "dummy",
				"locations": [
					{
						"path":"/home",
						"target":"http://www.tjrn.jus.br/"
					},
					{
						"path":"/stats",
						"target":"https://www.google.com"
					}
				]
				
			}
			`,
			expected: Config{
				BindIP:      "127.0.0.1",
				BindPort:    8080,
				EnableCache: false,
				CacheEngine: "dummy",
				Locations: []configtypes.Location{

					{
						Path:   "/home",
						Target: "http://www.tjrn.jus.br",
					},
					{
						Path:   "/stats",
						Target: "https://www.google.com",
					},
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			reader := strings.NewReader(tC.content)
			config, err := ParseConfig(reader)

			if err != nil {
				t.Errorf("%s", err.Error())
			}

			if config.BindIP != tC.expected.BindIP {
				t.Errorf("BindIp error expected %s got %s", tC.expected.BindIP, config.BindIP)
			}
			if config.BindPort != tC.expected.BindPort {
				t.Errorf("BindPort error expected %d got %d", tC.expected.BindPort, config.BindPort)
			}
			if len(config.Locations) != len(tC.expected.Locations) {
				t.Errorf("Error parsing Locations expected %d got %d", len(config.Locations), len(tC.expected.Locations))
			}

			if config.EnableCache != tC.expected.EnableCache {
				t.Errorf("Error parsing EnableCache expected %v got %v", config.EnableCache, tC.expected.EnableCache)
			}

			if config.CacheEngine != tC.expected.CacheEngine {
				t.Errorf("Error parsing EnableCache expected %s got %s", config.CacheEngine, tC.expected.CacheEngine)
			}

			if config.Locations[0].Path != tC.expected.Locations[0].Path {
				t.Errorf("Error parsing Locations expected %s got %s", config.Locations[0].Path, tC.expected.Locations[0].Path)
			}
		})
	}
}

func TestParseConfigFromFile(t *testing.T) {
	configFile := "../../config.json"

	f, err := os.Open(configFile)
	if err != nil {
		t.Errorf("Error opening config file.:%s", err.Error())
	}

	config, err := ParseConfig(f)
	if err != nil {
		t.Errorf("Error parsing config file.:%s", err.Error())
	}
	if config.BindIP != "127.0.0.1" && len(config.Locations) != 2 {
		t.Errorf("Error parsing config file.:%s", err.Error())
	}

	config, err = ParseConfigFromFile(configFile)
	if err != nil {
		t.Errorf("Error parsing config file.:%s", err.Error())
	}
}
