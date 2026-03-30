// Package s3 提供S3协议兼容的对象存储上传功能
package s3

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config S3上传配置
type Config struct {
	// Endpoint S3兼容存储的端点地址（留空则使用AWS S3）
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	// Region 存储区域
	Region string `json:"region" yaml:"region"`
	// Bucket 存储桶名称
	Bucket string `json:"bucket" yaml:"bucket"`
	// AccessKeyID 访问密钥ID
	AccessKeyID string `json:"access_key_id" yaml:"access_key_id"`
	// SecretAccessKey 访问密钥
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
	// KeyPrefix 对象键前缀（可选）
	KeyPrefix string `json:"key_prefix" yaml:"key_prefix"`
	// UsePathStyle 是否使用路径样式访问（非AWS S3时通常需要设置为true）
	UsePathStyle bool `json:"use_path_style" yaml:"use_path_style"`
	// PublicBaseURL 公开访问基础URL（可选，用于生成访问链接）
	PublicBaseURL string `json:"public_base_url" yaml:"public_base_url"`
}

// Uploader S3上传器
type Uploader struct {
	client *s3.Client
	cfg    *Config
}

// NewUploader 创建新的S3上传器
func NewUploader(cfg *Config) (*Uploader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("s3配置不能为空")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("s3存储桶名称不能为空")
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	// 构建AWS配置选项
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}

	// 配置静态凭证（如果提供）
	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	// 加载AWS配置
	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("加载AWS配置失败: %w", err)
	}

	// 创建S3客户端选项
	s3Opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = cfg.UsePathStyle
		},
	}

	// 配置自定义端点（用于MinIO、阿里云OSS等S3兼容存储）
	if cfg.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	}

	client := s3.NewFromConfig(awsCfg, s3Opts...)

	return &Uploader{
		client: client,
		cfg:    cfg,
	}, nil
}

// Upload 上传文件到S3存储，返回对象的访问URL
func (u *Uploader) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	// 构建完整的对象键
	fullKey := key
	if u.cfg.KeyPrefix != "" {
		fullKey = u.cfg.KeyPrefix + "/" + key
	}

	// 上传对象
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.cfg.Bucket),
		Key:         aws.String(fullKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("上传到S3失败: %w", err)
	}

	// 生成访问URL
	url := u.buildURL(fullKey)
	return url, nil
}

// buildURL 构建对象访问URL
func (u *Uploader) buildURL(key string) string {
	if u.cfg.PublicBaseURL != "" {
		return fmt.Sprintf("%s/%s", trimTrailingSlash(u.cfg.PublicBaseURL), key)
	}

	// 使用自定义端点
	if u.cfg.Endpoint != "" {
		endpoint := trimTrailingSlash(u.cfg.Endpoint)
		if u.cfg.UsePathStyle {
			return fmt.Sprintf("%s/%s/%s", endpoint, u.cfg.Bucket, key)
		}
		return fmt.Sprintf("%s/%s", endpoint, key)
	}

	// 默认AWS S3 URL格式
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.cfg.Bucket, u.cfg.Region, key)
}

// trimTrailingSlash 去除字符串末尾的斜杠
func trimTrailingSlash(s string) string {
	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

// GenerateObjectKey 生成唯一的对象键
// 格式: documents/{YYYY-MM-DD}/{timestamp}_{filename}.docx
func GenerateObjectKey(filename string) string {
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	timestamp := now.UnixNano() / int64(time.Millisecond)

	// 获取不含路径的文件名
	base := filepath.Base(filename)
	if base == "" || base == "." {
		base = "document.docx"
	}

	return fmt.Sprintf("documents/%s/%d_%s", dateStr, timestamp, base)
}
