package route

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/common/model"
	"muidea.com/magicCenter/application/module/kernel/modules/api/def"
	"muidea.com/magicCenter/foundation/net"
	"muidea.com/magicCenter/foundation/util"
)

// AppendCatalogRoute 追加User Route
func AppendCatalogRoute(routes []common.Route, modHub common.ModuleHub, sessionRegistry common.SessionRegistry) []common.Route {
	rt, _ := CreateGetCatalogByIDRoute(modHub)
	routes = append(routes, rt)

	rt, _ = CreateGetCatalogListRoute(modHub)
	routes = append(routes, rt)

	rt, _ = CreateCreateCatalogRoute(modHub, sessionRegistry)
	routes = append(routes, rt)

	rt, _ = CreateUpdateCatalogRoute(modHub, sessionRegistry)
	routes = append(routes, rt)

	rt, _ = CreateDestroyCatalogRoute(modHub, sessionRegistry)
	routes = append(routes, rt)

	return routes
}

// CreateGetCatalogByIDRoute 新建GetCatalog Route
func CreateGetCatalogByIDRoute(modHub common.ModuleHub) (common.Route, bool) {
	mod, found := modHub.FindModule(common.CotentModuleID)
	if !found {
		return nil, false
	}

	endPoint := mod.EndPoint()
	switch endPoint.(type) {
	case common.ContentHandler:
		i := catalogGetByIDRoute{contentHandler: endPoint.(common.ContentHandler)}
		return &i, true
	}

	return nil, false
}

// CreateGetCatalogListRoute 新建GetAllCatalog Route
func CreateGetCatalogListRoute(modHub common.ModuleHub) (common.Route, bool) {
	mod, found := modHub.FindModule(common.CotentModuleID)
	if !found {
		return nil, false
	}

	endPoint := mod.EndPoint()
	switch endPoint.(type) {
	case common.ContentHandler:
		i := catalogGetListRoute{contentHandler: endPoint.(common.ContentHandler)}
		return &i, true
	}

	return nil, false
}

// CreateCreateCatalogRoute 新建CreateCatalogRoute Route
func CreateCreateCatalogRoute(modHub common.ModuleHub, sessionRegistry common.SessionRegistry) (common.Route, bool) {
	mod, found := modHub.FindModule(common.CotentModuleID)
	if !found {
		return nil, false
	}

	endPoint := mod.EndPoint()
	switch endPoint.(type) {
	case common.ContentHandler:
		i := catalogCreateRoute{contentHandler: endPoint.(common.ContentHandler), sessionRegistry: sessionRegistry}
		return &i, true
	}

	return nil, false
}

// CreateUpdateCatalogRoute UpdateCatalogRoute Route
func CreateUpdateCatalogRoute(modHub common.ModuleHub, sessionRegistry common.SessionRegistry) (common.Route, bool) {
	mod, found := modHub.FindModule(common.CotentModuleID)
	if !found {
		return nil, false
	}

	endPoint := mod.EndPoint()
	switch endPoint.(type) {
	case common.ContentHandler:
		i := catalogUpdateRoute{contentHandler: endPoint.(common.ContentHandler), sessionRegistry: sessionRegistry}
		return &i, true
	}

	return nil, false
}

// CreateDestroyCatalogRoute DestroyCatalogRoute Route
func CreateDestroyCatalogRoute(modHub common.ModuleHub, sessionRegistry common.SessionRegistry) (common.Route, bool) {
	mod, found := modHub.FindModule(common.CotentModuleID)
	if !found {
		return nil, false
	}

	endPoint := mod.EndPoint()
	switch endPoint.(type) {
	case common.ContentHandler:
		i := catalogDestroyRoute{contentHandler: endPoint.(common.ContentHandler), sessionRegistry: sessionRegistry}
		return &i, true
	}

	return nil, false
}

type catalogGetByIDRoute struct {
	contentHandler common.ContentHandler
}

type catalogGetByIDResult struct {
	common.Result
	Catalog model.CatalogDetail
}

func (i *catalogGetByIDRoute) Method() string {
	return common.GET
}

func (i *catalogGetByIDRoute) Pattern() string {
	return net.JoinURL(def.URL, "content/catalog/[0-9]+/")
}

func (i *catalogGetByIDRoute) Handler() interface{} {
	return i.getCatalogHandler
}

func (i *catalogGetByIDRoute) getCatalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("getCatalogHandler")

	result := catalogGetByIDResult{}
	_, value := net.SplitRESTAPI(r.URL.Path)
	for true {
		id, err := strconv.Atoi(value)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效参数"
			break
		}

		catalog, ok := i.contentHandler.GetCatalogByID(id)
		if ok {
			result.Catalog = catalog
			result.ErrCode = 0
		} else {
			result.ErrCode = 1
			result.Reason = "对象不存在"
		}
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type catalogGetListRoute struct {
	contentHandler common.ContentHandler
}

type catalogGetListResult struct {
	common.Result
	Catalog []model.Summary
}

func (i *catalogGetListRoute) Method() string {
	return common.GET
}

