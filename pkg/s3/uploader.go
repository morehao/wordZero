// Package s3 提供S3协议兼容的对象存储上传功能，底层使用 github.com/ygpkg/yg-go/storage 实现
package s3

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"sync"
	"time"

	ygconfig "github.com/ygpkg/yg-go/config"
	ygstorage "github.com/ygpkg/yg-go/storage"
)

func init() {
	// 确保 .docx MIME 类型在所有操作系统上正确注册
	_ = mime.AddExtensionType(".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
}

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

// Uploader S3上传器，使用 github.com/ygpkg/yg-go/storage 实现上传操作
type Uploader struct {
	cfg   *Config
	once  sync.Once
	fs    *ygstorage.S3Fs
	fsErr error
}

var GlobalUploader *Uploader

func SetGlobalUploader(uploader *Uploader) {
	GlobalUploader = uploader
}

// NewUploader 创建新的S3上传器（延迟初始化底层存储连接）
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

	return &Uploader{cfg: cfg}, nil
}

// init 延迟初始化底层 yg-go S3Fs 存储器（首次上传时建立连接）
func (u *Uploader) init() error {
	u.once.Do(func() {
		s3cfg := ygconfig.S3StorageConfig{
			EndPoint:        u.cfg.Endpoint,
			Region:          u.cfg.Region,
			Bucket:          u.cfg.Bucket,
			AccessKeyID:     u.cfg.AccessKeyID,
			SecretAccessKey: u.cfg.SecretAccessKey,
			UsePathStyle:    u.cfg.UsePathStyle,
		}
		opt := ygconfig.StorageOption{}
		u.fs, u.fsErr = ygstorage.NewS3Fs(s3cfg, opt)
	})
	return u.fsErr
}

// Upload 上传文件到S3存储，返回对象的访问URL
func (u *Uploader) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	// 延迟初始化底层存储连接
	if err := u.init(); err != nil {
		return "", fmt.Errorf("初始化S3连接失败: %w", err)
	}

	// 构建完整的对象键
	fullKey := key
	if u.cfg.KeyPrefix != "" {
		fullKey = u.cfg.KeyPrefix + "/" + key
	}

	// 获取文件扩展名（用于设置 MIME 类型）
	ext := filepath.Ext(fullKey)

	// 构建 FileInfo
	fi := &ygstorage.FileInfo{
		StoragePath: fullKey,
		FileExt:     ext,
	}

	// 使用 yg-go storage 上传文件
	if err := u.fs.Save(ctx, fi, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("上传到S3失败: %w", err)
	}

	// 如果配置了 PublicBaseURL，则覆盖生成的 URL
	if u.cfg.PublicBaseURL != "" {
		return fmt.Sprintf("%s/%s", trimTrailingSlash(u.cfg.PublicBaseURL), fullKey), nil
	}

	return fi.PublicURL, nil
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
