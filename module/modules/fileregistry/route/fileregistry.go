package route

import (
	"log"
	"net/http"

	"github.com/muidea/magicCenter/common"
	"github.com/muidea/magicCenter/module/modules/fileregistry/def"
	common_const "github.com/muidea/magicCommon/common"
	"github.com/muidea/magicCommon/foundation/net"
)

// AppendFileRegistryRoute 追加FileRegistry路由
func AppendFileRegistryRoute(routes []common.Route, fileRegistryHandler common.FileRegistryHandler) []common.Route {
	route := createUploadFileRoute(fileRegistryHandler)
	routes = append(routes, route)

	rt := createDownloadFileRoute(fileRegistryHandler)
	routes = append(routes, rt)

	rt = createDeleteFileRoute(fileRegistryHandler)
	routes = append(routes, rt)

	return routes
}

func createUploadFileRoute(fileRegistryHandler common.FileRegistryHandler) common.Route {
	return &uploadFileRoute{fileRegistryHandler: fileRegistryHandler}
}

func createDownloadFileRoute(fileRegistryHandler common.FileRegistryHandler) common.Route {
	return &downloadFileRoute{fileRegistryHandler: fileRegistryHandler}
}

func createDeleteFileRoute(fileRegistryHandler common.FileRegistryHandler) common.Route {
	return &deleteFileRoute{fileRegistryHandler: fileRegistryHandler}
}

type uploadFileRoute struct {
	fileRegistryHandler common.FileRegistryHandler
}

func (i *uploadFileRoute) Method() string {
	return common.POST
}

func (i *uploadFileRoute) Pattern() string {
	return net.JoinURL(def.URL, def.PostFile)
}

func (i *uploadFileRoute) Handler() interface{} {
	return i.uploadFileHandler
}

func (i *uploadFileRoute) AuthGroup() int {
	return common_const.UserAuthGroup.ID
}

func (i *uploadFileRoute) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("uploadFileHandler")

	i.fileRegistryHandler.UploadFile(w, r)
}

type downloadFileRoute struct {
	fileRegistryHandler common.FileRegistryHandler
}

func (i *downloadFileRoute) Method() string {
	return common.GET
}

func (i *downloadFileRoute) Pattern() string {
	return net.JoinURL(def.URL, def.GetFile)
}

func (i *downloadFileRoute) Handler() interface{} {
	return i.downloadFileHandler
}

func (i *downloadFileRoute) AuthGroup() int {
	return common_const.VisitorAuthGroup.ID
}

func (i *downloadFileRoute) downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("downloadFileHandler")

	i.fileRegistryHandler.DownloadFile(w, r)
}

type deleteFileRoute struct {
	fileRegistryHandler common.FileRegistryHandler
}

func (i *deleteFileRoute) Method() string {
	return common.DELETE
}

func (i *deleteFileRoute) Pattern() string {
	return net.JoinURL(def.URL, def.DeleteFile)
}

func (i *deleteFileRoute) Handler() interface{} {
	return i.deleteFileHandler
}

func (i *deleteFileRoute) AuthGroup() int {
	return common_const.UserAuthGroup.ID
}

func (i *deleteFileRoute) deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteFileHandler")

	i.fileRegistryHandler.DeleteFile(w, r)
}
