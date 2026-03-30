// Package api 提供WordZero的HTTP服务配置
package api

import (
	"fmt"

	"github.com/zerx-lab/wordZero/internal/s3"
)

// Config HTTP服务配置
type Config struct {
	// Host 监听地址，默认为0.0.0.0
	Host string `json:"host" yaml:"host"`
	// Port 监听端口，默认为8080
	Port int `json:"port" yaml:"port"`
	// ReadTimeoutSeconds 读取超时（秒），默认30
	ReadTimeoutSeconds int `json:"read_timeout_seconds" yaml:"read_timeout_seconds"`
	// WriteTimeoutSeconds 写入超时（秒），默认60
	WriteTimeoutSeconds int `json:"write_timeout_seconds" yaml:"write_timeout_seconds"`
	// IdleTimeoutSeconds 空闲超时（秒），默认120
	IdleTimeoutSeconds int `json:"idle_timeout_seconds" yaml:"idle_timeout_seconds"`
	// S3Config S3存储配置
	S3Config s3.Config `json:"s3" yaml:"s3"`
	// TemplateDownloadTimeoutSeconds 模板下载超时（秒），默认30
	TemplateDownloadTimeoutSeconds int `json:"template_download_timeout_seconds" yaml:"template_download_timeout_seconds"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:                           "0.0.0.0",
		Port:                           8080,
		ReadTimeoutSeconds:             30,
		WriteTimeoutSeconds:            60,
		IdleTimeoutSeconds:             120,
		TemplateDownloadTimeoutSeconds: 30,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("端口号无效: %d", c.Port)
	}
	if c.S3Config.Bucket == "" {
		return fmt.Errorf("S3存储桶名称不能为空")
	}
	return nil
}
