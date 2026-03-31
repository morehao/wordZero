// Package wordzero 提供WordZero的SDK客户端，用于以编程方式生成Word文档并上传至S3存储
package wordzero

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	ygconfig "github.com/ygpkg/yg-go/config"
	"github.com/zerx-lab/wordZero/pkg/document"
	"github.com/zerx-lab/wordZero/pkg/generator"
	"github.com/zerx-lab/wordZero/pkg/s3storage"
)

type Client struct {
	s3Uploader *s3storage.Uploader
	httpClient *http.Client
}

type Config struct {
	S3Config    ygconfig.S3StorageConfig `json:"s3" yaml:"s3"`
	HTTPTimeout time.Duration            `json:"http_timeout" yaml:"http_timeout"`
}

func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("SDK配置不能为空")
	}

	uploader, err := s3storage.NewUploader(cfg.S3Config)
	if err != nil {
		return nil, fmt.Errorf("创建S3上传器失败: %w", err)
	}

	// 配置HTTP客户端（用于下载模板）
	timeout := cfg.HTTPTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	httpClient := &http.Client{Timeout: timeout}

	return &Client{
		s3Uploader: uploader,
		httpClient: httpClient,
	}, nil
}

// GenerateFromContent 根据内容生成Word文档并上传至S3
// 返回文档的S3访问URL
func (c *Client) GenerateFromContent(ctx context.Context, req *ContentRequest) (*GenerateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("请求参数验证失败: %w", err)
	}

	// 创建新文档
	doc := document.New()

	// 添加内容
	for _, item := range req.Content {
		if err := addContentItem(doc, item); err != nil {
			return nil, fmt.Errorf("添加内容失败: %w", err)
		}
	}

	// 将文档转换为字节
	docBytes, err := doc.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("生成文档失败: %w", err)
	}

	// 上传到S3
	filename := req.Filename
	if filename == "" {
		filename = "document.docx"
	}
	key := s3storage.GenerateObjectKey(filename)
	url, err := c.s3Uploader.Upload(ctx, key, docBytes, "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		return nil, fmt.Errorf("上传文档失败: %w", err)
	}

	return &GenerateResponse{
		URL:      url,
		Filename: filename,
		Size:     int64(len(docBytes)),
	}, nil
}

// GenerateFromTemplate 根据模板URL生成Word文档并上传至S3
// 返回文档的S3访问URL
func (c *Client) GenerateFromTemplate(ctx context.Context, req *TemplateRequest) (*GenerateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("请求参数验证失败: %w", err)
	}

	// 下载模板文件
	templateBytes, err := c.downloadTemplate(ctx, req.TemplateURL)
	if err != nil {
		return nil, fmt.Errorf("下载模板失败: %w", err)
	}

	// 使用共享生成器从模板字节渲染文档
	tplData := generator.BuildTemplateData(req.Variables)
	doc, err := generator.RenderTemplateFromBytes(templateBytes, tplData)
	if err != nil {
		return nil, fmt.Errorf("渲染模板失败: %w", err)
	}

	// 将文档转换为字节
	docBytes, err := doc.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("生成文档失败: %w", err)
	}

	// 上传到S3
	filename := req.Filename
	if filename == "" {
		filename = "document.docx"
	}
	key := s3storage.GenerateObjectKey(filename)
	url, err := c.s3Uploader.Upload(ctx, key, docBytes, "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		return nil, fmt.Errorf("上传文档失败: %w", err)
	}

	return &GenerateResponse{
		URL:      url,
		Filename: filename,
		Size:     int64(len(docBytes)),
	}, nil
}

// downloadTemplate 通过HTTP下载模板文件
func (c *Client) downloadTemplate(ctx context.Context, templateURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, templateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载模板失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载模板失败，HTTP状态码: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取模板数据失败: %w", err)
	}

	return data, nil
}

// addContentItem 向文档添加内容项
func addContentItem(doc *document.Document, item ContentItem) error {
	switch item.Type {
	case ContentTypeHeading:
		level := item.HeadingLevel
		if level < 1 || level > 6 {
			level = 1
		}
		doc.AddHeadingParagraph(item.Text, level)
	case ContentTypeParagraph:
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
	case ContentTypePageBreak:
		doc.AddPageBreak()
	default:
		// 默认作为普通段落处理
		doc.AddParagraph(item.Text)
	}
	return nil
}
