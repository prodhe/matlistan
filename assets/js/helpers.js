/* helper functions */

function sortAlpha(a, b) {
    return (a.toLowerCase() > b.toLowerCase()) ? 1 : -1;
}

function safe2html(str) {
    str = str.replace(/</g, "&lt;");
    str = str.replace(/>/g, "&gt;");
    return str;
}

function addObjToList(obj) {
    $(obj)
        .clone()
        .data("placement", "left")
        .tooltip()
        .appendTo($('#userlist'));
    $('#userlist').trigger("sortupdate");
}

function clearUserlist() {
    $('#userlist').children().remove();
    $('#userlist').trigger("sortupdate");
}