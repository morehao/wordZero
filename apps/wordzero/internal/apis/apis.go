package apis

import (
	"github.com/ygpkg/yg-go/apis/runtime/server"
)

func RegistryRouter(eng *server.Router) {
	eng.HandleDoc("wordzero")

	eng.P("wordzero.GenerateFromContent", GenerateFromContent)
	eng.P("wordzero.GenerateFromTemplate", GenerateFromTemplate)
}
