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
| `type` | string | 内容类型：`paragraph`（段落）、`heading`（标题）、`page_break`（分页符）、`table`（表格） |
| `text` | string | 文本内容 |
| `heading_level` | int | 标题级别 1-6（仅 `heading` 类型有效） |
| `bold` | bool | 是否加粗（仅 `paragraph` 类型有效） |
| `italic` | bool | 是否斜体（仅 `paragraph` 类型有效） |
| `alignment` | string | 对齐方式：`left`、`center`、`right`、`both`（仅 `paragraph` 类型有效） |
| `table_data` | array | 表格数据，二维数组（仅 `table` 类型有效），如 `[["列1", "列2"], ["行1", "行2"]]` |
| `table_width` | int | 表格宽度，单位：磅（1磅≈3.53毫米），默认 9000（仅 `table` 类型有效） |

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
    },
    {
      "type": "heading",
      "text": "销售数据统计",
      "heading_level": 2
    },
    {
      "type": "table",
      "table_width": 9000,
      "table_data": [
        ["产品名称", "销量（件）", "销售额（元）"],
        ["iPhone 15", "1250", "8750000"],
        ["MacBook Pro", "380", "5320000"],
        ["iPad Air", "620", "2480000"],
        ["AirPods Pro", "2100", "4200000"]
      ]
    },
    {
      "type": "paragraph",
      "text": "结论：iPhone 15 销量最高，AirPods Pro 销量增长最快。",
      "bold": true
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
    "github.com/zerx-lab/wordZero/wordzero"
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
    "github.com/zerx-lab/wordZero/wordzero"
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
        {
            Type:         sdk.ContentTypeTable,
            TableData:    [][]string{
                {"产品名称", "销量", "销售额"},
                {"iPhone 15", "1250", "8750000"},
                {"MacBook Pro", "380", "5320000"},
                {"iPad Air", "620", "2480000"},
            },
            TableWidth: 9000,
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
| `Type` | `ContentType` | 内容类型：`ContentTypeParagraph`、`ContentTypeHeading`、`ContentTypePageBreak`、`ContentTypeTable` |
| `Text` | string | 文本内容 |
| `HeadingLevel` | int | 标题级别 1-6（仅 `ContentTypeHeading` 有效） |
| `Bold` | bool | 是否加粗（仅 `ContentTypeParagraph` 有效） |
| `Italic` | bool | 是否斜体（仅 `ContentTypeParagraph` 有效） |
| `Alignment` | string | 对齐方式：`left`、`center`、`right`、`both`（仅 `ContentTypeParagraph` 有效） |
| `TableData` | `[][]string` | 表格数据，二维数组（仅 `ContentTypeTable` 有效） |
| `TableWidth` | int | 表格宽度，单位：磅，默认 9000（仅 `ContentTypeTable` 有效） |

### 通过模板生成文档 (SDK)

```go
import (
    "context"
    "fmt"
    "log"
    "github.com/zerx-lab/wordZero/wordzero"
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

---

## 错误码说明

### HTTP 状态码

| 状态码 | 说明 | 常见原因 |
|--------|------|----------|
| `200` | 请求成功 | 文档生成并上传成功 |
| `400` | 请求参数错误 | 请求体格式错误、必填字段缺失 |
| `401` | 认证失败 | 认证信息无效或缺失 |
| `413` | 请求体过大 | 模板文件或请求内容超出限制 |
| `422` | 请求格式正确但无法处理 | 请求参数值无效 |
| `500` | 服务器内部错误 | 文档生成失败、S3上传失败等 |
| `502` | 上游服务错误 | S3服务不可用 |
| `504` | 网关超时 | 模板下载超时 |

### 业务错误码

响应体中的 `code` 和 `message` 字段用于描述具体错误原因：

| code | message | 说明 | 解决方案 |
|------|---------|------|----------|
| `bad_request` | 内容列表不能为空 | `content` 数组为空 | 确保 `content` 数组至少包含一个元素 |
| `bad_request` | 模板URL不能为空 | `template_url` 为空 | 提供有效的模板文件 URL |
| `bad_request` | 未知内容类型: xxx | `type` 字段值无效 | 使用有效的类型：`paragraph`、`heading`、`page_break`、`table` |
| `bad_request` | 表格数据不能为空 | `table` 类型的 `table_data` 为空 | 确保 `table_data` 至少包含一行数据 |
| `bad_request` | 表格列数不能为0 | `table` 类型的 `table_data` 每行都为空 | 确保 `table_data` 至少包含一列数据 |
| `bad_request` | 标题级别必须在 1-6 之间 | `heading_level` 超出范围 | 确保 `heading_level` 在 1-6 范围内 |
| `internal_error` | 下载模板失败 | 模板 URL 无法访问 | 检查模板 URL 是否可访问 |
| `internal_error` | 渲染模板失败 | 模板语法错误或变量不匹配 | 检查模板语法和变量名称 |
| `internal_error` | 生成文档失败 | 文档生成过程中出错 | 检查内容格式是否正确 |
| `internal_error` | 上传文档失败 | S3 上传失败 | 检查 S3 配置是否正确 |

### 错误响应示例

**400 错误 - 请求参数错误：**

```json
{
  "code": "bad_request",
  "message": "内容列表不能为空",
  "request_id": "req_abc123"
}
```

**500 错误 - 服务器内部错误：**

```json
{
  "code": "internal_error",
  "message": "上传文档失败: connection refused",
  "request_id": "req_def456"
}
```

### 错误处理建议

1. **4xx 错误**：通常是请求问题，检查请求参数是否正确
2. **5xx 错误**：可能是服务端临时问题，建议实现重试机制
3. **S3 错误**：检查 S3 配置、网络连接、存储桶权限等

---

## 更多使用示例

### 高级内容生成示例

#### 多级标题文档

```json
{
  "filename": "user_manual.docx",
  "content": [
    {"type": "heading", "text": "用户手册", "heading_level": 1},
    {"type": "paragraph", "text": "本文档介绍产品使用方法"},
    {"type": "heading", "text": "第一章：快速开始", "heading_level": 2},
    {"type": "heading", "text": "1.1 安装", "heading_level": 3},
    {"type": "paragraph", "text": "请按照以下步骤安装..."},
    {"type": "heading", "text": "1.2 配置", "heading_level": 3},
    {"type": "paragraph", "text": "首次启动后需要配置..."},
    {"type": "heading", "text": "第二章：高级功能", "heading_level": 2},
    {"type": "heading", "text": "2.1 自定义设置", "heading_level": 3},
    {"type": "paragraph", "text": "高级用户可以自定义..."},
    {"type": "page_break"},
    {"type": "heading", "text": "附录", "heading_level": 1},
    {"type": "paragraph", "text": "常见问题解答..."}
  ]
}
```

#### 包含表格的文档

```json
{
  "filename": "sales_report.docx",
  "content": [
    {"type": "heading", "text": "2024年销售报告", "heading_level": 1},
    {"type": "paragraph", "text": "本报告展示本年度销售数据汇总。"},
    {"type": "heading", "text": "一季度销售数据", "heading_level": 2},
    {
      "type": "table",
      "table_width": 9000,
      "table_data": [
        ["产品类别", "销售额（元）", "占比"],
        ["电子产品", "5,250,000", "35%"],
        ["办公设备", "3,150,000", "21%"],
        ["家居用品", "2,100,000", "14%"],
        ["其他", "4,500,000", "30%"]
      ]
    },
    {"type": "heading", "text": "结论", "heading_level": 2},
    {"type": "paragraph", "text": "电子产品销售额占比最高，是公司主要收入来源。", "bold": true}
  ]
}
```

#### 格式化文本文档

```json
{
  "filename": "notice.docx",
  "content": [
    {"type": "heading", "text": "重要通知", "heading_level": 1},
    {"type": "paragraph", "text": "尊敬的用户：", "alignment": "left"},
    {"type": "paragraph", "text": "    感谢您使用我们的服务。本通知旨在告知您关于系统维护的最新安排。", "alignment": "both"},
    {"type": "paragraph", "text": "维护时间：2024年1月15日 00:00 - 06:00", "bold": true},
    {"type": "paragraph", "text": "如有疑问，请联系客服。", "italic": true, "alignment": "right"},
    {"type": "paragraph", "text": "--------------------------", "alignment": "center"},
    {"type": "paragraph", "text": "此致", "alignment": "right"},
    {"type": "paragraph", "text": "敬礼", "alignment": "right"}
  ]
}
```

### 模板高级用法示例

#### 条件渲染 - VIP 客户合同

```json
{
  "filename": "contract.docx",
  "template_url": "https://example.com/templates/vip_contract.docx",
  "variables": {
    "customer_name": "张三",
    "contract_type": "VIP服务协议",
    "is_vip": true,
    "discount": "8折",
    "service_level": "专属客服",
    "validity_period": "24个月"
  }
}
```

模板内容示例（`vip_contract.docx`）：

```
服务合同

客户名称：{{customer_name}}
合同类型：{{contract_type}}

{{#if is_vip}}
VIP客户专享：
- 折扣：{{discount}}
- 服务级别：{{service_level}}
{{else}}
标准客户：
- 折扣：9折
- 服务级别：标准客服
{{/if}}

合同期限：{{validity_period}}
```

#### 循环遍历 - 产品列表

```json
{
  "filename": "product_list.docx",
  "template_url": "https://example.com/templates/product_list.docx",
  "variables": {
    "title": "产品目录",
    "category": "电子产品",
    "products": [
      {"name": "手机X1", "price": "3999", "stock": 100},
      {"name": "平板Pro", "price": "4999", "stock": 50},
      {"name": "无线耳机", "price": "899", "stock": 200}
    ]
  }
}
```

模板内容示例：

```
{{title}}

类别：{{category}}

产品列表：
{{#each products}}
- 产品名称：{{name}}
  价格：￥{{price}}
  库存：{{stock}}件
{{/each}}
```

#### 嵌套循环 - 部门员工列表

```json
{
  "filename": "org_chart.docx",
  "template_url": "https://example.com/templates/org_chart.docx",
  "variables": {
    "company_name": "示例科技",
    "departments": [
      {
        "name": "技术部",
        "manager": "李经理",
        "employees": [
          {"name": "张三", "role": "高级工程师"},
          {"name": "李四", "role": "工程师"}
        ]
      },
      {
        "name": "市场部",
        "manager": "王经理",
        "employees": [
          {"name": "赵六", "role": "市场专员"},
          {"name": "钱七", "role": "策划专员"}
        ]
      }
    ]
  }
}
```

模板内容示例：

```
{{company_name}} 组织架构

{{#each departments}}
部门：{{name}}
负责人：{{manager}}
成员：
{{#each employees}}
  - {{name}} ({{role}})
{{/each}}
---
{{/each}}
```

### 实际业务场景示例

#### 合同生成

```json
{
  "filename": "contract_20240101.docx",
  "template_url": "https://example.com/templates/contract.docx",
  "variables": {
    "contract_no": "CT-2024-001",
    "party_a": "甲方科技有限公司",
    "party_b": "乙方企业有限公司",
    "sign_date": "2024年1月1日",
    "amount": "人民币壹佰万元整",
    "payment_method": "按季度支付",
    "is_long_term": true,
    "items": [
      {"name": "软件开发服务", "price": "500000"},
      {"name": "技术支持服务", "price": "300000"},
      {"name": "培训服务", "price": "200000"}
    ]
  }
}
```

#### 报告生成

```json
{
  "filename": "monthly_report.docx",
  "template_url": "https://example.com/templates/monthly_report.docx",
  "variables": {
    "report_title": "2024年1月月度报告",
    "department": "销售部",
    "author": "张三",
    "total_sales": "1,850,000",
    "growth_rate": "15%",
    "new_customers": 45,
    "is_completed": true,
    "highlights": [
      {"title": "华东区业绩突破", "content": "销售额同比增长30%"},
      {"title": "新客户开拓", "content": "本月新增客户45家"},
      {"title": "产品创新", "content": "推出2款新产品"}
    ]
  }
}
```

### HTTP API 调用示例

#### curl - 通过内容生成文档

```bash
curl -X POST http://localhost:8080/api/v1/documents/content \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "test.docx",
    "content": [
      {"type": "heading", "text": "测试文档", "heading_level": 1},
      {"type": "paragraph", "text": "这是使用 curl 生成的 Word 文档。"},
      {"type": "paragraph", "text": "支持多种格式：", "bold": true},
      {"type": "paragraph", "text": "- 加粗文本", "bold": true},
      {"type": "paragraph", "text": "- 斜体文本", "italic": true},
      {"type": "paragraph", "text": "- 居中对齐", "alignment": "center"}
    ]
  }'
```

#### curl - 通过模板生成文档

```bash
curl -X POST http://localhost:8080/api/v1/documents/template \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "template_test.docx",
    "template_url": "https://example.com/templates/sample.docx",
    "variables": {
      "title": "测试文档",
      "author": "张三",
      "date": "2024-01-01"
    }
  }'
```

#### curl - 批量生成

```bash
# 循环生成多个文档
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/documents/content \
    -H "Content-Type: application/json" \
    -d "{
      \"filename\": \"batch_$i.docx\",
      \"content\": [
        {\"type\": \"heading\", \"text\": \"文档 $i\", \"heading_level\": 1},
        {\"type\": \"paragraph\", \"text\": \"这是第 $i 个批量生成的文档。\"}
      ]
    }"
done
```

### SDK 使用示例

#### Go SDK - 合同生成

```go
import (
    "context"
    "fmt"
    "log"
    "github.com/zerx-lab/wordZero/wordzero"
)

func main() {
    client, _ := sdk.NewClient(&sdk.Config{
        S3Config: s3.Config{ /* 配置 */ },
    })

    resp, err := client.GenerateFromTemplate(context.Background(), &sdk.TemplateRequest{
        Filename:    "contract.docx",
        TemplateURL: "https://example.com/templates/contract.docx",
        Variables: map[string]interface{}{
            "contract_no":  "CT-2024-001",
            "party_a":      "甲方科技有限公司",
            "party_b":      "乙方企业有限公司",
            "sign_date":    "2024年1月1日",
            "amount":       "人民币壹佰万元整",
            "payment_method": "按季度支付",
            "is_long_term": true,
            "items": []interface{}{
                map[string]interface{}{"name": "软件开发服务", "price": "500000"},
                map[string]interface{}{"name": "技术支持服务", "price": "300000"},
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("合同生成成功: %s\n", resp.URL)
}
```

---

## SDK 高级功能

### 自定义 HTTP 客户端

SDK 允许自定义 HTTP 客户端，用于下载模板文件。你可以根据需要配置超时、代理、TLS 等选项。

```go
import (
    "crypto/tls"
    "net/http"
    "net/url"
    "time"
    "github.com/zerx-lab/wordZero/wordzero"
    "github.com/zerx-lab/wordZero/internal/s3"
)

func main() {
    // 自定义 HTTP 传输配置
    transport := &http.Transport{
        // 配置代理
        Proxy: func(req *http.Request) (*url.URL, error) {
            return url.Parse("http://proxy.example.com:8080")
        },
        // 跳过 TLS 验证（仅用于开发环境）
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        // 配置连接池
        MaxIdleConns:        10,
        MaxIdleConnsPerHost: 5,
        IdleConnTimeout:     30 * time.Second,
    }

    // 自定义 HTTP 客户端
    httpClient := &http.Client{
        Timeout:   60 * time.Second, // 模板下载超时 60 秒
        Transport: transport,
    }

    // 使用自定义 HTTP 客户端创建 SDK 客户端
    // 注意：SDK 内部使用 http.DefaultClient，你需要通过配置传递
    // 这里展示配置方式，实际使用请参考具体 SDK 版本
}
```

### 错误处理与重试机制

在实际生产环境中，建议实现错误处理和重试机制，以提高系统稳定性。

```go
import (
    "context"
    "errors"
    "time"
    "github.com/zerx-lab/wordZero/wordzero"
)

const maxRetries = 3
const retryDelay = time.Second

func generateWithRetry(client *sdk.Client, req *sdk.TemplateRequest) (*sdk.GenerateResponse, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        resp, err := client.GenerateFromTemplate(context.Background(), req)
        if err == nil {
            return resp, nil
        }

        lastErr = err

        // 判断是否需要重试
        if shouldRetry(err) {
            // 指数退避
            delay := retryDelay * time.Duration(i+1)
            time.Sleep(delay)
            continue
        }

        // 非重试错误，直接返回
        return nil, err
    }

    return nil, errors.New("达到最大重试次数: " + lastErr.Error())
}

func shouldRetry(err error) bool {
    // 根据错误类型判断是否需要重试
    errStr := err.Error()
    retryableErrors := []string{
        "connection refused",
        "timeout",
        "temporary failure",
    }

    for _, pattern := range retryableErrors {
        if contains(errStr, pattern) {
            return true
        }
    }
    return false
}

func contains(s, substr string) bool {
    return len(s) >= len(substr) && 
           (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
    for i := start; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
```

### 超时控制

为长时间运行的文档生成任务设置超时。

```go
import (
    "context"
    "time"
    "github.com/zerx-lab/wordZero/wordzero"
)

func main() {
    // 创建带超时控制的上下文
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    client, _ := sdk.NewClient(&sdk.Config{
        S3Config:   s3.Config{ /* 配置 */ },
        HTTPTimeout: 30 * time.Second,
    })

    resp, err := client.GenerateFromTemplate(ctx, &sdk.TemplateRequest{
        Filename:    "large_template.docx",
        TemplateURL: "https://example.com/large_template.docx",
        Variables:   map[string]interface{}{},
    })

    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            // 处理超时
            println("文档生成超时")
        } else {
            // 其他错误
            println("错误:", err.Error())
        }
    }

    println("文档 URL:", resp.URL)
}
```

### 并发处理

在需要批量生成文档时，可以使用并发来提高效率。

```go
import (
    "context"
    "sync"
    "github.com/zerx-lab/wordZero/wordzero"
)

type GenerateTask struct {
    Filename    string
    TemplateURL string
    Variables   map[string]interface{}
}

type GenerateResult struct {
    Filename string
    URL      string
    Size     int64
    Error    error
}

func batchGenerate(client *sdk.Client, tasks []GenerateTask) []GenerateResult {
    results := make([]GenerateResult, len(tasks))

    // 使用 WaitGroup 实现并发控制
    var wg sync.WaitGroup
    // 限制并发数为 5
    semaphore := make(chan struct{}, 5)

    for i, task := range tasks {
        wg.Add(1)
        
        go func(i int, task GenerateTask) {
            // 获取信号量
            semaphore <- struct{}{}
            defer wg.Done()
            defer func() { <-semaphore }()

            resp, err := client.GenerateFromTemplate(context.Background(), &sdk.TemplateRequest{
                Filename:    task.Filename,
                TemplateURL: task.TemplateURL,
                Variables:   task.Variables,
            })

            if err != nil {
                results[i] = GenerateResult{
                    Filename: task.Filename,
                    Error:    err,
                }
            } else {
                results[i] = GenerateResult{
                    Filename: task.Filename,
                    URL:      resp.URL,
                    Size:     resp.Size,
                }
            }
        }(i, task)
    }

    wg.Wait()
    return results
}
```

### S3 高级配置

#### 自定义对象键生成策略

可以通过配置实现自定义的对象键生成策略。

```go
import (
    "fmt"
    "time"
    "github.com/zerx-lab/wordZero/pkg/s3"
)

// 自定义对象键生成器
func customKeyGenerator(filename string) string {
    now := time.Now()
    dateStr := now.Format("2006-01-02")
    timestamp := now.UnixMilli()
    return fmt.Sprintf("custom/%s/%d_%s", dateStr, timestamp, filename)
}

// 在配置中使用
// 注意：需要修改 s3 包以支持自定义键生成函数
```

#### 自定义 Content-Type

文档默认使用 `application/vnd.openxmlformats-officedocument.wordprocessingml.document`，如需自定义可以修改上传逻辑。

```go
// SDK 内部使用固定的 Content-Type
// 如需自定义，可以修改 sdk/client.go 中的 Upload 调用
// 将最后一个参数改为自定义的 Content-Type
```

### 日志集成

在生产环境中，建议集成日志来监控 SDK 运行状态。

```go
import (
    "log"
    "github.com/zerx-lab/wordZero/wordzero"
)

func main() {
    client, err := sdk.NewClient(&sdk.Config{
        S3Config:   s3.Config{ /* 配置 */ },
        HTTPTimeout: 30,
    })
    if err != nil {
        // 使用你项目的日志框架
        log.Printf("创建 SDK 客户端失败: %v", err)
        return
    }

    // 生成文档并记录日志
    resp, err := client.GenerateFromContent(context.Background(), &sdk.ContentRequest{
        Filename: "test.docx",
        Content:  []sdk.ContentItem{ /* ... */ },
    })
    if err != nil {
        log.Printf("生成文档失败: %v", err)
        return
    }

    log.Printf("文档生成成功: filename=%s, url=%s, size=%d", 
        resp.Filename, resp.URL, resp.Size)
}
```

### 性能优化建议

1. **连接池复用**：创建一次 SDK 客户端并复用，避免频繁创建
2. **并发控制**：批量生成时使用信号量限制并发数，避免资源耗尽
3. **模板缓存**：对于常用模板，可以先下载到本地，减少网络请求
4. **超时设置**：根据实际需求合理设置超时时间
5. **错误重试**：对临时性错误实现重试机制，提高可用性
