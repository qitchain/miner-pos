package config

import (
	"encoding/base64"
)

type Config struct {
	Path []string `yaml:"path"`
	Rpc  struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Token    string `yaml:"token"`
	}
	Log struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	}
	FarmerKey        map[string]string `yaml:"farmerKey"`
	FarmerPrivateKey []string          `yaml:"farmerPrivateKey"`
}

func (c *Config) GetAuthorizationToken() string {
	if c.Rpc.Token != "" {
		return c.Rpc.Token
	}
	if c.Rpc.Username != "" || c.Rpc.Password != "" {
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.Rpc.Username+":"+c.Rpc.Password))
	}
	return ""
}
