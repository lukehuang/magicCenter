package handler

import (
	"strings"

	"muidea.com/magicCenter/application/common/dbhelper"
	"muidea.com/magicCenter/application/common/model"
	"muidea.com/magicCenter/application/module/kernel/modules/cas/dal"
)

type authGroupManager struct {
	module2AuthGroup map[string][]model.AuthGroup
}

func createAuthGroupManager() authGroupManager {
	authGroupManager := authGroupManager{module2AuthGroup: make(map[string][]model.AuthGroup)}
	authGroupManager.loadAllAuthGroup()

	return authGroupManager
}

func (i *authGroupManager) loadAllAuthGroup() bool {
	dbhelper, err := dbhelper.NewHelper()
	if err != nil {
		return false
	}

	authGroups := dal.GetAllAuthGroup(dbhelper)
	for _, authGroup := range authGroups {
		authGroups, found := i.module2AuthGroup[authGroup.Module]
		if !found {
			authGroups = []model.AuthGroup{}
		}

		authGroups = append(authGroups, authGroup)
		i.module2AuthGroup[authGroup.Module] = authGroups
	}

	return true
}

func (i *authGroupManager) queryAuthGroup(module string) ([]model.AuthGroup, bool) {
	if strings.ToLower(module) == "all" {
		authGroups := []model.AuthGroup{}
		for _, groups := range i.module2AuthGroup {
			authGroups = append(authGroups, groups...)
		}

		return authGroups, true
	}

	authGroups, found := i.module2AuthGroup[module]
	return authGroups, found
}

func (i *authGroupManager) insertAuthGroup(authGroups []model.AuthGroup) bool {
	dbhelper, err := dbhelper.NewHelper()
	if err != nil {
		return false
	}

	for _, authGroup := range authGroups {
		authGroup, ok := dal.InsertAuthGroup(dbhelper, authGroup)
		if ok {
			authGroups, found := i.module2AuthGroup[authGroup.Module]
			if !found {
				authGroups = []model.AuthGroup{}
			}

			authGroups = append(authGroups, authGroup)
			i.module2AuthGroup[authGroup.Module] = authGroups
		}
	}

	return true
}

func (i *authGroupManager) deleteAuthGroup(authGroups []model.AuthGroup) bool {
	dbhelper, err := dbhelper.NewHelper()
	if err != nil {
		return false
	}
	for _, v := range authGroups {
		dal.DeleteAuthGroup(dbhelper, v.ID)

		curAuthGroups, found := i.module2AuthGroup[v.Module]
		newAuthGroups := []model.AuthGroup{}
		if found {
			for _, c := range curAuthGroups {
				if c.ID != v.ID {
					newAuthGroups = append(newAuthGroups, c)
				}
			}
			if len(newAuthGroups) > 0 {
				i.module2AuthGroup[v.Module] = newAuthGroups
			}
		}
	}

	return true
}

func (i *authGroupManager) adjustUserAuthGroup(userID int, authGroup []int) bool {
	dbhelper, err := dbhelper.NewHelper()
	if err != nil {
		return false
	}
	return dal.UpateUserAuthorityGroup(dbhelper, userID, authGroup)
}

func (i *authGroupManager) getUserAuthGroup(userID int) ([]int, bool) {
	dbhelper, err := dbhelper.NewHelper()
	if err != nil {
		return []int{}, false
	}

	return dal.GetUserAuthorityGroup(dbhelper, userID), true
}
