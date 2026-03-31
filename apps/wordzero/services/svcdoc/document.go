package svcdoc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zerx-lab/wordZero/apps/wordzero/internal/dto/docdto"
	"github.com/zerx-lab/wordZero/pkg/document"
	"github.com/zerx-lab/wordZero/pkg/generator"
	"github.com/zerx-lab/wordZero/pkg/s3"
)

type documentService struct {
	s3Uploader *s3.Uploader
	httpClient *http.Client
}

func NewDocumentService(s3Cfg *s3.Config) (*documentService, error) {
	uploader, err := s3.NewUploader(s3Cfg)
	if err != nil {
		return nil, fmt.Errorf("创建S3上传器失败: %w", err)
	}

	return &documentService{
		s3Uploader: uploader,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func GenerateFromContent(ctx *gin.Context, req *docdto.GenerateFromContentRequest) (res *docdto.GenerateFromContentResponse, err error) {
	res = &docdto.GenerateFromContentResponse{}

	doc := document.New()

	for _, item := range req.Request.Content {
		if err := addDocumentContent(doc, item); err != nil {
			return nil, fmt.Errorf("添加内容失败: %w", err)
		}
	}

	docBytes, err := doc.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("生成文档失败: %w", err)
	}

	filename := req.Request.Filename
	if filename == "" {
		filename = "document.docx"
	}

	key := s3.GenerateObjectKey(filename)
	url, err := s3UploaderGlobal.Upload(ctx.Request.Context(), key, docBytes,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		return nil, fmt.Errorf("上传文档失败: %w", err)
	}

	res.Response.URL = url
	res.Response.Filename = filename
	res.Response.Size = int64(len(docBytes))
	return res, nil
}

func GenerateFromTemplate(ctx *gin.Context, req *docdto.GenerateFromTemplateRequest) (res *docdto.GenerateFromTemplateResponse, err error) {
	res = &docdto.GenerateFromTemplateResponse{}

	templateBytes, err := downloadTemplate(ctx.Request.Context(), req.Request.TemplateURL)
	if err != nil {
		return nil, fmt.Errorf("下载模板失败: %w", err)
	}

	tplData := generator.BuildTemplateData(req.Request.Variables)
	doc, err := generator.RenderTemplateFromBytes(templateBytes, tplData)
	if err != nil {
		return nil, fmt.Errorf("渲染模板失败: %w", err)
	}

	docBytes, err := doc.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("生成文档失败: %w", err)
	}

	filename := req.Request.Filename
	if filename == "" {
		filename = "document.docx"
	}

	key := s3.GenerateObjectKey(filename)
	url, err := s3UploaderGlobal.Upload(ctx.Request.Context(), key, docBytes,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	if err != nil {
		return nil, fmt.Errorf("上传文档失败: %w", err)
	}

	res.Response.URL = url
	res.Response.Filename = filename
	res.Response.Size = int64(len(docBytes))
	return res, nil
}

var s3UploaderGlobal *s3.Uploader

func SetGlobalUploader(uploader *s3.Uploader) {
	s3UploaderGlobal = uploader
}

func downloadTemplate(ctx context.Context, templateURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, templateURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
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

func addDocumentContent(doc *document.Document, item docdto.ContentItemRequest) error {
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
