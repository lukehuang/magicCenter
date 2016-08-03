package dal

import (
	"magiccenter/kernel/modules/account/model"
	"magiccenter/util/dbhelper"
	"testing"
)

func TestGroup(t *testing.T) {
	helper, err := dbhelper.NewHelper()
	if err != nil {
		panic("construct helper failed")
	}
	defer helper.Release()

	g1 := model.Group{}
	g1.Name = "test1"
	g1.Type = model.AdminGroup

	SaveGroup(helper, g1)

	g2 := model.Group{}
	g2.Name = "test2"
	g2.Type = model.CommonGroup
	SaveGroup(helper, g2)

	groups := QueryAllGroup(helper)
	if len(groups) < 2 {
		t.Error("SaveGroup failed")
		return
	}

	g11, found := QueryGroupByName(helper, "test1")
	if !found {
		t.Errorf("QueryGroupByName failed, name=%s", "test1")
		return
	}
	if !g11.AdminGroup() {
		t.Error("QueryGroupByName return invalid groupType")
		return
	}

	g111, found := QueryGroupByID(helper, g11.ID)
	if !found {
		t.Errorf("QueryGroupByID failed, id=%d", g11.ID)
		return
	}

	if g111.Name != "test1" {
		t.Error("QueryGroupByID return invalid groupName")
		return
	}

	ok := DeleteGroup(helper, g11.ID)
	if !ok {
		t.Errorf("DeleteGroup failed, id=%d", g11.ID)
		return
	}

	g22, found := QueryGroupByName(helper, "test2")
	if !found {
		t.Errorf("QueryGroupByName failed, name=%s", "test2")
		return
	}

	DeleteGroup(helper, g22.ID)
}
