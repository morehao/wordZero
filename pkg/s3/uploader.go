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
	"github.com/ygpkg/yg-go/logs"
	ygstorage "github.com/ygpkg/yg-go/storage"
)

func init() {
	_ = mime.AddExtensionType(".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
}

type Uploader struct {
	cfg   ygconfig.S3StorageConfig
	once  sync.Once
	fs    *ygstorage.S3Fs
	fsErr error
}

var GlobalUploader *Uploader

func SetGlobalUploader(uploader *Uploader) {
	GlobalUploader = uploader
}

func InitGlobalUploader(cfg ygconfig.S3StorageConfig) error {
	if cfg.Bucket == "" {
		logs.Warnf("[s3] S3 config not found, skip uploader init")
		return nil
	}

	uploader, err := NewUploader(cfg)
	if err != nil {
		logs.Errorf("[s3] init S3 uploader failed: %s", err)
		return err
	}
	SetGlobalUploader(uploader)
	logs.Infof("[s3] S3 uploader initialized, bucket: %s", cfg.Bucket)
	return nil
}

func NewUploader(cfg ygconfig.S3StorageConfig) (*Uploader, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("s3存储桶名称不能为空")
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	return &Uploader{cfg: cfg}, nil
}

func (u *Uploader) init() error {
	u.once.Do(func() {
		opt := ygconfig.StorageOption{}
		u.fs, u.fsErr = ygstorage.NewS3Fs(u.cfg, opt)
	})
	return u.fsErr
}

func (u *Uploader) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	if err := u.init(); err != nil {
		return "", fmt.Errorf("初始化S3连接失败: %w", err)
	}

	ext := filepath.Ext(key)

	fi := &ygstorage.FileInfo{
		StoragePath: key,
		FileExt:     ext,
	}

	if err := u.fs.Save(ctx, fi, bytes.NewReader(data)); err != nil {
		return "", fmt.Errorf("上传到S3失败: %w", err)
	}

	return fi.PublicURL, nil
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
