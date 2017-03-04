package bll

import (
	"magiccenter/common/model"
	"magiccenter/kernel/modules/content/dal"
	"magiccenter/system"
)

// QueryAllMedia 查询全部图像
func QueryAllMedia() []model.MediaDetail {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryAllMedia(helper)
}

// QueryMediaByID 查询指定图像
func QueryMediaByID(id int) (model.MediaDetail, bool) {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryMediaByID(helper, id)
}

// DeleteMediaByID 删除图像
func DeleteMediaByID(id int) bool {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.DeleteMediaByID(helper, id)
}

// QueryMediaByCatalog 查询指定分类的图像
func QueryMediaByCatalog(id int) []model.MediaDetail {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryMediaByCatalog(helper, id)
}

// QueryMediaByRang 查询指定范围图像
func QueryMediaByRang(begin, offset int) []model.MediaDetail {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryMediaByRang(helper, begin, offset)
}

// CreateMedia 新建图像
func CreateMedia(name, url, mediaType, desc string, uID int, catalogs []int) (model.MediaDetail, bool) {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.CreateMedia(helper, name, url, mediaType, desc, uID, catalogs)
}

// SaveMedia 保存图像
func SaveMedia(media model.MediaDetail) (model.MediaDetail, bool) {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.SaveMedia(helper, media)
}