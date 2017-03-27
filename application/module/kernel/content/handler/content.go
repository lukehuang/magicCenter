package handler

import (
	"muidea.com/magicCenter/application/common"
	"muidea.com/magicCenter/application/common/dbhelper"
	"muidea.com/magicCenter/application/common/model"
)

// CreateContentHandler 新建ContentHandler
func CreateContentHandler() common.ContentHandler {
	dbhelper, _ := dbhelper.NewHelper()
	i := impl{
		articleHandler: articleActionHandler{dbhelper: dbhelper},
		catalogHandler: catalogActionHandler{dbhelper: dbhelper},
		linkHandler:    linkActionHandler{dbhelper: dbhelper},
		mediaHandler:   mediaActionHandler{dbhelper: dbhelper}}

	return &i
}

type impl struct {
	articleHandler articleActionHandler
	catalogHandler catalogActionHandler
	linkHandler    linkActionHandler
	mediaHandler   mediaActionHandler
}

func (i *impl) GetAllArticle() []model.Summary {
	return i.articleHandler.getAllArticleSummary()
}

func (i *impl) GetArticleByID(id int) (model.ArticleDetail, bool) {
	return i.articleHandler.findArticleByID(id)
}

func (i *impl) GetArticleByCatalog(catalog int) []model.Summary {
	return i.articleHandler.findArticleByCatalog(catalog)
}

func (i *impl) CreateArticle(title, content, createDate string, catalog []int, author int) (model.Summary, bool) {
	return i.articleHandler.createArticle(title, content, createDate, catalog, author)
}

func (i *impl) SaveArticle(article model.ArticleDetail) (model.Summary, bool) {
	return i.articleHandler.saveArticle(article)
}

func (i *impl) DestroyArticle(id int) bool {
	return i.articleHandler.destroyArticle(id)
}

func (i *impl) GetAllCatalog() []model.Summary {
	return i.catalogHandler.getAllCatalog()
}

func (i *impl) GetCatalogByID(id int) (model.CatalogDetail, bool) {
	return i.catalogHandler.findCatalogByID(id)
}

func (i *impl) GetCatalogByCatalog(id int) []model.Summary {
	return i.catalogHandler.findCatalogByCatalog(id)
}

func (i *impl) CreateCatalog(name, description, createdate string, parent []int, author int) (model.Summary, bool) {
	return i.catalogHandler.createCatalog(name, description, createdate, parent, author)
}

func (i *impl) SaveCatalog(catalog model.CatalogDetail) (model.Summary, bool) {
	return i.catalogHandler.saveCatalog(catalog)
}

func (i *impl) DestroyCatalog(id int) bool {
	return i.catalogHandler.destroyCatalog(id)
}

func (i *impl) GetAllLink() []model.Summary {
	return i.linkHandler.getAllLink()
}

func (i *impl) GetLinkByID(id int) (model.LinkDetail, bool) {
	return i.linkHandler.findLinkByID(id)
}

func (i *impl) GetLinkByCatalog(catalog int) []model.Summary {
	return i.linkHandler.findLinkByCatalog(catalog)
}

func (i *impl) CreateLink(name, url, logo, createdate string, catalog []int, author int) (model.Summary, bool) {
	return i.linkHandler.createLink(name, url, logo, createdate, catalog, author)
}

func (i *impl) SaveLink(link model.LinkDetail) (model.Summary, bool) {
	return i.linkHandler.saveLink(link)
}

func (i *impl) DestroyLink(id int) bool {
	return i.linkHandler.destroyLink(id)
}

func (i *impl) GetAllMedia() []model.Summary {
	return i.mediaHandler.getAllMedia()
}

func (i *impl) GetMediaByID(id int) (model.MediaDetail, bool) {
	return i.mediaHandler.findMediaByID(id)
}

func (i *impl) GetMediaByCatalog(catalog int) []model.Summary {
	return i.mediaHandler.findMediaByCatalog(catalog)
}

func (i *impl) CreateMedia(name, url, desc, createdate string, catalog []int, author int) (model.Summary, bool) {
	return i.mediaHandler.createMedia(name, url, desc, createdate, catalog, author)
}

func (i *impl) SaveMedia(media model.MediaDetail) (model.Summary, bool) {
	return i.mediaHandler.saveMedia(media)
}

func (i *impl) DestroyMedia(id int) bool {
	return i.mediaHandler.destroyMedia(id)
}