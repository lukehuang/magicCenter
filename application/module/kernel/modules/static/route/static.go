package route

import (
	"log"
	"net/http"

	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/module/kernel/modules/static/def"
	"muidea.com/magicCenter/foundation/net"
)

// CreateStaticViewRoute 新建静态视图路由
func CreateStaticViewRoute(staticHandler common.StaticHandler) common.Route {
	i := &staticViewRoute{staticHandler: staticHandler}

	return i
}

// CreateStaticResRoute 新建静态资源路由
func CreateStaticResRoute(staticHandler common.StaticHandler) common.Route {
	i := &staticResRoute{staticHandler: staticHandler}

	return i
}

type staticViewRoute struct {
	staticHandler common.StaticHandler
}

func (i *staticViewRoute) Method() string {
	return common.GET
}

func (i *staticViewRoute) Pattern() string {
	return net.JoinURL(def.URL, "**.html")
}

func (i *staticViewRoute) Handler() interface{} {
	return i.getStaticViewHandler
}

func (i *staticViewRoute) getStaticViewHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getStaticViewHandler, path:%s", r.URL.Path)

	i.staticHandler.HandleView("default", w, r)
}

type staticResRoute struct {
	staticHandler common.StaticHandler
}

func (i *staticResRoute) Method() string {
	return common.GET
}

func (i *staticResRoute) Pattern() string {
	return net.JoinURL(def.URL, "**")
}

func (i *staticResRoute) Handler() interface{} {
	return i.getStaticResHandler
}

func (i *staticResRoute) getStaticResHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("getStaticResHandler, path:%s", r.URL.Path)
	i.staticHandler.HandleResource("default", w, r)
}