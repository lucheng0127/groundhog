package common

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

//ServerConfig configure of server
type ServerConfig struct {
	Port        int    `yaml:"port"`
	IPAddr      string `yaml:"ipAddr"`
	PubnetIface string `yaml:"iface"`
	VPNEnable   bool   `yaml:"vpnEnable"`
	MTU         int    `yaml:"mtu"`
}

// ClientConfig configure of client
type ClientConfig struct {
	ServerAddr string `yaml:"serverAddr"`
	SecretKey  string `yaml:"key"`
	LocalAddr  string `yaml:"localAddr"`
}

// GroundhogConfig configure of groundhog
type GroundhogConfig struct {
	Mode   string       `yaml:"mode"`
	Server ServerConfig `yaml:"server"`
	Client ClientConfig `yaml:"client"`
}

// ParseConfig read configure from yaml file and return config
func ParseConfig(filename string) (interface{}, error) {
	cfg := new(GroundhogConfig)
	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, InitailizeError(err.Error())
	}

	err = yaml.Unmarshal(cfgFile, &cfg)
	if err != nil {
		return nil, InitailizeError(err.Error())
	}

	switch cfg.Mode {
	case "server":
		return cfg.Server, nil
	case "client":
		return cfg.Client, nil
	default:
		return nil, InitailizeError("Unsupport launch mode, server or client.")
	}
}
