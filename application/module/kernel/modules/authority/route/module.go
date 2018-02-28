package route

import (
	"encoding/json"
	"log"
	"net/http"

	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/common/model"
	"muidea.com/magicCenter/application/module/kernel/modules/authority/def"
	"muidea.com/magicCenter/foundation/net"
)

// CreateGetModuleACLRoute 新建ModuleACLGetRoute
func CreateGetModuleACLRoute(authorityHandler common.AuthorityHandler) common.Route {
	i := moduleGetACLRoute{authorityHandler: authorityHandler}
	return &i
}

// CreateGetModuleUserAuthGroupRoute 新建ModuleUserGetRoute
func CreateGetModuleUserAuthGroupRoute(authorityHandler common.AuthorityHandler) common.Route {
	i := moduleGetUserAuthGroupRoute{authorityHandler: authorityHandler}
	return &i
}

// CreatePutModuleUserAuthGroupRoute 新建PutModuleUserRoute
func CreatePutModuleUserAuthGroupRoute(authorityHandler common.AuthorityHandler) common.Route {
	i := modulePutUserAuthGroupRoute{authorityHandler: authorityHandler}
	return &i
}

type moduleGetACLRoute struct {
	authorityHandler common.AuthorityHandler
}

type moduleGetACLResult struct {
	common.Result
	Module string      `json:"module"`
	ACL    []model.ACL `json:"acl"`
}

func (i *moduleGetACLRoute) Method() string {
	return common.GET
}

func (i *moduleGetACLRoute) Pattern() string {
	return net.JoinURL(def.URL, def.GetModuleACL)
}

func (i *moduleGetACLRoute) Handler() interface{} {
	return i.getModuleACLHandler
}

func (i *moduleGetACLRoute) AuthGroup() int {
	return common.MaintainerAuthGroup.ID
}

func (i *moduleGetACLRoute) getModuleACLHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("getModuleACLHandler")
	result := moduleGetACLResult{}

	for true {
		_, id := net.SplitRESTAPI(r.URL.Path)
		result.Module = id

		result.ACL = i.authorityHandler.QueryACLByModule(id)
		result.ErrorCode = common.Success

		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type moduleGetUserAuthGroupRoute struct {
	authorityHandler common.AuthorityHandler
}

type moduleGetUserAuthGroupResult struct {
	common.Result
	model.ModuleUserAuthGroupInfo
}

func (i *moduleGetUserAuthGroupRoute) Method() string {
	return common.GET
}

func (i *moduleGetUserAuthGroupRoute) Pattern() string {
	return net.JoinURL(def.URL, def.GetModuleUserAuthGroup)
}

func (i *moduleGetUserAuthGroupRoute) Handler() interface{} {
	return i.getModuleUserAuthGroupHandler
}

func (i *moduleGetUserAuthGroupRoute) AuthGroup() int {
	return common.MaintainerAuthGroup.ID
}

func (i *moduleGetUserAuthGroupRoute) getModuleUserAuthGroupHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("getModuleUserAuthGroupHandler")

	result := moduleGetUserAuthGroupResult{}
	for true {
		_, id := net.SplitRESTAPI(r.URL.Path)
		info := i.authorityHandler.QueryModuleUserAuthGroup(id)
		result.Module = id
		result.UserAuthGroups = info.UserAuthGroups
		result.ErrorCode = common.Success

		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type modulePutUserAuthGroupRoute struct {
	authorityHandler common.AuthorityHandler
}

type modulePutUserAuthGroupParam struct {
	UserAuthGroup []model.UserAuthGroup
}

type modulePutUserAuthGroupResult struct {
	common.Result
}

func (i *modulePutUserAuthGroupRoute) Method() string {
	return common.PUT
}

func (i *modulePutUserAuthGroupRoute) Pattern() string {
	return net.JoinURL(def.URL, def.PutModuleUserAuthGroup)
}

func (i *modulePutUserAuthGroupRoute) Handler() interface{} {
	return i.putModuleUserAuthGroupHandler
}

func (i *modulePutUserAuthGroupRoute) AuthGroup() int {
	return common.MaintainerAuthGroup.ID
}

func (i *modulePutUserAuthGroupRoute) putModuleUserAuthGroupHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("putModuleUserAuthGroupHandler")

	result := modulePutUserAuthGroupResult{}
	for true {
		_, id := net.SplitRESTAPI(r.URL.Path)

		param := &modulePutUserAuthGroupParam{}
		err := net.ParsePostJSON(r, param)
		if err != nil {
			result.ErrorCode = common.Failed
			result.Reason = "非法参数"
			break
		}

		ok := i.authorityHandler.UpdateModuleUserAuthGroup(id, param.UserAuthGroup)
		if ok {
			result.ErrorCode = common.Success
		} else {
			result.ErrorCode = common.Failed
			result.Reason = "更新模块用户信息失败"
		}

		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}
