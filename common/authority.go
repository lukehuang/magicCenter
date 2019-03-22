package common

import (
	"net/http"

	"github.com/muidea/magicCommon/model"
)

// AuthorityHandler 鉴权处理器
type AuthorityHandler interface {
	// 校验权限
	VerifyAuthority(res http.ResponseWriter, req *http.Request) bool

	// 查询所有ACL
	QueryAllACL() []model.ACL
	// 查询指定Module的ACL
	QueryACLByModule(module string) []model.ACL

	// 查询指定ACL
	QueryACLByID(id int) (model.ACL, bool)
	// 新增ACL
	InsertACL(url, method, module string, status int, authGroup int) (model.ACL, bool)

	// 删除ACL
	DeleteACL(id int) bool

	// 更新指定的ACL的授权组和状态
	UpdateACL(acl model.ACL) bool

	// 批量更新ACL状态
	UpdateACLStatus(enableList []int, disableList []int) bool

	// 查询指定ACL的授权组信息
	QueryACLAuthGroup(id int) (model.AuthGroup, bool)
	// 更新指定ACL的授权组信息
	UpdateACLAuthGroup(id, authGroup int) bool

	// 查询所有Module的用户信息
	QueryAllModuleUser() []model.ModuleUserInfo
	// 查询指定模块的拥有者
	QueryModuleUserAuthGroup(module string) []model.UserAuthGroup
	// 更新指定模块的拥有者
	UpdateModuleUserAuthGroup(module string, userAuthGroups []model.UserAuthGroup) bool

	// 查询所有用户的Module信息
	QueryAllUserModule() []model.UserModuleInfo
	// 查询指定用户使用的模块信息
	QueryUserModuleAuthGroup(user int) []model.ModuleAuthGroup
	// 更新指定用户使用的模块信息
	UpdateUserModuleAuthGroup(user int, moduleAuthGroups []model.ModuleAuthGroup) bool

	// 查询指定用户的ACL列表
	QueryUserACL(user int) []model.ACL

	// 查询所有授权组定义
	QueryAllAuthGroupDef() []model.AuthGroup
}
