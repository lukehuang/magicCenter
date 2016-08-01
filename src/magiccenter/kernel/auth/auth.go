package auth

import (
	"log"
	"magiccenter/configuration"
	"magiccenter/kernel/modules/account/bll"
	"magiccenter/kernel/modules/account/model"
	"magiccenter/router"
	"magiccenter/session"
	"martini"
	"net/http"
)

// AdminAuthVerify 管理权限校验器
func AdminAuthVerify() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request) bool {
		authID, found := configuration.GetOption(configuration.AuthorithID)
		if !found {
			panic("unexpected, can't fetch authorith id")
		}

		result := false
		session := session.GetSession(res, req)
		user, found := session.GetOption(authID)
		if found {
			for _, gid := range user.(model.UserDetail).Groups {
				group, found := bll.QueryGroupByID(gid)
				if found && group.AdminGroup() {
					result = true
				}
			}
		}

		return result
	}
}

// LoginAuthVerify 登陆权限校验器
func LoginAuthVerify() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request) bool {
		authID, found := configuration.GetOption(configuration.AuthorithID)
		if !found {
			panic("unexpected, can't fetch authorith id")
		}

		session := session.GetSession(res, req)
		_, found = session.GetOption(authID)
		return found
	}
}

// Authority 权限校验处理器
// 用于在路由过程中进行权限校验
func Authority() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {

		if !router.VerifyAuthority(res, req) {
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		c.Next()
	}
}
