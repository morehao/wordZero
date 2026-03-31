// Package s3storage_test 提供S3storage包的单元测试
package s3storage_test

import (
	"strings"
	"testing"

	ygconfig "github.com/ygpkg/yg-go/config"
	"github.com/zerx-lab/wordZero/pkg/s3storage"
)

func TestGenerateObjectKey(t *testing.T) {
	key := s3storage.GenerateObjectKey("test.docx")
	if key == "" {
		t.Fatal("生成的对象键不能为空")
	}
	if !strings.HasPrefix(key, "documents/") {
		t.Errorf("对象键应该以 'documents/' 开头，实际为: %s", key)
	}
	if !strings.HasSuffix(key, "test.docx") {
		t.Errorf("对象键应该以文件名结尾，实际为: %s", key)
	}
}

func TestGenerateObjectKeyEmpty(t *testing.T) {
	key := s3storage.GenerateObjectKey("")
	if key == "" {
		t.Fatal("生成的对象键不能为空")
	}
	if !strings.HasPrefix(key, "documents/") {
		t.Errorf("对象键应该以 'documents/' 开头，实际为: %s", key)
	}
}

func TestNewUploaderValidation(t *testing.T) {
	t.Run("空存储桶应该失败", func(t *testing.T) {
		cfg := ygconfig.S3StorageConfig{
			Region:          "us-east-1",
			Bucket:          "",
			AccessKeyID:     "test",
			SecretAccessKey: "test",
		}
		_, err := s3storage.NewUploader(cfg)
		if err == nil {
			t.Fatal("期望空存储桶导致错误，但成功创建了上传器")
		}
	})

	t.Run("有效配置应该成功", func(t *testing.T) {
		cfg := ygconfig.S3StorageConfig{
			EndPoint:        "http://localhost:9000",
			Region:          "us-east-1",
			Bucket:          "test-bucket",
			AccessKeyID:     "test-key",
			SecretAccessKey: "test-secret",
			UsePathStyle:    true,
		}
		uploader, err := s3storage.NewUploader(cfg)
		if err != nil {
			t.Fatalf("期望成功创建上传器，但失败: %v", err)
		}
		if uploader == nil {
			t.Fatal("上传器不应该为nil")
		}
	})
}
