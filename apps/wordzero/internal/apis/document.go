package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/ygpkg/yg-go/apis/errcode"
	"github.com/ygpkg/yg-go/logs"
	"github.com/zerx-lab/wordZero/apps/wordzero/internal/dto/docdto"
	"github.com/zerx-lab/wordZero/apps/wordzero/services/svcdoc"
)

// GenerateFromContent 从内容生成文档
// @Tags 文档生成
// @Summary 从内容生成文档
// @Description 根据提供的内容列表生成 Word 文档
// @Router /wordzero.GenerateFromContent [post]
// @Param request body docdto.GenerateFromContentRequest true "request"
// @Success 200 {object} docdto.GenerateFromContentResponse "response"
func GenerateFromContent(ctx *gin.Context, req *docdto.GenerateFromContentRequest, resp *docdto.GenerateFromContentResponse) {
	if req.Validity(resp); resp.Code != 0 {
		logs.ErrorContextf(ctx, "[GenerateFromContent] request invalid, req: %s, error message: %v", logs.JSON(req), resp.Message)
		return
	}

	res, err := svcdoc.GenerateFromContent(ctx, req)
	if err != nil {
		logs.ErrorContextf(ctx, "[GenerateFromContent] svcdoc.GenerateFromContent failed, err: %v", err)
		resp.Code = errcode.ErrCode_InternalError
		resp.Message = "生成文档失败"
		return
	}
	resp.Code = res.Code
	resp.Message = res.Message
	resp.Response = res.Response
}

// GenerateFromTemplate 从模板生成文档
// @Tags 文档生成
// @Summary 从模板生成文档
// @Description 使用模板和变量生成 Word 文档
// @Router /wordzero.GenerateFromTemplate [post]
// @Param request body docdto.GenerateFromTemplateRequest true "request"
// @Success 200 {object} docdto.GenerateFromTemplateResponse "response"
func GenerateFromTemplate(ctx *gin.Context, req *docdto.GenerateFromTemplateRequest, resp *docdto.GenerateFromTemplateResponse) {
	if req.Validity(resp); resp.Code != 0 {
		logs.ErrorContextf(ctx, "[GenerateFromTemplate] request invalid, req: %s, error message: %v", logs.JSON(req), resp.Message)
		return
	}

	res, err := svcdoc.GenerateFromTemplate(ctx, req)
	if err != nil {
		logs.ErrorContextf(ctx, "[GenerateFromTemplate] svcdoc.GenerateFromTemplate failed, err: %v", err)
		resp.Code = errcode.ErrCode_InternalError
		resp.Message = "生成文档失败"
		return
	}
	resp.Code = res.Code
	resp.Message = res.Message
	resp.Response = res.Response
}
