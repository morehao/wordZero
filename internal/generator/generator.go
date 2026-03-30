// Package generator 提供从内存中的模板文档生成Word文档的共享功能
package generator

import (
	"bytes"
	"io"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// RenderTemplateFromBytes 从模板字节数据渲染Word文档
// templateBytes: 模板文档的原始字节
// data: 模板变量数据
// 返回渲染后的文档对象
func RenderTemplateFromBytes(templateBytes []byte, data *document.TemplateData) (*document.Document, error) {
	// 从内存加载模板文档
	templateDoc, err := document.OpenFromMemory(io.NopCloser(bytes.NewReader(templateBytes)))
	if err != nil {
		return nil, err
	}

	// 使用TemplateEngine加载并渲染模板
	engine := document.NewTemplateEngine()
	_, err = engine.LoadTemplateFromDocument("tpl", templateDoc)
	if err != nil {
		return nil, err
	}

	// 渲染模板到文档
	return engine.RenderTemplateToDocument("tpl", data)
}

// BuildTemplateData 从通用map构建TemplateData
// - bool类型变量将作为条件变量
// - []interface{}类型变量将作为列表变量
// - 其他类型将作为普通变量
func BuildTemplateData(variables map[string]interface{}) *document.TemplateData {
	data := &document.TemplateData{
		Variables:  make(map[string]interface{}),
		Lists:      make(map[string][]interface{}),
		Conditions: make(map[string]bool),
	}

	for k, v := range variables {
		switch val := v.(type) {
		case bool:
			data.Conditions[k] = val
		case []interface{}:
			data.Lists[k] = val
		default:
			data.Variables[k] = v
		}
	}

	return data
}
