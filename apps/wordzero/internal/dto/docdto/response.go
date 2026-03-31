package docdto

import (
	"github.com/ygpkg/yg-go/apis/apiobj"
)

type GenerateFromContentResponse struct {
	apiobj.BaseResponse
	Response GenerateFromContentEmbedResponse `json:"response"`
}

type GenerateFromContentEmbedResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type GenerateFromTemplateResponse struct {
	apiobj.BaseResponse
	Response GenerateFromTemplateEmbedResponse `json:"response"`
}

type GenerateFromTemplateEmbedResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}
