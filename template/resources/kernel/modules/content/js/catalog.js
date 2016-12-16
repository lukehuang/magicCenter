var catalog = {
    catalogInfo: {},
    userInfo: {}
};

$(document).ready(function() {
    // 绑定表单提交事件处理器
    $("#catalog-Content .catalog-Form").submit(function() {
        var options = {
            beforeSubmit: showRequest, // pre-submit callback
            success: showResponse, // post-submit callback
            dataType: 'json' // 'xml', 'script', or 'json' (expected server response type) 
        };

        // pre-submit callback
        function showRequest() {
            //return false;
        }
        // post-submit callback
        function showResponse(result) {
            if (result.ErrCode > 0) {
                $("#catalog-Edit .alert-Info .content").html(result.Reason);
                $("#catalog-Edit .alert-Info").modal();
            } else {
                catalog.refreshCatalog();
            }
        }

        function validate() {
            var result = true

            $("#catalog-Content .catalog-Form .catalog-name").parent().find("span").remove();
            var title = $("#catalog-Content .catalog-Form .catalog-name").val();
            if (title.length == 0) {
                $("#catalog-Content .catalog-Form .catalog-name").parent().append("<span class=\"input-notification error png_bg\">请输入分类名</span>");
                result = false;
            }

            return result;
        }

        if (!validate()) {
            return false;
        }

        //提交表单
        $(this).ajaxSubmit(options);

        // !!! Important !!!
        // 为了防止普通浏览器进行表单提交和产生页面导航（防止页面刷新？）返回false
        return false;
    });

    $("#catalog-Content .catalog-Form button.reset").click(function() {
        $("#catalog-Edit .catalog-Form .catalog-id").val(-1);
        $("#catalog-Edit .catalog-Form .catalog-name").val("");
        $("#catalog-Edit .catalog-Form .catalog-parent").children().remove();
        for (var ii = 0; ii < catalog.catalogInfo.length; ++ii) {
            var ca = catalog.catalogInfo[ii];
            $("#catalog-Edit .catalog-Form .catalog-parent").append("<label><input type='checkbox' name='catalog-parent' value=" + ca.ID + "> </input>" + ca.Name + "</label> ");
        }
        if (ii == 0) {
            $("#catalog-Edit .catalog-Form .catalog-parent").append("<label><input type='checkbox' name='catalog-parent' readonly='readonly' value='-1' onclick='return false'> </input>-</label> ");
        }
    });
});

catalog.initialize = function() {

    catalog.fillCatalogView();
};

catalog.refreshCatalog = function() {
    $.get("/content/queryAllCatalog/", {}, function(result) {
        catalog.catalogInfo = result.Catalogs;

        catalog.fillCatalogView();
    }, "json");
};


catalog.fillCatalogView = function() {
    $("#catalog-List table tbody tr").remove();
    for (var ii = 0; ii < catalog.catalogInfo.length; ++ii) {
        var info = catalog.catalogInfo[ii];
        var trContent = catalog.constructCatalogItem(info);
        $("#catalog-List table tbody").append(trContent);
    }

    $("#catalog-Edit .catalog-Form .catalog-id").val(-1);
    $("#catalog-Edit .catalog-Form .catalog-name").val("");


    $("#catalog-Edit .catalog-Form .catalog-parent").children().remove();
    for (var ii = 0; ii < catalog.catalogInfo.length; ++ii) {
        var ca = catalog.catalogInfo[ii];
        $("#catalog-Edit .catalog-Form .catalog-parent").append("<label><input type='checkbox' name='catalog-parent' value=" + ca.ID + "> </input>" + ca.Name + "</label> ");
    }
    if (ii == 0) {
        $("#catalog-Edit .catalog-Form .catalog-parent").append("<label><input type='checkbox' name='catalog-parent' readonly='readonly' value='-1' onclick='return false'> </input>-</label> ");
    }
};


