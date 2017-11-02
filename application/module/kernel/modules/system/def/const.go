package def

import "muidea.com/magicCenter/application/common"

// ID 模块ID
const ID = common.SystemModuleID

// Name 模块名称
const Name = "Magic SystemConfig"

// Description 模块描述信息
const Description = "Magic 系统配置管理"

// URL 模块Url
const URL = "/system"

// GetSystemConfig 获取系统配置
const GetSystemConfig = "/config/"

// PutSystemConfig 更新系统配置
const PutSystemConfig = "/config/"

// GetSystemModule 获取系统模块列表
const GetSystemModule = "/module/"