// Package wordzero_test 提供wordzero包的单元测试
package wordzero_test

import (
	"testing"

	"github.com/zerx-lab/wordZero/wordzero"
)

// TestContentRequestValidation 测试内容请求验证
func TestContentRequestValidation(t *testing.T) {
	t.Run("空内容列表验证失败", func(t *testing.T) {
		req := &wordzero.ContentRequest{
			Content: []wordzero.ContentItem{},
		}
		if err := req.Validate(); err == nil {
			t.Fatal("期望验证失败，但验证通过了")
		}
	})

	t.Run("非空内容列表验证通过", func(t *testing.T) {
		req := &wordzero.ContentRequest{
			Content: []wordzero.ContentItem{
				{Type: wordzero.ContentTypeParagraph, Text: "测试段落"},
			},
		}
		if err := req.Validate(); err != nil {
			t.Fatalf("期望验证通过，但验证失败: %v", err)
		}
	})
}

// TestTemplateRequestValidation 测试模板请求验证
func TestTemplateRequestValidation(t *testing.T) {
	t.Run("空模板URL验证失败", func(t *testing.T) {
		req := &wordzero.TemplateRequest{
			TemplateURL: "",
		}
		if err := req.Validate(); err == nil {
			t.Fatal("期望验证失败，但验证通过了")
		}
	})

	t.Run("非空模板URL验证通过", func(t *testing.T) {
		req := &wordzero.TemplateRequest{
			TemplateURL: "http://example.com/template.docx",
		}
		if err := req.Validate(); err != nil {
			t.Fatalf("期望验证通过，但验证失败: %v", err)
		}
	})
}

// TestContentTypes 测试内容类型常量
func TestContentTypes(t *testing.T) {
	if wordzero.ContentTypeParagraph != "paragraph" {
		t.Errorf("ContentTypeParagraph 值错误: %s", wordzero.ContentTypeParagraph)
	}
	if wordzero.ContentTypeHeading != "heading" {
		t.Errorf("ContentTypeHeading 值错误: %s", wordzero.ContentTypeHeading)
	}
	if wordzero.ContentTypePageBreak != "page_break" {
		t.Errorf("ContentTypePageBreak 值错误: %s", wordzero.ContentTypePageBreak)
	}
}
