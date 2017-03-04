package bll

import (
	"magiccenter/common/model"
	"magiccenter/kernel/modules/content/dal"
	"magiccenter/system"
)

// QueryAllCatalog 查询全部分类
func QueryAllCatalog() []model.Catalog {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryAllCatalog(helper)
}

// QueryAllCatalogDetail 查询全部分类详情
func QueryAllCatalogDetail() []model.CatalogDetail {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryAllCatalogDetail(helper)
}

// QueryCatalogByID 查询指定分类
func QueryCatalogByID(id int) (model.CatalogDetail, bool) {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	catalog, result := dal.QueryCatalogByID(helper, id)
	return catalog, result
}

// QueryAvalibleParentCatalog 查询可用父类
func QueryAvalibleParentCatalog(id int) []model.Catalog {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QueryAvalibleParentCatalog(helper, id)
}

// QuerySubCatalog 查询子类
func QuerySubCatalog(id int) []model.Catalog {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.QuerySubCatalog(helper, id)
}

// DeleteCatalog 删除分类
func DeleteCatalog(id int) bool {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.DeleteCatalog(helper, id)
}

// CreateCatalog 新建分类
func CreateCatalog(name string, uID int, parents []int) (model.CatalogDetail, bool) {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.CreateCatalog(helper, name, parents, uID)
}

// UpdateCatalog 更新分类
func UpdateCatalog(catalog model.CatalogDetail) (model.CatalogDetail, bool) {
	helper, err := system.GetDBHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	return dal.SaveCatalog(helper, catalog)
}