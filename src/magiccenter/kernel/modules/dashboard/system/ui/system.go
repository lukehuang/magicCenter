package ui

import (
	"html/template"
	"log"
	"magiccenter/common"
	"magiccenter/system"
	"net/http"
)

const passwordMark = "********"

// SystemSettingView 系统设置视图
type SystemSettingView struct {
	SystemInfo common.SystemInfo
}

// SystemSettingViewHandler 系统设置视图处理器
func SystemSettingViewHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("SystemSettingViewHandler")

	w.Header().Set("content-type", "text/html")
	w.Header().Set("charset", "utf-8")

	t, err := template.ParseFiles("template/html/dashboard/system/setting.html")
	if err != nil {
		panic("parse files failed")
	}

	view := SystemSettingView{}
	configuration := system.GetConfiguration()
	view.SystemInfo = configuration.GetSystemInfo()
	view.SystemInfo.MailPassword = passwordMark

	t.Execute(w, view)
}
