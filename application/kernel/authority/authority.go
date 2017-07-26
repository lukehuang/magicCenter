package authority

/*
鉴权

实现service.Authority接口，提供管理员鉴权处理器和登陆鉴权处理器，返回martini的鉴权handler

应用端通过System获取接口对象
*/
import (
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"muidea.com/magicCenter/application/common"
)

// Authority 权限控制服务
// 控制系统各个模块的访问权限
type Authority interface {
	Verify(authorityHandler common.CASHandler, res http.ResponseWriter, req *http.Request) bool
}

// VerifyHandler 权限校验处理器
// 用于在路由过程中进行权限校验
func VerifyHandler(modHub common.ModuleHub, authority Authority) martini.Handler {
	var authorityHandler common.CASHandler
	mod, ok := modHub.FindModule(common.CASModuleID)
	if ok {
		authorityHandler = nil
		entryPoint := mod.EntryPoint()
		switch entryPoint.(type) {
		case common.CASHandler:
			authorityHandler = entryPoint.(common.CASHandler)
		}
	}

	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
		if authority.Verify(authorityHandler, res, req) {
			// 拥有权限，继续执行
			c.Next()
		} else {
			// 没有权限
		}
	}
}

// CreateAuthority 创建Authority
func CreateAuthority() Authority {
	impl := &impl{}

	return impl
}

type impl struct {
}

func (i *impl) Verify(authorityHandler common.CASHandler, res http.ResponseWriter, req *http.Request) bool {
	if authorityHandler == nil {
		return false
	}

	return authorityHandler.VerifyAccount(res, req)
}
