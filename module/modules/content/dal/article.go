package dal

import (
	"database/sql"
	"fmt"

	"muidea.com/magicCenter/common/dbhelper"
	"muidea.com/magicCenter/common/resource"
	common_const "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/foundation/util"
	"muidea.com/magicCommon/model"
)

func loadArticleID(helper dbhelper.DBHelper) int {
	var maxID sql.NullInt64
	sql := fmt.Sprintf(`select max(id) from content_article`)
	helper.Query(sql)
	defer helper.Finish()

	if helper.Next() {
		helper.GetValue(&maxID)
	}

	return int(maxID.Int64)
}

// QueryAllArticleSummary 查询所有文章摘要
func QueryAllArticleSummary(helper dbhelper.DBHelper) []model.Summary {
	summaryList := []model.Summary{}

	ress := resource.QueryResourceByType(helper, model.ARTICLE)
	for _, v := range ress {
		summary := model.Summary{Unit: model.Unit{ID: v.RId(), Name: v.RName()}, Description: v.RDescription(), Type: v.RType(), CreateDate: v.RCreateDate(), Creater: v.ROwner()}

		for _, r := range v.Relative() {
			summary.Catalog = append(summary.Catalog, *r.CatalogUnit())
		}
		// 如果Catalog没有父分类，则认为其父分类为BuildContentCatalog
		if len(summary.Catalog) == 0 {
			summary.Catalog = append(summary.Catalog, *common_const.SystemContentCatalog.CatalogUnit())
		}

		summaryList = append(summaryList, summary)
	}

	return summaryList
}

// QueryArticles 查询指定文章
func QueryArticles(helper dbhelper.DBHelper, ids []int) []model.Article {
	articleList := []model.Article{}

	if len(ids) == 0 {
		return articleList
	}

	sql := fmt.Sprintf(`select id, title from content_article where id in(%s)`, util.IntArray2Str(ids))
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		article := model.Article{}
		helper.GetValue(&article.ID, &article.Title)

		articleList = append(articleList, article)
	}

	return articleList
}

// QueryArticleByID 查询指定文章
func QueryArticleByID(helper dbhelper.DBHelper, id int) (model.ArticleDetail, bool) {
	ar := model.ArticleDetail{}

	sql := fmt.Sprintf(`select id, title, content, creater, createdate from content_article where id = %d`, id)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&ar.ID, &ar.Title, &ar.Content, &ar.Creater, &ar.CreateDate)
		result = true
	}
	helper.Finish()

	if result {
		ress := resource.QueryRelativeResource(helper, ar.ID, model.ARTICLE)
		for _, r := range ress {
			ar.Catalog = append(ar.Catalog, *r.CatalogUnit())
		}

		// 如果Catalog没有父分类，则认为其父分类为BuildContentCatalog
		if len(ar.Catalog) == 0 {
			ar.Catalog = append(ar.Catalog, *common_const.SystemContentCatalog.CatalogUnit())
		}
	}

	return ar, result
}

// QueryArticleSummaryByCatalog 查询指定分类下的所有文章摘要
func QueryArticleSummaryByCatalog(helper dbhelper.DBHelper, catalog model.CatalogUnit) []model.Summary {
	summaryList := []model.Summary{}
	resList := resource.QueryReferenceResource(helper, catalog.ID, catalog.Type, model.ARTICLE)
	for _, r := range resList {
		summary := model.Summary{Unit: model.Unit{ID: r.RId(), Name: r.RName()}, Description: r.RDescription(), Type: r.RType(), CreateDate: r.RCreateDate(), Creater: r.ROwner()}
		summaryList = append(summaryList, summary)
	}

	for index, value := range summaryList {
		summary := &summaryList[index]
		ress := resource.QueryRelativeResource(helper, value.ID, value.Type)
		for _, r := range ress {
			summary.Catalog = append(summary.Catalog, *r.CatalogUnit())
		}

		// 如果Catalog没有父分类，则认为其父分类为BuildContentCatalog
		if len(summary.Catalog) == 0 {
			summary.Catalog = append(summary.Catalog, *common_const.SystemContentCatalog.CatalogUnit())
		}
	}

	return summaryList
}

// CreateArticle 保存文章
func CreateArticle(helper dbhelper.DBHelper, title, content string, catalogs []model.CatalogUnit, creater int, createDate string) (model.Summary, bool) {
	desc := util.ExtractSummary(content)
	article := model.Summary{Unit: model.Unit{Name: title}, Description: desc, Type: model.ARTICLE, Catalog: catalogs, CreateDate: createDate, Creater: creater}

	id := allocArticleID()
	result := false
	helper.BeginTransaction()
	for {
		// insert
		sql := fmt.Sprintf(`insert into content_article (id, title,content,creater,createdate) values (%d, '%s','%s',%d,'%s')`, id, title, content, creater, createDate)
		_, result = helper.Execute(sql)
		if !result {
			break
		}

		article.ID = id
		res := resource.CreateSimpleRes(article.ID, model.ARTICLE, article.Name, desc, article.CreateDate, article.Creater)
		for _, c := range article.Catalog {
			if c.ID != common_const.SystemContentCatalog.ID && c.Type != model.CATALOG {
				ca, ok := resource.QueryResourceByID(helper, c.ID, c.Type)
				if ok {
					res.AppendRelative(ca)
				} else {
					result = false
					break
				}
			}
		}

		if result {
			result = resource.CreateResource(helper, res, true)
		}

		break
	}

	if result {
		helper.Commit()
	} else {
		helper.Rollback()
	}

	return article, result
}

// SaveArticle 保存文章
func SaveArticle(helper dbhelper.DBHelper, article model.ArticleDetail) (model.Summary, bool) {
	desc := util.ExtractSummary(article.Content)
	summary := model.Summary{Unit: model.Unit{ID: article.ID, Name: article.Title}, Description: desc, Type: model.ARTICLE, Catalog: article.Catalog, CreateDate: article.CreateDate, Creater: article.Creater}
	result := false

	helper.BeginTransaction()
	for {
		// modify
		sql := fmt.Sprintf(`update content_article set title ='%s', content ='%s', creater =%d, createdate ='%s' where id=%d`, article.Title, article.Content, article.Creater, article.CreateDate, article.ID)
		_, result = helper.Execute(sql)

		if result {
			res, ok := resource.QueryResourceByID(helper, article.ID, model.ARTICLE)
			if !ok {
				result = false
				break
			}

			res.UpdateName(article.Title)
			res.UpdateDescription(desc)

			res.ResetRelative()
			for _, c := range article.Catalog {
				if c.ID != common_const.SystemContentCatalog.ID && c.Type != model.CATALOG {
					ca, ok := resource.QueryResourceByID(helper, c.ID, c.Type)
					if ok {
						res.AppendRelative(ca)
					} else {
						result = false
						break
					}
				}
			}

			if result {
				result = resource.SaveResource(helper, res, true)
			}

			break
		}

		break
	}

	if result {
		helper.Commit()
	} else {
		helper.Rollback()
	}

	return summary, result
}

// DeleteArticle 删除文章
func DeleteArticle(helper dbhelper.DBHelper, id int) bool {
	result := false
	helper.BeginTransaction()

	for {
		sql := fmt.Sprintf(`delete from content_article where id=%d`, id)

		_, result = helper.Execute(sql)
		if result {
			res, ok := resource.QueryResourceByID(helper, id, model.ARTICLE)
			if ok {
				result = resource.DeleteResource(helper, res, true)
			} else {
				result = ok
			}
		}

		break
	}

	if result {
		helper.Commit()
	} else {
		helper.Rollback()
	}

	return result
}