catalog.constructCatalogItem = function(ca) {
    var tr = document.createElement("tr");
    tr.setAttribute("class", "catalog");

    var titleTd = document.createElement("td");
    var titleLink = document.createElement("a");
    titleLink.setAttribute("class", "edit");
    titleLink.setAttribute("href", "#editCatalog");
    titleLink.setAttribute("onclick", "catalog.editCatalog('/content/queryCatalog/?id=" + ca.ID + "'); return false;");
    titleLink.innerHTML = ca.Name;
    titleTd.appendChild(titleLink);
    tr.appendChild(titleTd);

    var parentTd = document.createElement("td");
    var parentCatalogs = "";
    if (ca.Parent) {
        for (var ii = 0; ii < ca.Parent.length;) {
            var cid = ca.Parent[ii++];
            for (var jj = 0; jj < catalog.catalogInfo.length;) {
                var singleCatalog = catalog.catalogInfo[jj++];
                if (singleCatalog.ID == cid) {
                    parentCatalogs += singleCatalog.Name;
                    if (ii < ca.Parent.length) {
                        parentCatalogs += ",";
                    }
                    break;
                }
            }
        }
    }
    parentCatalogs = parentCatalogs.length == 0 ? '-' : parentCatalogs;
    parentTd.innerHTML = parentCatalogs;
    tr.appendChild(parentTd);

    var createrTd = document.createElement("td");
    var createrValue = "-"
    for (var ii = 0; ii < catalog.userInfo.length; ii++) {
        var author = catalog.userInfo[ii];
        if (author.ID == ca.Creater) {
            createrValue = author.Name;
            break;
        }
    }
    createrTd.innerHTML = createrValue;
    tr.appendChild(createrTd);

    var editTd = document.createElement("td");
    var editLink = document.createElement("a");
    editLink.setAttribute("class", "edit");
    editLink.setAttribute("href", "#editCatalog");
    editLink.setAttribute("onclick", "catalog.editCatalog('/content/queryCatalog/?id=" + ca.ID + "'); return false;");
    var editImage = document.createElement("img");
    editImage.setAttribute("src", "/common/images/pencil.png");
    editImage.setAttribute("alt", "Edit");
    editLink.appendChild(editImage);
    editTd.appendChild(editLink);

    var deleteLink = document.createElement("a");
    deleteLink.setAttribute("class", "delete");
    deleteLink.setAttribute("href", "#deleteCatalog");
    deleteLink.setAttribute("onclick", "catalog.deleteCatalog('/content/deleteCatalog/?id=" + ca.ID + "'); return false;");
    var deleteImage = document.createElement("img");
    deleteImage.setAttribute("src", "/common/images/cross.png");
    deleteImage.setAttribute("alt", "Delete");
    deleteLink.appendChild(deleteImage);
    editTd.appendChild(deleteLink);

    tr.appendChild(editTd);

    return tr;
};

catalog.editCatalog = function(editUrl) {
    $.get(editUrl, {}, function(result) {
        if (result.ErrCode > 0) {
            $("#catalog-List .alert-Info .content").html(result.Reason);
            $("#catalog-List .alert-Info").modal();
            return
        }

        $("#catalog-Edit .catalog-Form .catalog-id").val(result.Catalog.ID);
        $("#catalog-Edit .catalog-Form .catalog-name").val(result.Catalog.Name);
        $("#catalog-Edit .catalog-Form .catalog-parent").children().remove();
        if (catalog.catalogInfo) {
            var hasParentCatalog = false;
            for (var ii = 0; ii < catalog.catalogInfo.length;) {
                var ca = catalog.catalogInfo[ii++];
                if (ca.ID < result.Catalog.ID) {
                    $("#catalog-Edit .catalog-Form .catalog-parent").append("<label><input type='checkbox' name='catalog-parent' value=" + ca.ID + "> </input>" + ca.Name + "</label> ");
                    hasParentCatalog = true;
                }
            }

            if (!hasParentCatalog) {
                $("#catalog-Edit .catalog-Form .catalog-parent").append("<label><input type='checkbox' name='catalog-parent' readonly='readonly' value='-1' onclick='return false'> </input> - </label> ");
            }
        }

        if (result.Catalog.Parent) {
            for (var ii = 0; ii < result.Catalog.Parent.length; ++ii) {
                var ca = result.Catalog.Parent[ii];
                $("#catalog-Edit .catalog-Form .catalog-parent input").filter("[value=" + ca + "]").prop("checked", true);
            }
        }

        $("#catalog-Content .content-header .nav .catalog-Edit").find("a").trigger("click");
    }, "json");
};

catalog.deleteCatalog = function(deleteUrl) {
    $.get(deleteUrl, {}, function(result) {
        if (result.ErrCode > 0) {
            $("#catalog-List .alert-Info .content").html(result.Reason);
            $("#catalog-List .alert-Info").modal();
            return;
        }

        catalog.refreshCatalog();
    }, "json");
};