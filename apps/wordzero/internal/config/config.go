package config

import (
	"github.com/ygpkg/yg-go/config"
	"github.com/zerx-lab/wordZero/pkg/s3"
)

type WordZeroConfig struct {
	config.CoreConfig `yaml:",inline"`
	Storage           StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	S3            s3.Config `yaml:"s3"`
	KeyPrefix     string    `yaml:"key_prefix"`
	PublicBaseURL string    `yaml:"public_base_url"`
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

func (c *WordZeroConfig) GetS3Config() *s3.Config {
	return &s3.Config{
		Endpoint:        c.Storage.S3.Endpoint,
		Region:          c.Storage.S3.Region,
		Bucket:          c.Storage.S3.Bucket,
		AccessKeyID:     c.Storage.S3.AccessKeyID,
		SecretAccessKey: c.Storage.S3.SecretAccessKey,
		KeyPrefix:       c.Storage.KeyPrefix,
		PublicBaseURL:   c.Storage.PublicBaseURL,
		UsePathStyle:    c.Storage.S3.UsePathStyle,
	}
}
