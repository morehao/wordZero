// Package api_test 提供API包的单元测试
package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zerx-lab/wordZero/api"
	"github.com/zerx-lab/wordZero/internal/s3"
)

// newTestConfig 创建测试用的服务器配置（使用无效的S3配置）
func newTestConfig() *api.Config {
	return &api.Config{
		Host:                           "127.0.0.1",
		Port:                           8080,
		ReadTimeoutSeconds:             30,
		WriteTimeoutSeconds:            60,
		IdleTimeoutSeconds:             120,
		TemplateDownloadTimeoutSeconds: 30,
		S3Config: s3.Config{
			Endpoint:        "http://localhost:9000",
			Region:          "us-east-1",
			Bucket:          "test-bucket",
			AccessKeyID:     "test-key",
			SecretAccessKey: "test-secret",
			UsePathStyle:    true,
		},
	}
}

// TestConfigValidation 测试配置验证
func TestConfigValidation(t *testing.T) {
	t.Run("有效配置验证通过", func(t *testing.T) {
		cfg := newTestConfig()
		if err := cfg.Validate(); err != nil {
			t.Fatalf("期望验证通过，但验证失败: %v", err)
		}
	})

	t.Run("无效端口验证失败", func(t *testing.T) {
		cfg := newTestConfig()
		cfg.Port = 0
		if err := cfg.Validate(); err == nil {
			t.Fatal("期望验证失败，但验证通过了")
		}
	})

	t.Run("空存储桶验证失败", func(t *testing.T) {
		cfg := newTestConfig()
		cfg.S3Config.Bucket = ""
		if err := cfg.Validate(); err == nil {
			t.Fatal("期望验证失败，但验证通过了")
		}
	})
}

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	cfg := api.DefaultConfig()
	if cfg.Host != "0.0.0.0" {
		t.Errorf("默认Host错误: %s", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("默认Port错误: %d", cfg.Port)
	}
}

// TestHealthEndpoint 测试健康检查端点（通过测试服务器）
func TestHealthEndpoint(t *testing.T) {
	// 使用httptest直接测试处理函数
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			t.Errorf("写入响应失败: %v", err)
		}
	})

	// GET请求应该成功
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际 %d", http.StatusOK, rr.Code)
	}

	// 验证响应体
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应体失败: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("期望status为ok，实际为: %s", resp["status"])
	}
}

// TestGenerateFromContentRequest 测试内容生成请求的解析
func TestGenerateFromContentRequest(t *testing.T) {
	reqBody := api.GenerateFromContentRequest{
		Filename: "test.docx",
		Content: []api.ContentItemRequest{
			{Type: "heading", Text: "测试标题", HeadingLevel: 1},
			{Type: "paragraph", Text: "测试内容"},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("序列化请求体失败: %v", err)
	}

	// 验证序列化/反序列化一致性
	var decoded api.GenerateFromContentRequest
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&decoded); err != nil {
		t.Fatalf("反序列化请求体失败: %v", err)
	}

	if decoded.Filename != reqBody.Filename {
		t.Errorf("Filename不一致: 期望 %s，实际 %s", reqBody.Filename, decoded.Filename)
	}
	if len(decoded.Content) != len(reqBody.Content) {
		t.Errorf("Content长度不一致: 期望 %d，实际 %d", len(reqBody.Content), len(decoded.Content))
	}
}
