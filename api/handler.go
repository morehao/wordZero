// Package api 提供WordZero的HTTP服务请求处理器
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zerx-lab/wordZero/internal/generator"
	"github.com/zerx-lab/wordZero/internal/s3"
	"github.com/zerx-lab/wordZero/pkg/document"
)

// handler HTTP请求处理器
type handler struct {
	s3Uploader *s3.Uploader
	httpClient *http.Client
}

// newHandler 创建新的请求处理器
func newHandler(cfg *Config) (*handler, error) {
	uploader, err := s3.NewUploader(&cfg.S3Config)
	if err != nil {
		return nil, fmt.Errorf("创建S3上传器失败: %w", err)
	}

	timeout := time.Duration(cfg.TemplateDownloadTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &handler{
		s3Uploader: uploader,
		httpClient: &http.Client{Timeout: timeout},
	}, nil
}

// handleGenerateFromContent 处理通过内容生成文档的请求
// POST /api/v1/documents/content
func (h *handler) handleGenerateFromContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	// 解析请求体
	var req GenerateFromContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("请求体解析失败: %v", err))
		return
	}
	defer r.Body.Close()

	// 验证请求
	if len(req.Content) == 0 {
		writeError(w, http.StatusBadRequest, "内容列表不能为空")
		return
	}

	// 创建新文档
	doc := document.New()

	// 添加内容
	for _, item := range req.Content {
		if err := addDocumentContent(doc, item); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("添加内容失败: %v", err))
			return
		}
	}

	// 将文档转换为字节
	docBytes, err := doc.ToBytes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("生成文档失败: %v", err))
		return
	}

	// 上传到S3
	filename := req.Filename
	if filename == "" {
		filename = "document.docx"
	}
	key := s3.GenerateObjectKey(filename)
	url, err := h.s3Uploader.Upload(r.Context(), key, docBytes,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("上传文档失败: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, GenerateResponse{
		URL:      url,
		Filename: filename,
		Size:     int64(len(docBytes)),
	})
}

// handleGenerateFromTemplate 处理通过模板生成文档的请求
// POST /api/v1/documents/template
func (h *handler) handleGenerateFromTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	// 解析请求体
	var req GenerateFromTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("请求体解析失败: %v", err))
		return
	}
	defer r.Body.Close()

	// 验证请求
	if req.TemplateURL == "" {
		writeError(w, http.StatusBadRequest, "模板URL不能为空")
		return
	}

	// 下载模板文件
	templateBytes, err := h.downloadTemplate(r.Context(), req.TemplateURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("下载模板失败: %v", err))
		return
	}

	// 使用共享生成器从模板字节渲染文档
	tplData := generator.BuildTemplateData(req.Variables)
	doc, err := generator.RenderTemplateFromBytes(templateBytes, tplData)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("渲染模板失败: %v", err))
		return
	}

	// 将文档转换为字节
	docBytes, err := doc.ToBytes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("生成文档失败: %v", err))
		return
	}

	// 上传到S3
	filename := req.Filename
	if filename == "" {
		filename = "document.docx"
	}
	key := s3.GenerateObjectKey(filename)
	url, err := h.s3Uploader.Upload(r.Context(), key, docBytes,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("上传文档失败: %v", err))
		return
	}

	writeJSON(w, http.StatusOK, GenerateResponse{
		URL:      url,
		Filename: filename,
		Size:     int64(len(docBytes)),
	})
}

// downloadTemplate 通过HTTP下载模板文件
func (h *handler) downloadTemplate(ctx context.Context, templateURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, templateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器返回错误状态码: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return data, nil
}

// addDocumentContent 向文档添加内容项
func addDocumentContent(doc *document.Document, item ContentItemRequest) error {
	switch item.Type {
	case "heading":
		level := item.HeadingLevel
		if level < 1 || level > 6 {
			level = 1
		}
		doc.AddHeadingParagraph(item.Text, level)
	case "paragraph", "":
		para := doc.AddParagraph(item.Text)
		if item.Bold {
			para.SetBold(true)
		}
		if item.Italic {
			para.SetItalic(true)
		}
		if item.Alignment != "" {
			para.SetAlignment(document.AlignmentType(item.Alignment))
		}
	case "page_break":
		doc.AddPageBreak()
	default:
		return fmt.Errorf("未知内容类型: %s", item.Type)
	}
	return nil
}
