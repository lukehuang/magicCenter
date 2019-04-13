package dal

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/muidea/magicCenter/common/dbhelper"
	"github.com/muidea/magicCenter/common/resource"
	common_const "github.com/muidea/magicCommon/common"
	"github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/util"
	"github.com/muidea/magicCommon/model"
)

func loadMediaID(helper dbhelper.DBHelper) int {
	var maxID sql.NullInt64
	sql := fmt.Sprintf(`select max(id) from content_media`)
	helper.Query(sql)
	defer helper.Finish()

	if helper.Next() {
		helper.GetValue(&maxID)
	}

	return int(maxID.Int64)
}

// QueryAllMedia 查询所有图像
func QueryAllMedia(helper dbhelper.DBHelper, filter *def.Filter) ([]model.Summary, int) {
	summaryList := []model.Summary{}

	ress, resCount := resource.QueryResourceByType(helper, model.MEDIA, filter)
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

	return summaryList, resCount
}

// QueryMedias 查询指定文章
func QueryMedias(helper dbhelper.DBHelper, ids []int) []model.Media {
	mediaList := []model.Media{}

	if len(ids) == 0 {
		return mediaList
	}

	sql := fmt.Sprintf(`select id, name from content_media where id in(%s)`, util.IntArray2Str(ids))
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		media := model.Media{}
		helper.GetValue(&media.ID, &media.Name)

		mediaList = append(mediaList, media)
	}

	return mediaList
}

// QueryMediaByCatalog 查询指定分类的图像
func QueryMediaByCatalog(helper dbhelper.DBHelper, catalog model.CatalogUnit, filter *def.Filter) ([]model.Summary, int) {
	summaryList := []model.Summary{}

	resList, resCount := resource.QueryReferenceResource(helper, catalog.ID, catalog.Type, model.MEDIA, filter)
	for _, r := range resList {
		summary := model.Summary{Unit: model.Unit{ID: r.RId(), Name: r.RName()}, Description: r.RDescription(), Type: r.RType(), CreateDate: r.RCreateDate(), Creater: r.ROwner()}
		summaryList = append(summaryList, summary)
	}

	for index, value := range summaryList {
		summary := &summaryList[index]
		ress, _ := resource.QueryRelativeResource(helper, value.ID, value.Type, nil)
		for _, r := range ress {
			summary.Catalog = append(summary.Catalog, *r.CatalogUnit())
		}

		// 如果Catalog没有父分类，则认为其父分类为BuildContentCatalog
		if len(summary.Catalog) == 0 {
			summary.Catalog = append(summary.Catalog, *common_const.SystemContentCatalog.CatalogUnit())
		}
	}

	return summaryList, resCount
}

// QueryMediaByID 查询指定的图像
func QueryMediaByID(helper dbhelper.DBHelper, id int) (model.MediaDetail, bool) {
	media := model.MediaDetail{}

	sql := fmt.Sprintf(`select id, name, description, fileToken, createdate, creater, expiration from content_media where id = %d`, id)
	helper.Query(sql)

	result := false
	if helper.Next() {
		helper.GetValue(&media.ID, &media.Name, &media.Description, &media.FileToken, &media.CreateDate, &media.Creater, &media.Expiration)
		result = true
	}
	helper.Finish()

	if result {
		ress, _ := resource.QueryRelativeResource(helper, id, model.MEDIA, nil)
		for _, r := range ress {
			media.Catalog = append(media.Catalog, *r.CatalogUnit())
		}

		// 如果Catalog没有父分类，则认为其父分类为BuildContentCatalog
		if len(media.Catalog) == 0 {
			media.Catalog = append(media.Catalog, *common_const.SystemContentCatalog.CatalogUnit())
		}
	}

	return media, result
}

