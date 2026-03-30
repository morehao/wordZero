// Package api 提供WordZero的HTTP服务请求/响应模型
package api

// GenerateFromContentRequest 通过内容生成文档的HTTP请求体
type GenerateFromContentRequest struct {
	// Filename 输出文件名（可选，默认为document.docx）
	Filename string `json:"filename,omitempty"`
	// Content 文档内容列表
	Content []ContentItemRequest `json:"content"`
}

// ContentItemRequest 文档内容项请求
type ContentItemRequest struct {
	// Type 内容类型（paragraph/heading/page_break）
	Type string `json:"type"`
	// Text 文本内容
	Text string `json:"text,omitempty"`
	// HeadingLevel 标题级别（1-6，仅对heading类型有效）
	HeadingLevel int `json:"heading_level,omitempty"`
	// Bold 是否加粗
	Bold bool `json:"bold,omitempty"`
	// Italic 是否斜体
	Italic bool `json:"italic,omitempty"`
	// Alignment 对齐方式（left/center/right/both）
	Alignment string `json:"alignment,omitempty"`
}

// GenerateFromTemplateRequest 通过模板生成文档的HTTP请求体
type GenerateFromTemplateRequest struct {
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

// GenerateResponse 文档生成成功响应
type GenerateResponse struct {
	// URL 上传到S3后的访问URL
	URL string `json:"url"`
	// Filename 生成的文件名
	Filename string `json:"filename"`
	// Size 文件大小（字节）
	Size int64 `json:"size"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	// Error 错误信息
	Error string `json:"error"`
}
