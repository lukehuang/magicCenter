package fileregistry

import (
	"github.com/muidea/magicCenter/common"
	"github.com/muidea/magicCenter/module/modules/fileregistry/def"
	"github.com/muidea/magicCenter/module/modules/fileregistry/handler"
	"github.com/muidea/magicCenter/module/modules/fileregistry/route"
	common_const "github.com/muidea/magicCommon/common"
)

type fileRegistry struct {
	routes             []common.Route
	fileRegistryHanler common.FileRegistryHandler
}

// LoadModule 加载Static模块
func LoadModule(configuration common.Configuration, sessionRegistry common.SessionRegistry, moduleHub common.ModuleHub) {
	fileRegistryHanler := handler.CreateFileRegistryHandler(configuration, sessionRegistry, moduleHub)

	instance := &fileRegistry{fileRegistryHanler: fileRegistryHanler}

	instance.routes = route.AppendFileRegistryRoute(instance.routes, instance.fileRegistryHanler)

	moduleHub.RegisterModule(instance)
}

// ID ID
func (instance *fileRegistry) ID() string {
	return def.ID
}

// Name 名称
func (instance *fileRegistry) Name() string {
	return def.Name
}

// Description 名称
func (instance *fileRegistry) Description() string {
	return def.Description
}

func (instance *fileRegistry) Group() string {
	return "resource"
}

func (instance *fileRegistry) Type() int {
	return common_const.INTERNAL
}

func (instance *fileRegistry) Status() int {
	return common_const.ACTIVE
}

func (instance *fileRegistry) EntryPoint() interface{} {
	return instance.fileRegistryHanler
}

// Route 路由信息
func (instance *fileRegistry) Routes() []common.Route {
	return instance.routes
}

// Startup 启动模块
func (instance *fileRegistry) Startup() bool {
	return true
}

// Cleanup 清除模块
func (instance *fileRegistry) Cleanup() {

}
