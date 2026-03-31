package docdto

import (
	"github.com/ygpkg/yg-go/apis/apiobj"
	"github.com/ygpkg/yg-go/apis/errcode"
)

type GenerateFromContentRequest struct {
	apiobj.BaseRequest
	Request GenerateFromContentEmbedRequest `json:"request"`
}

type GenerateFromContentEmbedRequest struct {
	Filename string               `json:"filename,omitempty"`
	Content  []ContentItemRequest `json:"content"`
}

type ContentItemRequest struct {
	Type         string `json:"type"`
	Text         string `json:"text,omitempty"`
	HeadingLevel int    `json:"heading_level,omitempty"`
	Bold         bool   `json:"bold,omitempty"`
	Italic       bool   `json:"italic,omitempty"`
	Alignment    string `json:"alignment,omitempty"`
}

func (opt *GenerateFromContentRequest) Validity(resp *GenerateFromContentResponse) {
	if len(opt.Request.Content) == 0 {
		resp.Code = errcode.ErrCode_BadRequest
		resp.Message = "内容列表不能为空"
		return
	}
}

type GenerateFromTemplateRequest struct {
	apiobj.BaseRequest
	Request GenerateFromTemplateEmbedRequest `json:"request"`
}

type GenerateFromTemplateEmbedRequest struct {
	Filename    string                 `json:"filename,omitempty"`
	TemplateURL string                 `json:"template_url"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
}

func (opt *GenerateFromTemplateRequest) Validity(resp *GenerateFromTemplateResponse) {
	if opt.Request.TemplateURL == "" {
		resp.Code = errcode.ErrCode_BadRequest
		resp.Message = "模板URL不能为空"
		return
	}
}
