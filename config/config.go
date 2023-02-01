package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	JsonIPService string `yaml:"jsonIPService" json:"jsonIPService"`
	JsonQuery     string `yaml:"jsonQuery" json:"jsonQuery"`
	AccessToken   string `yaml:"accessToken" json:"accessToken"`
	Domain        string `yaml:"domain" json:"domain"`
	RecordName    string `yaml:"recordName" json:"recordName"`
	RecordType    string `yaml:"recordType" json:"recordType"`
	ZoneID        string `yaml:"zoneId" json:"zoneId"`
	DnsID         string `yaml:"dnsId" json:"dnsId"`
}

func ReadConfig(configFile string) ([]byte, error) {
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

func ParseConfig(config []byte) (*Config, error) {
	var conf Config
	if err := yaml.Unmarshal(config, &conf); err != nil {
		return nil, err
	}

	if conf.AccessToken == "" {
		return nil, errors.New("could not read bearerToken from config file")
	}
	if conf.Domain == "" {
		return nil, errors.New("could not read domain from config file")
	}
	if conf.RecordName == "" {
		return nil, errors.New("could not read recordName from config file")
	}
	if conf.RecordType == "" {
		return nil, errors.New("could not read recordType from config file")
	}

	return &conf, nil
}
