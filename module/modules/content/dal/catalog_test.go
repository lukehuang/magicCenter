package dal

import (
	"log"
	"testing"

	"github.com/muidea/magicCenter/common/dbhelper"
	"github.com/muidea/magicCommon/model"
)

func TestCatalog(t *testing.T) {
	log.Print("------------------TestCatalog--------------------")

	helper, err := dbhelper.NewHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	ca := model.CatalogDetail{}
	ca.ID = 12
	ca.Name = "testCatalog"
	ca.Creater = 3
	ca.Catalog = append(ca.Catalog, model.CatalogUnit{ID: 10, Type: "catalog"})
	catalog, ret := SaveCatalog(helper, ca, true)
	if !ret {
		t.Error("SaveCatalog failed")
		return
	}

	ca, found := QueryCatalogByID(helper, catalog.ID)
	if !found {
		t.Error("QueryCatalogByID failed")
	}
	if ca.Creater != 3 {
		t.Error("QueryCatalogByID failed")
	}

	ca.Catalog = append(ca.Catalog, model.CatalogUnit{ID: 8, Type: "catalog"})
	ca.Catalog = append(ca.Catalog, model.CatalogUnit{ID: 9, Type: "catalog"})

	catalog, ret = SaveCatalog(helper, ca, true)
	if !ret {
		t.Error("SaveCatalog failed")
	}

	ca, found = QueryCatalogByID(helper, catalog.ID)
	if !found {
		t.Error("QueryCatalogByID failed")
	}

	if len(ca.Catalog) != 3 {
		t.Error("QueryCatalogByID failed")
	}

	ret = DeleteCatalog(helper, ca.ID)
	if !ret {
		t.Error("DeleteCatalog failed")
	}

	catalogs := QueryAllCatalog(helper)
	if len(catalogs) != 3 {
		t.Error("QueryAllCatalog failed")
	}

	catalogDetails := QueryAllCatalog(helper)
	if len(catalogDetails) != 3 {
		t.Error("QueryAllCatalogDetail")
	}
}
