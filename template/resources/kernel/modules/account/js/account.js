var account = {
    userInfos: {},
    groupInfos: {}
};

$(document).ready(function() {

    // 绑定表单提交事件处理器
    $("#account-Content .account-Form").submit(function() {
        var options = {
            beforeSubmit: showRequest,
            success: showResponse,
            dataType: "json"
        };

        function showRequest() {}

        function showResponse(result) {

            if (result.ErrCode > 0) {
                $("#account-Edit .alert-Info .content").html(result.Reason);
                $("#account-Edit .alert-Info").modal();
            } else {
                account.refreshAccount();
                account.fillAccountListView();
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

    $("#account-Content .account-Form .account-account").change(function() {
        $("#account-Content .account-Form .account-account").parent().removeClass("has-error");
        var account = $("#account-Content .account-Form .account-account").val();
        $.get("/account/checkAccount/?account=" + account, {}, function(result) {
            if (result.ErrCode == 0) {
                return;
            }
            $("#account-Content .account-Form .account-account").parent().addClass("has-error");
            console.log(result);
        }, "json");
    });

    $("#account-Content .account-Form button.reset").click(function() {
        $("#account-Edit .account-Form .account-account").val("");
        $("#account-Edit .account-Form .account-email").val("");
        $("#account-Edit .account-Form .account-group").prop("checked", false);
        $("#account-Edit .account-Form .account-id").val("-1");
    });
});

account.initialize = function() {
    account.fillAccountListView();

    account.fillGroupListView();
};

account.getAccountListView = function() {
    return $("#account-List table");
};

account.getGroupListView = function() {
    return $("#account-Edit .account-Form .account-group");
};

account.constructAccountItem = function(userInfo) {
    var tr = document.createElement("tr");

    var nameTd = document.createElement("td");
    nameTd.innerHTML = userInfo.Account;
    tr.appendChild(nameTd);

    var emailTd = document.createElement("td");
    emailTd.innerHTML = userInfo.Email;
    tr.appendChild(emailTd);

    var groupTd = document.createElement("td");
    var groups = "";
    for (var ii = 0; ii < userInfo.Groups.length;) {
        var gid = userInfo.Groups[ii++];
        for (var jj = 0; jj < account.groupInfos.length;) {
            var group = account.groupInfos[jj++];
            if (group.ID == gid) {
                groups += group.Name;

                if (ii < userInfo.Groups.length) {
                    groups += ",";
                }
                break;
            }
        }
    }
    groupTd.innerHTML = groups;
    tr.appendChild(groupTd);

    var statusTd = document.createElement("td");
    switch (userInfo.Status) {
        case 0:
            statusTd.innerHTML = "新建";
            break;
        case 1:
            statusTd.innerHTML = "激活";
            break;
        case 2:
            statusTd.innerHTML = "未激活";
            break;
        case 3:
            statusTd.innerHTML = "锁定";
            break;
        default:
            statusTd.innerHTML = "异常";
            break;
    }
    tr.appendChild(statusTd);

    var opTd = document.createElement("td");
    var editLink = document.createElement("a");
    editLink.setAttribute("class", "edit");
    editLink.setAttribute("href", "#editAccount");
    editLink.setAttribute("onclick", "account.editAccount('/account/queryUser/?id=" + userInfo.ID + "'); return false");
    var editImage = document.createElement("img");
    editImage.setAttribute("src", "/common/images/pencil.png");
    editImage.setAttribute("alt", "Edit");
    editLink.appendChild(editImage);
    opTd.appendChild(editLink);

    var deleteLink = document.createElement("a");
    deleteLink.setAttribute("class", "delete");
    deleteLink.setAttribute("href", "#deleteAccount");
    deleteLink.setAttribute("onclick", "account.deleteAccount('/account/deleteUser/?id=" + userInfo.ID + "'); return false;");
    var deleteImage = document.createElement("img");
    deleteImage.setAttribute("src", "/common/images/cross.png");
    deleteImage.setAttribute("alt", "Delete");
    deleteLink.appendChild(deleteImage);
    opTd.appendChild(deleteLink);
    tr.appendChild(opTd);

    return tr;
};

account.fillAccountListView = function() {
    var userListView = account.getAccountListView();

    $(userListView).find("tbody tr").remove();
    for (var ii = 0; ii < account.userInfos.length; ++ii) {
        var userInfo = account.userInfos[ii];
        var trContent = account.constructAccountItem(userInfo);

        $(userListView).find("tbody").append(trContent);
    }

    $("#account-Edit .account-Form .account-account").val("");
    $("#account-Edit .account-Form .account-email").val("");
    $("#account-Edit .account-Form .account-group").prop("checked", false);
    $("#account-Edit .account-Form .account-id").val("-1");
};

account.constructGroupItem = function(group) {
    var label = document.createElement("label");

    var chk = document.createElement("input");
    chk.setAttribute("type", "checkbox");
    chk.setAttribute("name", "account-group");
    chk.setAttribute("class", "account-group");
    chk.setAttribute("value", group.ID);
    label.appendChild(chk);

    var span = document.createElement("span");
    span.innerHTML = group.Name;
    label.appendChild(span);
    label.setAttribute("class", "text-center");

    return label;
};

account.fillGroupListView = function() {
    var groupListView = account.getGroupListView();

    $(groupListView).find("label").remove();
    for (var ii = 0; ii < account.groupInfos.length; ++ii) {
        var cur = account.groupInfos[ii];
        var label = account.constructGroupItem(cur);
        $(groupListView).append(label);
    }
}

account.refreshAccount = function() {
    $.get("/account/queryAllUser/", {}, function(result) {
        account.userInfos = result.Users;

        account.fillAccountListView();
    }, "json");
};

account.editAccount = function(editUrl) {
    $.get(editUrl, {}, function(result) {
        if (result.ErrCode > 0) {
            $("#account-List .alert-Info .content").html(result.Reason);
            $("#account-List .alert-Info").modal();
            return
        }

        $("#account-Edit .account-Form .account-account").val(result.User.Account);
        $("#account-Edit .account-Form .account-email").val(result.User.Email);
        $("#account-Edit .account-Form .account-id").val(result.User.ID);

        var groupListView = account.getGroupListView();
        for (var ii = 0; ii < result.User.Groups.length; ++ii) {
            var gid = result.User.Groups[ii];
            $(groupListView).find("input ").filter("[value=" + gid + "]").prop("checked", true);
        }

        $("#account-Content .content-header .nav .account-Edit").find("a").trigger("click");
    }, "json");
}

account.deleteAccount = function(deleteUrl) {
    $.get(deleteUrl, {}, function(result) {
        if (result.ErrCode > 0) {
            $("#account-List .alert-Info .content").html(result.Reason);
            $("#account-List .alert-Info").modal();
            return
        }

        account.refreshAccount();
    }, "json");
};