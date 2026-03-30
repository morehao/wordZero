// Package sdk 提供WordZero的SDK客户端模型定义
package sdk

import "fmt"

// ContentType 内容类型
type ContentType string

const (
	// ContentTypeParagraph 普通段落
	ContentTypeParagraph ContentType = "paragraph"
	// ContentTypeHeading 标题
	ContentTypeHeading ContentType = "heading"
	// ContentTypePageBreak 分页符
	ContentTypePageBreak ContentType = "page_break"
)

// ContentItem 文档内容项
type ContentItem struct {
	// Type 内容类型
	Type ContentType `json:"type"`
	// Text 文本内容
	Text string `json:"text,omitempty"`
	// HeadingLevel 标题级别（1-6，仅对ContentTypeHeading有效）
	HeadingLevel int `json:"heading_level,omitempty"`
	// Bold 是否加粗
	Bold bool `json:"bold,omitempty"`
	// Italic 是否斜体
	Italic bool `json:"italic,omitempty"`
	// Alignment 对齐方式（left/center/right/both）
	Alignment string `json:"alignment,omitempty"`
}

// ContentRequest 通过内容生成文档的请求
type ContentRequest struct {
	// Filename 输出文件名（可选，默认为document.docx）
	Filename string `json:"filename,omitempty"`
	// Content 文档内容列表
	Content []ContentItem `json:"content"`
}

// Validate 验证请求参数
func (r *ContentRequest) Validate() error {
	if len(r.Content) == 0 {
		return fmt.Errorf("内容列表不能为空")
	}
	return nil
}

// TemplateRequest 通过模板生成文档的请求
type TemplateRequest struct {
	// Filename 输出文件名（可选，默认为document.docx）
	Filename string `json:"filename,omitempty"`
	// TemplateURL 模板文件的URL（支持HTTP/HTTPS）
	TemplateURL string `json:"template_url"`
	// Variables 模板变量
	// - bool类型变量将作为条件变量
	// - []interface{}类型变量将作为列表变量
	// - 其他类型将作为普通变量
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// Validate 验证请求参数
func (r *TemplateRequest) Validate() error {
	if r.TemplateURL == "" {
		return fmt.Errorf("模板URL不能为空")
	}
	return nil
}

// GenerateResponse 文档生成响应
type GenerateResponse struct {
	// URL 上传到S3后的访问URL
	URL string `json:"url"`
	// Filename 生成的文件名
	Filename string `json:"filename"`
	// Size 文件大小（字节）
	Size int64 `json:"size"`
}
