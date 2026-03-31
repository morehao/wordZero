package conf

import (
	"github.com/ygpkg/yg-go/config"
)

type WordZeroConfig struct {
	config.CoreConfig `yaml:",inline"`
	Storage           StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	config.S3StorageConfig `yaml:",inline"`
}

var std *WordZeroConfig

func Conf() *WordZeroConfig {
	return std
}

func LoadConfig(configFile string) (*WordZeroConfig, error) {
	cfg := &WordZeroConfig{}
	err := config.LoadYamlLocalFile(configFile, cfg)
	if err != nil {
		return nil, err
	}
	std = cfg
	return cfg, nil
}

func (c *WordZeroConfig) GetS3Config() config.S3StorageConfig {
	return c.Storage.S3StorageConfig
}