// DeleteMediaByID 删除图像
func DeleteMediaByID(helper dbhelper.DBHelper, id int) bool {
	result := false
	helper.BeginTransaction()

	for {
		sql := fmt.Sprintf(`delete from content_media where id =%d`, id)
		_, result = helper.Execute(sql)
		if result {
			res, ok := resource.QueryResourceByID(helper, id, model.MEDIA)
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

func createSingle(helper dbhelper.DBHelper, name, description, fileToken, createDate string, expiration, creater int, catalogs []model.CatalogUnit) (model.Summary, bool) {
	media := model.Summary{Unit: model.Unit{Name: name}, Description: description, Type: model.MEDIA, Catalog: catalogs, CreateDate: createDate, Creater: creater}

	id := allocMediaID()
	result := false
	for {
		// insert
		sql := fmt.Sprintf(`insert into content_media (id, name, description, fileToken, createdate, creater, expiration) values (%d, '%s','%s','%s','%s',%d,%d)`, id, name, description, fileToken, createDate, creater, expiration)
		_, result = helper.Execute(sql)
		if !result {
			break
		}

		media.ID = id
		res := resource.CreateSimpleRes(media.ID, model.MEDIA, media.Name, media.Description, media.CreateDate, media.Creater)
		constCatalogUnit := common_const.SystemContentCatalog.CatalogUnit()
		for _, c := range media.Catalog {
			if c.ID == constCatalogUnit.ID && c.Type == constCatalogUnit.Type {
				continue
			}

			ca, ok := resource.QueryResourceByID(helper, c.ID, c.Type)
			if ok {
				res.AppendRelative(ca)
			} else {
				result = false
				break
			}
		}

		if result {
			result = resource.CreateResource(helper, res, true)
		} else {
			media.ID = -1
		}

		break
	}

	return media, result
}

// CreateMedia 新建文件
func CreateMedia(helper dbhelper.DBHelper, name, description, fileToken, createDate string, expiration, creater int, catalogs []model.CatalogUnit) (model.Summary, bool) {
	media := model.Summary{Unit: model.Unit{Name: name}, Description: description, Type: model.MEDIA, Catalog: catalogs, CreateDate: createDate, Creater: creater}
	result := false
	helper.BeginTransaction()

	media, result = createSingle(helper, name, description, fileToken, createDate, expiration, creater, catalogs)

	if result {
		helper.Commit()
	} else {
		helper.Rollback()
	}

	return media, result
}

// BatchCreateMedia 批量新建文件
func BatchCreateMedia(helper dbhelper.DBHelper, medias []model.MediaItem, createDate string, creater int) ([]model.Summary, bool) {
	summaryList := []model.Summary{}

	for _, val := range medias {
		summary, _ := CreateMedia(helper, val.Name, val.Description, val.FileToken, createDate, val.Expiration, creater, val.Catalog)
		summaryList = append(summaryList, summary)
	}

	return summaryList, true
}

// SaveMedia 保存文件
func SaveMedia(helper dbhelper.DBHelper, media model.MediaDetail) (model.Summary, bool) {
	summary := model.Summary{Unit: model.Unit{ID: media.ID, Name: media.Name}, Type: model.MEDIA, Catalog: media.Catalog, CreateDate: media.CreateDate, Creater: media.Creater}
	result := false
	helper.BeginTransaction()
	for {
		// modify
		sql := fmt.Sprintf(`update content_media set name='%s', description='%s', fileToken ='%s', createdate='%s', creater=%d where id=%d`, media.Name, media.Description, media.FileToken, media.CreateDate, media.Creater, media.ID)
		_, result = helper.Execute(sql)
		if result {
			res, ok := resource.QueryResourceByID(helper, media.ID, model.MEDIA)
			if !ok {
				result = false
				break
			}

			res.UpdateName(media.Name)
			res.UpdateDescription(media.Description)
			res.ResetRelative()
			constCatalogUnit := common_const.SystemContentCatalog.CatalogUnit()
			for _, c := range media.Catalog {
				if c.ID == constCatalogUnit.ID && c.Type == constCatalogUnit.Type {
					continue
				}

				ca, ok := resource.QueryResourceByID(helper, c.ID, c.Type)
				if ok {
					res.AppendRelative(ca)
				} else {
					result = false
					break
				}
			}

			if result {
				result = resource.SaveResource(helper, res, true)
			}
		} else {
			log.Printf("update content_media failed, sql:%s", sql)
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

// LoadMediaExpiration 加载所有Media有效期
func LoadMediaExpiration(helper dbhelper.DBHelper) map[int]time.Time {
	expirationMap := map[int]time.Time{}

	sql := fmt.Sprintf(`select id, createdate, expiration from content_media`)
	helper.Query(sql)
	defer helper.Finish()

	for helper.Next() {
		id := -1
		createDate := ""
		expiration := -1
		helper.GetValue(&id, &createDate, &expiration)

		layout := "2006-01-02 15:04:05"
		cd, err := time.Parse(layout, createDate)
		if err != nil {
			cd, _ = time.Parse(layout, "1970-01-01 00:00:01")
		}

		expirationMap[id] = cd.Add(time.Duration(expiration*24) * time.Hour)
	}

	return expirationMap
}
