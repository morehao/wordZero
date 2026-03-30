# WordZero HTTP API & SDK 使用文档

WordZero 提供了两种集成方式：**HTTP API 服务** 和 **Go SDK**。两者均支持以内容或模板方式生成 Word 文档，并将结果自动上传至 S3 兼容存储（如 MinIO、阿里云 OSS、腾讯云 COS 等）。

---

## 目录

- [HTTP API 服务](#http-api-服务)
  - [启动服务](#启动服务)
  - [配置说明](#配置说明)
  - [接口列表](#接口列表)
    - [健康检查](#健康检查)
    - [通过内容生成文档](#通过内容生成文档)
    - [通过模板生成文档](#通过模板生成文档)
- [Go SDK](#go-sdk)
  - [安装](#安装)
  - [初始化客户端](#初始化客户端)
  - [通过内容生成文档](#通过内容生成文档-sdk)
  - [通过模板生成文档](#通过模板生成文档-sdk)

---

## HTTP API 服务

HTTP 服务基于 [gin](https://github.com/gin-gonic/gin) 框架实现，S3 上传使用 [ygpkg/yg-go/storage](https://github.com/ygpkg/yg-go) 包。

### 启动服务

**方式一：使用配置文件**

```bash
./wordzero-server -config config.json
```

配置文件示例（`config.example.json`）：

```json
{
  "host": "0.0.0.0",
  "port": 8080,
  "read_timeout_seconds": 30,
  "write_timeout_seconds": 60,
  "idle_timeout_seconds": 120,
  "template_download_timeout_seconds": 30,
  "s3": {
    "endpoint": "http://localhost:9000",
    "region": "us-east-1",
    "bucket": "wordzero",
    "access_key_id": "minioadmin",
    "secret_access_key": "minioadmin",
    "key_prefix": "documents",
    "use_path_style": true,
    "public_base_url": "http://localhost:9000/wordzero"
  }
}
```

**方式二：使用环境变量**

```bash
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export S3_ENDPOINT=http://localhost:9000
export S3_REGION=us-east-1
export S3_BUCKET=wordzero
export S3_ACCESS_KEY_ID=minioadmin
export S3_SECRET_ACCESS_KEY=minioadmin
export S3_KEY_PREFIX=documents
export S3_USE_PATH_STYLE=true
export S3_PUBLIC_BASE_URL=http://localhost:9000/wordzero

./wordzero-server
```

**方式三：Docker Compose**

```bash
docker-compose up -d
```

### 配置说明

| 字段 | 环境变量 | 说明 | 默认值 |
|------|----------|------|--------|
| `host` | `SERVER_HOST` | 监听地址 | `0.0.0.0` |
| `port` | `SERVER_PORT` | 监听端口 | `8080` |
| `read_timeout_seconds` | — | 读取超时（秒） | `30` |
| `write_timeout_seconds` | — | 写入超时（秒） | `60` |
| `idle_timeout_seconds` | — | 空闲超时（秒） | `120` |
| `template_download_timeout_seconds` | — | 模板下载超时（秒） | `30` |
| `s3.endpoint` | `S3_ENDPOINT` | S3 端点地址（留空则使用 AWS S3） | — |
| `s3.region` | `S3_REGION` | 存储区域 | `us-east-1` |
| `s3.bucket` | `S3_BUCKET` | 存储桶名称 | — |
| `s3.access_key_id` | `S3_ACCESS_KEY_ID` | 访问密钥 ID | — |
| `s3.secret_access_key` | `S3_SECRET_ACCESS_KEY` | 访问密钥 | — |
| `s3.key_prefix` | `S3_KEY_PREFIX` | 对象键前缀（可选） | — |
| `s3.use_path_style` | `S3_USE_PATH_STYLE` | 路径样式 URL（MinIO 等需设为 true） | `false` |
| `s3.public_base_url` | `S3_PUBLIC_BASE_URL` | 公开访问基础 URL（可选） | — |

---

### 接口列表

#### 健康检查

```
GET /health
```

**响应示例：**

```json
{"status": "ok"}
```

---

#### 通过内容生成文档

```
POST /api/v1/documents/content
Content-Type: application/json
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `filename` | string | 否 | 输出文件名，默认 `document.docx` |
| `content` | array | 是 | 文档内容列表 |

**`content` 数组元素字段：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `type` | string | 内容类型：`paragraph`（段落）、`heading`（标题）、`page_break`（分页符） |
| `text` | string | 文本内容 |
| `heading_level` | int | 标题级别 1-6（仅 `heading` 类型有效） |
| `bold` | bool | 是否加粗（仅 `paragraph` 类型有效） |
| `italic` | bool | 是否斜体（仅 `paragraph` 类型有效） |
| `alignment` | string | 对齐方式：`left`、`center`、`right`、`both`（仅 `paragraph` 类型有效） |

**请求示例：**

```json
{
  "filename": "report.docx",
  "content": [
    {
      "type": "heading",
      "text": "年度报告",
      "heading_level": 1
    },
    {
      "type": "paragraph",
      "text": "本报告总结了本年度的业务情况。",
      "bold": false,
      "alignment": "left"
    },
    {
      "type": "heading",
      "text": "第一章：业务概述",
      "heading_level": 2
    },
    {
      "type": "paragraph",
      "text": "本章节介绍业务概况。",
      "italic": true
    },
    {
      "type": "page_break"
    },
    {
      "type": "paragraph",
      "text": "（续页内容）",
      "alignment": "center"
    }
  ]
}
```

**成功响应（200）：**

```json
{
  "url": "http://localhost:9000/wordzero/documents/2024-01-01/1704067200000_report.docx",
  "filename": "report.docx",
  "size": 5432
}
```

**错误响应（4xx/5xx）：**

```json
{
  "error": "错误信息描述"
}
```

**curl 示例：**

```bash
curl -X POST http://localhost:8080/api/v1/documents/content \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "test.docx",
    "content": [
      {"type": "heading", "text": "测试文档", "heading_level": 1},
      {"type": "paragraph", "text": "这是一段测试内容。"}
    ]
  }'
```

---

#### 通过模板生成文档

```
POST /api/v1/documents/template
Content-Type: application/json
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `filename` | string | 否 | 输出文件名，默认 `document.docx` |
| `template_url` | string | 是 | 模板文件的 HTTP/HTTPS URL |
| `variables` | object | 否 | 模板变量（见下方说明） |

**`variables` 变量类型规则：**

| 变量值类型 | 在模板中的用途 |
|-----------|---------------|
| `string` / `number` / 其他基本类型 | 普通变量，用 `{{variable_name}}` 引用 |
| `bool` | 条件变量，用 `{{#if variable_name}}...{{/if}}` 引用 |
| `array`（`[]interface{}`） | 列表变量，用 `{{#each variable_name}}...{{/each}}` 引用 |

**请求示例：**

```json
{
  "filename": "contract.docx",
  "template_url": "https://example.com/templates/contract.docx",
  "variables": {
    "company_name": "示例科技有限公司",
    "date": "2024-01-01",
    "amount": "100,000",
    "is_vip": true,
    "items": [
      {"name": "产品A", "price": "1000"},
      {"name": "产品B", "price": "2000"}
    ]
  }
}
```

**成功响应（200）：**

```json
{
  "url": "http://localhost:9000/wordzero/documents/2024-01-01/1704067200000_contract.docx",
  "filename": "contract.docx",
  "size": 8192
}
```

**curl 示例：**

```bash
curl -X POST http://localhost:8080/api/v1/documents/template \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "output.docx",
    "template_url": "https://example.com/templates/my_template.docx",
    "variables": {
      "title": "我的文档",
      "author": "张三"
    }
  }'
```

---

## Go SDK

SDK 提供了与 HTTP API 功能相同的 Go 语言编程接口，适合在 Go 项目中直接使用。

### 安装

```bash
go get github.com/zerx-lab/wordZero
```

### 初始化客户端

```go
import (
    "github.com/zerx-lab/wordZero/sdk"
    "github.com/zerx-lab/wordZero/internal/s3"
    "time"
)

client, err := sdk.NewClient(&sdk.Config{
    S3Config: s3.Config{
        Endpoint:        "http://localhost:9000",
        Region:          "us-east-1",
        Bucket:          "wordzero",
        AccessKeyID:     "minioadmin",
        SecretAccessKey: "minioadmin",
        KeyPrefix:       "documents",
        UsePathStyle:    true,
        PublicBaseURL:   "http://localhost:9000/wordzero",
    },
    HTTPTimeout: 30 * time.Second,
})
if err != nil {
    log.Fatalf("创建 SDK 客户端失败: %v", err)
}
```

### 通过内容生成文档 (SDK)

```go
import (
    "context"
    "fmt"
    "log"
    "github.com/zerx-lab/wordZero/sdk"
)

resp, err := client.GenerateFromContent(context.Background(), &sdk.ContentRequest{
    Filename: "report.docx",
    Content: []sdk.ContentItem{
        {
            Type:         sdk.ContentTypeHeading,
            Text:         "年度报告",
            HeadingLevel: 1,
        },
        {
            Type:      sdk.ContentTypeParagraph,
            Text:      "本报告总结了本年度的业务情况。",
            Alignment: "left",
        },
        {
            Type:      sdk.ContentTypeParagraph,
            Text:      "重要提示：请仔细阅读。",
            Bold:      true,
            Italic:    true,
            Alignment: "center",
        },
        {
            Type: sdk.ContentTypePageBreak,
        },
        {
            Type: sdk.ContentTypeHeading,
            Text: "第二章：详细内容",
            HeadingLevel: 2,
        },
    },
})
if err != nil {
    log.Fatalf("生成文档失败: %v", err)
}

fmt.Printf("文档生成成功: URL=%s, 大小=%d 字节\n", resp.URL, resp.Size)
```

**`ContentItem` 字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `Type` | `ContentType` | 内容类型：`ContentTypeParagraph`、`ContentTypeHeading`、`ContentTypePageBreak` |
| `Text` | string | 文本内容 |
| `HeadingLevel` | int | 标题级别 1-6（仅 `ContentTypeHeading` 有效） |
| `Bold` | bool | 是否加粗（仅 `ContentTypeParagraph` 有效） |
| `Italic` | bool | 是否斜体（仅 `ContentTypeParagraph` 有效） |
| `Alignment` | string | 对齐方式：`left`、`center`、`right`、`both`（仅 `ContentTypeParagraph` 有效） |

### 通过模板生成文档 (SDK)

```go
import (
    "context"
    "fmt"
    "log"
    "github.com/zerx-lab/wordZero/sdk"
)

resp, err := client.GenerateFromTemplate(context.Background(), &sdk.TemplateRequest{
    Filename:    "contract.docx",
    TemplateURL: "https://example.com/templates/contract.docx",
    Variables: map[string]interface{}{
        // 普通文本变量（在模板中使用 {{company_name}}）
        "company_name": "示例科技有限公司",
        "date":         "2024-01-01",
        // 条件变量（在模板中使用 {{#if is_vip}}...{{/if}}）
        "is_vip": true,
        // 列表变量（在模板中使用 {{#each items}}...{{/each}}）
        "items": []interface{}{
            map[string]interface{}{"name": "产品A", "price": "1000"},
            map[string]interface{}{"name": "产品B", "price": "2000"},
        },
    },
})
if err != nil {
    log.Fatalf("生成文档失败: %v", err)
}

fmt.Printf("文档生成成功: URL=%s, 大小=%d 字节\n", resp.URL, resp.Size)
```

**`TemplateRequest` 字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `Filename` | string | 否 | 输出文件名，默认 `document.docx` |
| `TemplateURL` | string | 是 | 模板文件的 HTTP/HTTPS URL |
| `Variables` | `map[string]interface{}` | 否 | 模板变量（bool 为条件变量，`[]interface{}` 为列表变量，其他为普通变量） |

**`GenerateResponse` 响应字段说明：**

| 字段 | 类型 | 说明 |
|------|------|------|
| `URL` | string | 上传到 S3 后的访问 URL |
| `Filename` | string | 生成的文件名 |
| `Size` | int64 | 文件大小（字节） |

---

## 模板语法参考

WordZero 模板支持以下语法：

| 语法 | 说明 | 示例 |
|------|------|------|
| `{{variable}}` | 普通变量替换 | `{{company_name}}` → `示例科技有限公司` |
| `{{#if condition}}...{{/if}}` | 条件块 | `{{#if is_vip}}VIP客户{{/if}}` |
| `{{#each list}}...{{/each}}` | 循环块 | `{{#each items}}{{name}}{{/each}}` |

---

## S3 兼容存储配置示例

### MinIO（本地开发）

```json
{
  "endpoint": "http://localhost:9000",
  "region": "us-east-1",
  "bucket": "wordzero",
  "access_key_id": "minioadmin",
  "secret_access_key": "minioadmin",
  "use_path_style": true,
  "public_base_url": "http://localhost:9000/wordzero"
}
```

### 阿里云 OSS

```json
{
  "endpoint": "https://oss-cn-hangzhou.aliyuncs.com",
  "region": "cn-hangzhou",
  "bucket": "your-bucket",
  "access_key_id": "your-access-key-id",
  "secret_access_key": "your-access-key-secret",
  "use_path_style": false,
  "public_base_url": "https://your-bucket.oss-cn-hangzhou.aliyuncs.com"
}
```

### AWS S3

```json
{
  "region": "us-east-1",
  "bucket": "your-bucket",
  "access_key_id": "your-access-key-id",
  "secret_access_key": "your-secret-access-key"
}
```
