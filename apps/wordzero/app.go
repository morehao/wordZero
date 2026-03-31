package wordzero

import (
	"sync"

	"github.com/ygpkg/yg-go/apis/runtime/server"
	"github.com/zerx-lab/wordZero/apps/wordzero/internal/apis"
	"gorm.io/gorm"
)

// @title wordZero API
// @description Word 文档生成服务 API
// @host localhost:8080
// @BasePath /apis/p
// @schemes http
// @accept json
// @produce json

// @param Env header string true "dev"

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

// Routers 路由注册
func Routers(eng *server.Router) error {
	apis.RegistryRouter(eng)
	return nil
}

// Migrates 补全数据表及数据库索引
func Migrates(db *gorm.DB) error {
	return nil
}

var onceStart sync.Once

// RunJob 启动定时任务
func RunJob() error {
	return nil
}
