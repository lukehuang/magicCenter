var authgroup = {};

authgroup.constructAclListlView = function(aclList, authGroupList, moduleList) {
    var aclListView = new Array();
    var offset = 0;
    for (var i = 0; i < aclList.length; ++i) {
        var curAcl = aclList[i];

        var curModule = null;
        for (var idx = 0; idx < moduleList.length; ++idx) {
            var mod = moduleList[idx];
            if (mod.ID == curAcl.Module) {
                curModule = mod;
                break;
            }
        }

        var curAuthGroup = "";
        for (var idx = 0; idx < authGroupList.length; ++idx) {
            var cur = authGroupList[idx];
            for (var ii = 0; ii < curAcl.AuthGroup.length; ++ii) {
                var id = curAcl.AuthGroup[ii];
                if ((cur.ID == id) && (cur.Module == curModule.ID)) {
                    if (curAuthGroup.length > 0) {
                        curAuthGroup += ",";
                    }

                    curAuthGroup += cur.Name;
                }
            }
        }

        var view = {
            ID: curAcl.ID,
            URL: curAcl.URL,
            Method: curAcl.Method,
            Module: curModule.Name,
            AuthGroup: curAuthGroup
        };

        aclListView[offset++] = view;
    }

    return aclListView;
};

authgroup.constructAclEditView = function(aclList, authGroupList, moduleList) {
    var aclListView = new Array();
    var ii = 0;
    for (var aclIdx = 0; aclIdx < aclList.length; ++aclIdx) {
        var curAcl = aclList[aclIdx];

        for (var idx = 0; idx < moduleList.length; ++idx) {
            var curModule = moduleList[idx];

            if (curAcl.Module == curModule.ID) {
                var view = {
                    ID: curAcl.ID,
                    URL: curAcl.URL,
                    Method: curAcl.Method,
                    Status: curAcl.Status,
                    Module: curModule.Name,
                    ModuleID: curModule.ID
                }

                aclListView[ii++] = view;
            }
        }
    }

    return aclListView;
};

authgroup.updateListAclVM = function(aclList) {
    authgroup.listVM.acls = aclList;
};

authgroup.updateEditAclVM = function(aclList) {
    authgroup.editVM.acls = aclList;
};

authgroup.updateEditAuthGroupVM = function(authGroupList) {
    authgroup.editVM.authGroups = authGroupList;
};

// 加载全部的Module
authgroup.getAllModulesAction = function(callBack) {
    $.ajax({
        type: "GET",
        url: "/dashboard/module/",
        data: {},
        dataType: "json",
        success: function(data) {
            if (callBack != null) {
                callBack(data.ErrCode, data.Module);
            }
        }
    });
};

// 获取全部的AuthGroup
authgroup.getAllAuthGroupsAction = function(callBack) {
    $.ajax({
        type: "GET",
        url: "/cas/authgroup/?module=all",
        data: {},
        dataType: "json",
        success: function(data) {
            if (callBack != null) {
                callBack(data.ErrCode, data.AuthGroup);
            }
        }
    });
};

// 加载全部已经定义的ACL
authgroup.getAllAclsAction = function(callBack) {
    $.ajax({
        type: "GET",
        url: "/cas/acl/?module=all&status=1",
        data: {},
        dataType: "json",
        success: function(data) {
            if (callBack != null) {
                callBack(data.ErrCode, data.ACLs);
            }
        }
    });
};

// 更新ACL对应的AuthGroup
authgroup.updateAclAuthGroupAction = function(aclID, authGroups, callBack) {
    var strAuthGroups = "";
    for (var ii = 0; ii < authGroups.length; ++ii) {
        strAuthGroups += authGroups[ii];
        strAuthGroups += ",";
    }

    $.ajax({
        type: "POST",
        url: "/cas/acl/authgroup/",
        data: { "acl-id": aclID, "acl-authgroup": strAuthGroups },
        dataType: "json",
        success: function(data) {
            if (callBack != null) {
                callBack(data.ErrCode);
            }
        }
    });
};

authgroup.loadData = function(callBack) {
    var getAllAclsCallBack = function(errCode, aclList) {
        if (errCode != 0) {
            return;
        }

        authgroup.acls = aclList;
        if (callBack != null) {
            callBack()
        }
    };

    var getAllAuthGroupsCallBack = function(errCode, authGroupList) {
        if (errCode != 0) {
            return;
        }

        authgroup.authGroups = authGroupList;
        authgroup.getAllAclsAction(getAllAclsCallBack)
    };

    var getAllModulesCallBack = function(errCode, moduleList) {
        if (errCode != 0) {
            return;
        }

        authgroup.modules = moduleList;
        authgroup.getAllAuthGroupsAction(getAllAuthGroupsCallBack);
    };

    // 加载完成
    authgroup.getAllModulesAction(getAllModulesCallBack);
};

authgroup.refreshAclListView = function(aclList, authGroupList, moduleList) {
    var aclsView = authgroup.constructAclListlView(aclList, authGroupList, moduleList);
    authgroup.updateListAclVM(aclsView);
};

authgroup.refreshAclEidtView = function(filterVal, aclList, authGroupList, moduleList) {
    var aclsView = authgroup.constructAclEditView(aclList, authGroupList, moduleList);
    var filterAclsView = new Array();
    var filterModuleList = new Array();
    var offset = 0;
    for (var idx = 0; idx < aclsView.length; ++idx) {
        var cur = aclsView[idx];
        if (cur.URL.search(filterVal) != -1) {
            if (filterModuleList.length == 0) {
                filterModuleList[filterModuleList.length] = cur.ModuleID;
            } else if (filterModuleList[0] != cur.ModuleID) {
                filterModuleList[filterModuleList.length] = cur.ModuleID;
            }

            filterAclsView[offset++] = cur;
            continue;
        }

        if (cur.Module.search(filterVal) != -1) {
            if (filterModuleList.length == 0) {
                filterModuleList[filterModuleList.length] = cur.ModuleID;
            } else if (filterModuleList[0] != cur.ModuleID) {
                filterModuleList[filterModuleList.length] = cur.ModuleID;
            }

            filterAclsView[offset++] = cur;
            continue;
        }
    }

    if (filterModuleList.length == 1) {
        var filterAuthgroupsView = new Array();
        for (var ii = 0; ii < authGroupList.length; ++ii) {
            var cur = authGroupList[ii];
            if (cur.Module == filterModuleList[0]) {
                filterAuthgroupsView[filterAuthgroupsView.length] = cur;
            }
        }

        authgroup.updateEditAclVM(filterAclsView);
        authgroup.updateEditAuthGroupVM(filterAuthgroupsView);
    }
};

$(document).ready(function() {
    authgroup.listVM = avalon.define({
        $id: "authgroup-List",
        acls: []
    });

    authgroup.editVM = avalon.define({
        $id: "authgroup-Edit",
        acls: [],
        authGroups: []
    });

    // 过滤出符合要求的ACL
    $("#authgroup-Edit .filter-button").click(
        function() {
            var filterVal = $("#authgroup-Edit .authgroup-filter").val();
            authgroup.refreshAclEidtView(filterVal, authgroup.acls, authgroup.authGroups, authgroup.modules);
        }
    );

    $("#authgroup-Edit .adjust-button").click(
        function() {}
    );

    authgroup.loadData(function() {
        authgroup.refreshAclListView(authgroup.acls, authgroup.authGroups, authgroup.modules);
    })
});