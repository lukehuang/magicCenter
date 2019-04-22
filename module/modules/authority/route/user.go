package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/muidea/magicCenter/common"
	"github.com/muidea/magicCenter/module/modules/authority/def"
	common_const "github.com/muidea/magicCommon/common"
	common_def "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/model"
)

// CreateQueryUserRoute 新建GetUserModuleAuthGroupRoute
func CreateQueryUserRoute(authorityHandler common.AuthorityHandler, accountHandler common.AccountHandler, moduleHub common.ModuleHub) common.Route {
	i := userGetRoute{authorityHandler: authorityHandler, accountHandler: accountHandler, moduleHub: moduleHub}
	return &i
}

// CreateGetUserByIDRoute 新建UserModulePutRoute
func CreateGetUserByIDRoute(authorityHandler common.AuthorityHandler, accountHandler common.AccountHandler, moduleHub common.ModuleHub) common.Route {
	i := userGetByIDRoute{authorityHandler: authorityHandler, accountHandler: accountHandler, moduleHub: moduleHub}
	return &i
}

// CreatePutUserRoute 新建UserACLGetRoute
func CreatePutUserRoute(authorityHandler common.AuthorityHandler, accountHandler common.AccountHandler, moduleHub common.ModuleHub) common.Route {
	i := userPutRoute{authorityHandler: authorityHandler, accountHandler: accountHandler, moduleHub: moduleHub}
	return &i
}

type userGetRoute struct {
	authorityHandler common.AuthorityHandler
	accountHandler   common.AccountHandler
	moduleHub        common.ModuleHub
}

func (i *userGetRoute) Method() string {
	return common.GET
}

func (i *userGetRoute) Pattern() string {
	return net.JoinURL(def.URL, def.QueryUser)
}

func (i *userGetRoute) Handler() interface{} {
	return i.getHandler
}

func (i *userGetRoute) AuthGroup() int {
	return common_const.UserAuthGroup.ID
}

func (i *userGetRoute) getHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("getHandler")

	result := common_def.GetUserModuleInfoListResult{}
	for true {
		userModuleInfo := i.authorityHandler.QueryAllUserModule()
		for _, val := range userModuleInfo {
			user, ok := i.accountHandler.FindUserByID(val.User)
			if ok {
				view := model.UserModuleInfoView{}
				view.User = user.User

				for _, v := range val.Module {
					mod, _ := i.moduleHub.FindModule(v)
					info := model.Module{ID: mod.ID(), Name: mod.Name()}

					view.Module = append(view.Module, info)
				}

				result.User = append(result.User, view)
			}
		}

		result.ErrorCode = common_def.Success
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type userGetByIDRoute struct {
	authorityHandler common.AuthorityHandler
	accountHandler   common.AccountHandler
	moduleHub        common.ModuleHub
}

func (i *userGetByIDRoute) Method() string {
	return common.GET
}

func (i *userGetByIDRoute) Pattern() string {
	return net.JoinURL(def.URL, def.GetUserByID)
}

func (i *userGetByIDRoute) Handler() interface{} {
	return i.getByIDHandler
}

func (i *userGetByIDRoute) AuthGroup() int {
	return common_const.MaintainerAuthGroup.ID
}

func (i *userGetByIDRoute) getByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("getByIDHandler")

	result := common_def.GetUserAuthGroupInfoResult{}
	for true {
		_, strID := net.SplitRESTAPI(r.URL.Path)
		id, err := strconv.Atoi(strID)
		if err != nil {
			result.ErrorCode = common_def.Failed
			result.Reason = "非法参数"
			break
		}

		user, ok := i.accountHandler.FindUserByID(id)
		if ok {
			result.User.UserDetail = user

			result.User.Group = i.accountHandler.GetGroups(user.Group)
		}

		moduleAuthGroups := i.authorityHandler.QueryUserModuleAuthGroup(id)
		for _, val := range moduleAuthGroups {
			view := model.ModuleAuthGroupView{}

			mod, _ := i.moduleHub.FindModule(val.Module)
			view.Module.ID = mod.ID()
			view.Module.Name = mod.Name()
			view.AuthGroup = common_const.GetAuthGroup(val.AuthGroup)

			result.User.ModuleAuthGroup = append(result.User.ModuleAuthGroup, view)
		}

		result.ErrorCode = common_def.Success

		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type userPutRoute struct {
	authorityHandler common.AuthorityHandler
	accountHandler   common.AccountHandler
	moduleHub        common.ModuleHub
}

func (i *userPutRoute) Method() string {
	return common.PUT
}

func (i *userPutRoute) Pattern() string {
	return net.JoinURL(def.URL, def.PutUser)
}

func (i *userPutRoute) Handler() interface{} {
	return i.putHandler
}

func (i *userPutRoute) AuthGroup() int {
	return common_const.UserAuthGroup.ID
}

func (i *userPutRoute) putHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("putHandler")

	result := common_def.UpdateModuleAuthGroupResult{}
	for true {
		_, strID := net.SplitRESTAPI(r.URL.Path)
		id, err := strconv.Atoi(strID)
		if err != nil {
			result.ErrorCode = common_def.Failed
			result.Reason = "非法参数"
			break
		}

		param := &common_def.UpdateModuleAuthGroupParam{}
		err = net.ParseJSONBody(r, param)
		if err != nil {
			result.ErrorCode = common_def.Failed
			result.Reason = "非法参数"
			break
		}

		ok := i.authorityHandler.UpdateUserModuleAuthGroup(id, param.ModuleAuthGroup)
		if ok {
			result.ErrorCode = common_def.Success
		} else {
			result.ErrorCode = common_def.Failed
			result.Reason = "更新用户模块授权信息失败"
		}

		result.ErrorCode = common_def.Success
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}
