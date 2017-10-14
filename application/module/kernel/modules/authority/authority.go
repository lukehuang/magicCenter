package authority

import (
	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/module/kernel/modules/authority/def"
	"muidea.com/magicCenter/application/module/kernel/modules/authority/handler"
)

type authority struct {
	moduleHub        common.ModuleHub
	sessionRegistry  common.SessionRegistry
	routes           []common.Route
	authorityHandler common.AuthorityHandler
}

// LoadModule 加载模块
func LoadModule(configuration common.Configuration, sessionRegistry common.SessionRegistry, moduleHub common.ModuleHub) {
	instance := &authority{
		moduleHub:        moduleHub,
		sessionRegistry:  sessionRegistry,
		authorityHandler: handler.CreateAuthorityHandler(moduleHub, sessionRegistry)}

	moduleHub.RegisterModule(instance)
}

func (instance *authority) ID() string {
	return def.ID
}

func (instance *authority) Name() string {
	return def.Name
}

func (instance *authority) Description() string {
	return def.Description
}

func (instance *authority) Group() string {
	return "kernel"
}

func (instance *authority) Type() int {
	return common.KERNEL
}

func (instance *authority) Status() int {
	return 0
}

func (instance *authority) EntryPoint() interface{} {
	return instance.authorityHandler
}

// Route 路由信息
func (instance *authority) Routes() []common.Route {
	return instance.routes
}

// Startup 启动模块
func (instance *authority) Startup() bool {
	return true
}

// Cleanup 清除模块
func (instance *authority) Cleanup() {

}