func (i *catalogGetListRoute) Pattern() string {
	return net.JoinURL(def.URL, "content/catalog/")
}

func (i *catalogGetListRoute) Handler() interface{} {
	return i.getCatalogListHandler
}

func (i *catalogGetListRoute) getCatalogListHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("getCatalogListHandler")

	result := catalogGetListResult{}
	for true {
		catalog := r.URL.Query()["catalog"]
		if len(catalog) < 1 {
			result.Catalog = i.contentHandler.GetAllCatalog()
			result.ErrCode = 0
			break
		}

		id, err := strconv.Atoi(catalog[0])
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效参数"
			break
		}

		result.Catalog = i.contentHandler.GetCatalogByCatalog(id)
		result.ErrCode = 0
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type catalogCreateRoute struct {
	contentHandler  common.ContentHandler
	sessionRegistry common.SessionRegistry
}

type catalogCreateResult struct {
	common.Result
	Catalog model.Summary
}

func (i *catalogCreateRoute) Method() string {
	return common.POST
}

func (i *catalogCreateRoute) Pattern() string {
	return net.JoinURL(def.URL, "content/catalog/")
}

func (i *catalogCreateRoute) Handler() interface{} {
	return i.createCatalogHandler
}

func (i *catalogCreateRoute) createCatalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("createCatalogHandler")

	session := i.sessionRegistry.GetSession(w, r)
	result := catalogCreateResult{}
	for true {
		user, found := session.GetAccount()
		if !found {
			result.ErrCode = 1
			result.Reason = "无效权限"
			break
		}

		r.ParseForm()

		name := r.FormValue("catalog-name")
		description := r.FormValue("catalog-description")
		createdate := time.Now().Format("2006-01-02 15:04:05")
		catalogs, _ := util.Str2IntArray(r.FormValue("catalog-catalog"))
		catalog, ok := i.contentHandler.CreateCatalog(name, description, createdate, catalogs, user.ID)
		if !ok {
			result.ErrCode = 1
			result.Reason = "新建失败"
			break
		}
		result.ErrCode = 0
		result.Catalog = catalog
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type catalogUpdateRoute struct {
	contentHandler  common.ContentHandler
	sessionRegistry common.SessionRegistry
}

type catalogUpdateResult struct {
	common.Result
	Catalog model.Summary
}

func (i *catalogUpdateRoute) Method() string {
	return common.PUT
}

func (i *catalogUpdateRoute) Pattern() string {
	return net.JoinURL(def.URL, "content/catalog/[0-9]+/")
}

func (i *catalogUpdateRoute) Handler() interface{} {
	return i.updateCatalogHandler
}

func (i *catalogUpdateRoute) updateCatalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("updateCatalogHandler")

	session := i.sessionRegistry.GetSession(w, r)
	result := catalogCreateResult{}
	_, value := net.SplitRESTAPI(r.URL.Path)
	for true {
		id, err := strconv.Atoi(value)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效参数"
			break
		}

		user, found := session.GetAccount()
		if !found {
			result.ErrCode = 1
			result.Reason = "无效权限"
			break
		}

		r.ParseForm()
		catalog := model.CatalogDetail{}
		catalog.ID = id
		catalog.Name = r.FormValue("catalog-name")
		catalog.Description = r.FormValue("catalog-description")
		catalog.CreateDate = time.Now().Format("2006-01-02 15:04:05")
		catalog.Catalog, _ = util.Str2IntArray(r.FormValue("catalog-catalog"))
		catalog.Creater = user.ID
		summmary, ok := i.contentHandler.SaveCatalog(catalog)
		if !ok {
			result.ErrCode = 1
			result.Reason = "更新失败"
			break
		}
		result.ErrCode = 0
		result.Catalog = summmary
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}

type catalogDestroyRoute struct {
	contentHandler  common.ContentHandler
	sessionRegistry common.SessionRegistry
}

type catalogDestroyResult struct {
	common.Result
}

func (i *catalogDestroyRoute) Method() string {
	return common.DELETE
}

func (i *catalogDestroyRoute) Pattern() string {
	return net.JoinURL(def.URL, "content/catalog/[0-9]+/")
}

func (i *catalogDestroyRoute) Handler() interface{} {
	return i.deleteCatalogHandler
}

func (i *catalogDestroyRoute) deleteCatalogHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteCatalogHandler")

	session := i.sessionRegistry.GetSession(w, r)
	result := catalogCreateResult{}
	_, value := net.SplitRESTAPI(r.URL.Path)
	for true {
		id, err := strconv.Atoi(value)
		if err != nil {
			result.ErrCode = 1
			result.Reason = "无效参数"
			break
		}
		_, found := session.GetAccount()
		if !found {
			result.ErrCode = 1
			result.Reason = "无效权限"
			break
		}

		ok := i.contentHandler.DestroyCatalog(id)
		if !ok {
			result.ErrCode = 1
			result.Reason = "删除失败"
			break
		}
		result.ErrCode = 0
		break
	}

	b, err := json.Marshal(result)
	if err != nil {
		panic("json.Marshal, failed, err:" + err.Error())
	}

	w.Write(b)
}
