$(document).ready(function () {
    $(".item").each(function () {
        $(this).attr('data-id-team', $(this).parent().attr('data-id-team'));
    })
    $('[data-toggle="tooltip"]').tooltip();
    var role = $("#role").attr('data-role');
    var id_user = $("#id_user").attr('data-id');
    handleDrag(role);

    // Listen msg from server
    socket.addEventListener("message", function(e) {
        var msg = JSON.parse(e.data);
        if (msg.MoveType === undefined) {
            onlineUsers(msg);
        }
        // movetype = 0 when move member on team page
        if (msg.MoveType == 0) {
            var item = $("[data-id-member=" + msg.member_id + "][data-id-team=" + msg.team_id_previous + "]");
            if (item.find(".fa-star").length > 0) {
                item.children(".fa-star").remove();
            }
            $(".box[data-id-team=" + msg.team_id + "]").append(item);
            item.attr('data-id-team', msg.team_id);
            handleDrag(role);
        }
    });

    $(".box").on("ss-added", function (e, selected) {
        var member_id = $(selected).attr('data-id-member');
        var team_id_previous = $(selected).attr('data-id-team');
        if ($(this).find("[data-id-member=" + member_id + "]").length == 1) {
            if ($(selected).find(".fa-star").length > 0) {
                $(selected).children(".fa-star").remove();
            }
            var team_id = $(this).attr('data-id-team');
            $(selected).attr('data-id-team', team_id);
            $.ajax({
                url: "/admin/moveMember",
                type: 'POST',
                data: {
                    teamID: team_id,
                    memberID: member_id,
                    teamIDPrevious: team_id_previous,
                },
            })
                .done(function (res) {
                    if (res == team_id) {
                        // Send to server
                        socket.send(
                            JSON.stringify({
                                team_id: team_id,
                                team_id_previous: team_id_previous,
                                member_id: member_id,
                                user_role: role,
                                user_id: id_user,
                                MoveType: 0,
                            })
                        );
                    } else if (res == "false") {
                        // res = false when update database failed
                        $(selected).attr("data-id-team", team_id_previous);
                        $(".box[data-id-team=" + team_id_previous + "]").append($(selected));
                        handleDrag(role);
                    } else {
                        // when drop one member to different teams at the same time
                        $(selected).attr("data-id-team", res);
                        $(".box[data-id-team=" + res + "]").append($(selected));
                        handleDrag(role);
                    }
                })
                .fail(function () {
                    console.log("error");
                });
        } else {
            $(".box[data-id-team=" + team_id_previous + "]").append($(selected));
            handleDrag(role);
        }
    });

});

// Drag & Drop member
function handleDrag(role) {
    if (role == 2) {
        $(".box").shapeshift()
    } else {
        $(".box").shapeshift({
            enableDrag: false
        });
    }
}

// Show user online
function onlineUsers(arr) {
    arr = getUnique(arr, 'ID');
    var lengthUser = 3;
    var div = $("#navList");
    div.empty();
    $('.tooltip').remove();
    var min = arr.length > lengthUser ? lengthUser : arr.length;
    for (i = 0; i < min; i++) {
        var div1 = $("<div>").addClass("rounded-circle").appendTo($(div));
        var img = $("<img>").attr({
            "data-toggle": "tooltip",
            "title": arr[i].Name,
            "class": "rounded-circle onlineUser",
            "src": arr[i].PictureURL.String
        }).appendTo($(div1));
    }
    if (arr.length > lengthUser) {
        var div2 = $("<div>").addClass("rounded-circle").appendTo($(div));
        var div3 = $("<div>").attr({
            "class": "numberCircle",
            "data-toggle": "dropdown"
        }).text("+" + (arr.length - 3).toString()).appendTo($(div2));
        var div4 = $("<div>").attr({
            "class": "dropdown-menu",
            "style": "width:210px"
        }).appendTo($(div2));
        for (i = 0; i < arr.length; i++) {
            var a = $("<a>").attr({
                "class": "dropdown-item",
                "href": "#"
            }).appendTo($(div4));
            var img1 = $("<img>").attr({
                "class": "rounded-circle onlineUser",
                "src": arr[i].PictureURL.String
            }).appendTo($(a));
            var h5 = $("<h5>").text(arr[i].Name).appendTo($(a));
        }
    }
    $('[data-toggle="tooltip"]').tooltip();
}

function getUnique(arr, comp) {
    const unique = arr
        .map(e => e[comp])
        .map((e, i, final) => final.indexOf(e) === i && i)
        .filter(e => arr[e]).map(e => arr[e]);
    return unique;
}
