var group = {
    groupInfos: {}
};

$(document).ready(function() {

    // 绑定表单提交事件处理器
    $("#group-Content .group-Form").submit(function() {
        var options = {
            beforeSubmit: showRequest,
            success: showResponse,
            dataType: "json"
        };

        function showRequest() {}

        function showResponse(result) {
            console.log(result);

            if (result.ErrCode > 0) {
                $("#group-Edit .alert-Info .content").html(result.Reason);
                $("#group-Edit .alert-Info").modal();
            } else {
                group.groupInfos = result.Groups;

                group.fillGroupListView();
            }
        }

        function validate() {
            var result = true;
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

    $("#group-Edit .group-Form button.reset").click(
        function() {
            $("#group-Edit .group-Form .group-id").val("");
            $("#group-Edit .group-Form .group-name").val("");
        });
});

group.initialize = function() {
    group.fillGroupListView();
};

group.getGroupListView = function() {
    return $("#group-List table");
};

group.constructGroupItem = function(groupInfo) {
    var tr = document.createElement("tr");

    var nameTd = document.createElement("td");
    nameTd.innerHTML = groupInfo.Name;
    tr.appendChild(nameTd);

    var userCountTd = document.createElement("td");
    userCountTd.innerHTML = groupInfo.UserCount;
    tr.appendChild(userCountTd);

    var createrlTd = document.createElement("td");
    createrlTd.innerHTML = groupInfo.Creater.Name;
    tr.appendChild(createrlTd);

    var opTd = document.createElement("td");
    var editLink = document.createElement("a");
    editLink.setAttribute("class", "edit");
    editLink.setAttribute("href", "#editGroup");
    editLink.setAttribute("onclick", "group.editGroup('/admin/account/queryGroup/?id=" + groupInfo.Id + "'); return false");
    var editImage = document.createElement("img");
    editImage.setAttribute("src", "/resources/admin/images/pencil.png");
    editImage.setAttribute("alt", "Edit");
    editLink.appendChild(editImage);
    opTd.appendChild(editLink);

    var deleteLink = document.createElement("a");
    deleteLink.setAttribute("class", "delete");
    deleteLink.setAttribute("href", "#deleteGroup");
    deleteLink.setAttribute("onclick", "group.deleteGroup('/admin/account/deleteGroup/?id=" + groupInfo.Id + "'); return false;");
    var deleteImage = document.createElement("img");
    deleteImage.setAttribute("src", "/resources/admin/images/cross.png");
    deleteImage.setAttribute("alt", "Delete");
    deleteLink.appendChild(deleteImage);
    opTd.appendChild(deleteLink);
    tr.appendChild(opTd);

    return tr;
};

group.fillGroupListView = function() {
    var groupListView = group.getGroupListView();

    $(groupListView).find("tbody tr").remove();
    for (var ii = 0; ii < group.groupInfos.length; ++ii) {
        var groupInfo = group.groupInfos[ii];
        var trContent = group.constructGroupItem(groupInfo);

        $(groupListView).find("tbody").append(trContent);
    }

    $("#group-Edit .group-Form .group-id").val("");
    $("#group-Edit .group-Form .group-name").val("");
};

group.editGroup = function(editUrl) {
    $.get(editUrl, {}, function(result) {
        if (result.ErrCode > 0) {
            $("#group-List .alert-Info .content").html(result.Reason);
            $("#group-List .alert-Info").modal();
            return
        }

        $("#group-Edit .group-Form .group-name").val(result.Group.Name);
        $("#group-Edit .group-Form .group-id").val(result.Group.Id);
        $("#group-Content .content-header .nav .group-Edit").find("a").trigger("click");
    }, "json");
}

group.deleteGroup = function(deleteUrl) {
    $.get(deleteUrl, {}, function(result) {
        if (result.ErrCode > 0) {
            $("#group-List .alert-Info .content").html(result.Reason);
            $("#group-List .alert-Info").modal();
            return
        }

        group.groupInfos = result.Groups;
        group.fillGroupListView();
    }, "json");
};