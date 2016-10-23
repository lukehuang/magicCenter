package loader

import (
	"magiccenter/kernel/modules/account"
	"magiccenter/kernel/modules/cache"
	"magiccenter/kernel/modules/content"
	"magiccenter/kernel/modules/dashboard"
	"magiccenter/kernel/modules/dashboard/contentmanage"
	"magiccenter/kernel/modules/dashboard/modulemanage"
	"magiccenter/kernel/modules/dashboard/systemsetting"
	"magiccenter/kernel/modules/mail"
)

// LoadAllModules 加载所有模块
func LoadAllModules() {

	mail.LoadModule()

	cache.LoadModule()

	dashboard.LoadModule()

	systemsetting.LoadModule()

	account.LoadModule()

	content.LoadModule()

	contentmanage.LoadModule()

	modulemanage.LoadModule()
}