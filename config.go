package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BearerToken     string `yaml:"bearerToken" json:"bearerToken"`
	Domain          string `yaml:"domain" json:"domain"`
	RecordNameToset string `yaml:"recordNameToSet" json:"recordNameToSet"`
	ZoneID          string `yaml:"zoneId" json:"zoneId"`
	DnsID           string `yaml:"dnsId" json:"dnsId"`
}

// purposely does not use ioutil
func readConfig(configFile string) ([]byte, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := make([]byte, 4096)

	i, err := f.Read(b)
	if err != nil {
		return nil, err
	}

	return b[:i], nil
}

func parseConfig(config []byte) (*Config, error) {
	var conf Config
	if err := yaml.Unmarshal(config, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
