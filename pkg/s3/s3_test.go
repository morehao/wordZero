// Package s3_test 提供S3包的单元测试
package s3_test

import (
	"strings"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/s3"
)

// TestGenerateObjectKey 测试对象键生成
func TestGenerateObjectKey(t *testing.T) {
	key := s3.GenerateObjectKey("test.docx")
	if key == "" {
		t.Fatal("生成的对象键不能为空")
	}
	// 检查格式是否正确
	if !strings.HasPrefix(key, "documents/") {
		t.Errorf("对象键应该以 'documents/' 开头，实际为: %s", key)
	}
	if !strings.HasSuffix(key, "test.docx") {
		t.Errorf("对象键应该以文件名结尾，实际为: %s", key)
	}
}

// TestGenerateObjectKeyEmpty 测试空文件名的对象键生成
func TestGenerateObjectKeyEmpty(t *testing.T) {
	key := s3.GenerateObjectKey("")
	if key == "" {
		t.Fatal("生成的对象键不能为空")
	}
	if !strings.HasPrefix(key, "documents/") {
		t.Errorf("对象键应该以 'documents/' 开头，实际为: %s", key)
	}
}

// TestNewUploaderValidation 测试上传器配置验证
func TestNewUploaderValidation(t *testing.T) {
	t.Run("空配置应该失败", func(t *testing.T) {
		_, err := s3.NewUploader(nil)
		if err == nil {
			t.Fatal("期望错误，但成功创建了上传器")
		}
	})

	t.Run("空存储桶应该失败", func(t *testing.T) {
		cfg := &s3.Config{
			Region:          "us-east-1",
			Bucket:          "",
			AccessKeyID:     "test",
			SecretAccessKey: "test",
		}
		_, err := s3.NewUploader(cfg)
		if err == nil {
			t.Fatal("期望空存储桶导致错误，但成功创建了上传器")
		}
	})

	t.Run("有效配置应该成功", func(t *testing.T) {
		cfg := &s3.Config{
			Endpoint:        "http://localhost:9000",
			Region:          "us-east-1",
			Bucket:          "test-bucket",
			AccessKeyID:     "test-key",
			SecretAccessKey: "test-secret",
			UsePathStyle:    true,
		}
		uploader, err := s3.NewUploader(cfg)
		if err != nil {
			t.Fatalf("期望成功创建上传器，但失败: %v", err)
		}
		if uploader == nil {
			t.Fatal("上传器不应该为nil")
		}
	})
}
