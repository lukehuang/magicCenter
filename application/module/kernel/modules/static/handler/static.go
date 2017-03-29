package handler

import (
	"net/http"
	"os"
	"path"

	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/module/kernel/modules/static/util"
)

// CreateStaticHandler 新建StaticHandler
func CreateStaticHandler(rootPath string) common.StaticHandler {
	i := impl{rootPath: rootPath}

	return &i
}

type impl struct {
	rootPath string
}

func (i *impl) HandleResource(basePath string, w http.ResponseWriter, r *http.Request) {
	fullPath := util.MergePath(i.rootPath, basePath, r.URL.Path)
	_, err := os.Stat(fullPath)
	if err == nil || os.IsExist(err) {
		filePath, fileName := path.Split(fullPath)
		dir := http.Dir(filePath)
		f, err := dir.Open(fileName)
		if err != nil {
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil || fi.IsDir() {
			return
		}

		http.ServeContent(w, r, fullPath, fi.ModTime(), f)
	} else {
		fullPath = path.Join(i.rootPath, "404-res.html")
	}
}
