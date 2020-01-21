$(document).ready(function () {
    $('[data-toggle="tooltip"]').tooltip();

    // export seat_master to pdf
    $("#export_pdf").click(function () {
        var seat_name = $(".nav-link.active").text();
        var active_tab = $(".seattab.active").attr('id');
        $('.' + active_tab).css('background', '#ffffff');
        var seat = document.getElementsByClassName(active_tab);
        html2canvas(seat, {
            onrendered: function (canvas) {
                var imgData = canvas.toDataURL(
                    'image/png');
                const ratio = 0.2645833333;
                var width = canvas.width * ratio;
                var height = canvas.height * ratio;
                var scale = width / height;
                var orientation = (scale <= 0.9) ? "p" : "l";
                var width_actual = (scale <= 0.9) ? 190 : 280;
                var height_actual = width_actual / scale;
                if (orientation == "p" && height_actual > 280) {
                    height_actual = 260;
                    width_actual = height_actual * scale;
                }
                if (orientation == "l" && height_actual > 190) {
                    height_actual = 170;
                    width_actual = height_actual * scale;
                }
                var margin_x = (width >= height) ? (297 - width_actual) / 2 : (210 - width_actual) / 2;
                var doc = new jsPDF(orientation, "mm", "a4");
                doc.addImage(imgData, 'PNG', margin_x, 10, width_actual, height_actual);
                var text_x = (orientation == "l") ? 148.5 : 105;
                doc.text(seat_name, text_x, 8, 'center');
                doc.save(seat_name + '.pdf');
            }
        });
    });
    // Listen msg from server
    socket.addEventListener("message", function (e) {
        var msg = JSON.parse(e.data);
        // movetype = 1 when move member on seat page
        if (msg.MoveType == 1) {
            if (msg.DepartmentID > 0) {
                var position = $(".seattab-" + msg.SeatMasterID).find($(".box-department[data-department-id = " + msg.DepartmentID + "]"));
                var item = $(".seattab-" + msg.SeatMasterID).find($(".child[data-member-id = " + msg.MemberID + "]"));
                item.appendTo($(position));
            } else {
                appendMemberToSeat(msg);
            }

        }
    });

    //display seat_master on user seat
    if ($("#check_page[data-page=seat]").length > 0) {
        var $seattab = $(".seattab-content").children();
        $.each($seattab, function (key) {
            var id = $(this).attr("data-seatmaster-id");
            $.ajax({
                type: "POST",
                url: "/seats/getSeats",
                data: {
                    id: id
                },
                success: function (res) {
                    res = JSON.parse(res);
                    var seats = JSON.parse(res.SeatMaster);
                    var seatmaster_id = res.ID;
                    $(".seattab" + key).css({
                        width: seats[3][0] + 'px',
                        height: seats[3][1] + 'px',
                        margin: 'auto'
                    });
                    for (i = 0; i < seats[0].length; i++) {
                        var seat = createStaticTable(seats[0][i][0], seats[0][i][2], seats[0][i][1], seats[0][i][3], seats[0][i][4], seatmaster_id);
                        seat.appendTo($(".seattab" + key));
                    }
                    $.ajax({
                        url: "/getAllSeat",
                        type: 'POST',
                        data: {
                            id: seatmaster_id,
                        },
                    }).done(function (res) {
                        var res = JSON.parse(res);
                        $.each(res, function (key, value) {
                            appendMemberToSeat(value);
                        })
                    });

                    for (i = 0; i < seats[1].length; i++) {
                        var obj = createStaticObj(seats[1][i][1], seats[1][i][2], seats[1][i][3], seats[1][i][4], seats[1][i][5]);
                        obj.appendTo($(".seattab" + key));
                    }
                    for (i = 0; i < seats[2].length; i++) {
                        var door = createStaticDoor(seats[2][i][0], seats[2][i][1], seats[2][i][2]);
                        door.appendTo($(".seattab" + key));
                    }
                }
            });
        })
    }

    // Handle move member
    var check_click_child_box;
    var original_position;
    var move = true;
    var move_step2 = true;
    var swap_member;
    var arr_seat = [];
    var check_click_child_seat = 0;
    var role = $("#role").attr('data-role');
    $(".child").click(function () {
        if (role == 2 && !check_click_child_box) {
            if (move) {
                original_position = $(this).parent();
                handleMove(this, true);
            }
            else {
                // mouse move step 2
                check_click_child_seat += 1;
                if (move_step2 && check_click_child_seat == 1) {
                    handleMove(this, false);
                }

            }
        }
    });

    // event seat-added when move member from department to seat
    $(".seattab-content .child").on("seat-added", function (e, selected) {
        var seat_parent = $(selected).parent();
        var member_id = $(selected).attr("data-member-id");
        var table_number = $(seat_parent).closest("table").attr("data-table-id");
        var seatmaster_id = $(seat_parent).closest(".seattab").attr("data-seatmaster-id");
        var col = $(seat_parent).parent().children().index($(seat_parent));
        var row = $(seat_parent).parent().parent().children().index($(seat_parent).parent());
        $.ajax({
            url: "/admin/addSeat",
            type: "POST",
            data: {
                memberID: member_id,
                tableNumber: table_number,
                seatMasterID: seatmaster_id,
                row: row,
                column: col,
            },
        }).done(function () {
            var res = JSON.parse(getSeat(selected));
            $(selected).attr("data-seat-id", res.ID);
            sendDataSocket(res);
        })
    });

    // event seat-updated when move member from seat to seat
    $(".seattab-content .child").on("seat-updated", function (e, selected) {
        var seat_parent = $(selected).parent();
        var seat_id = $(selected).attr("data-seat-id");
        var member_id = $(selected).attr("data-member-id");
        var table_number = $(seat_parent).closest("table").attr("data-table-id");
        var seatmaster_id = $(seat_parent).closest(".seattab").attr("data-seatmaster-id");
        var col = $(seat_parent).parent().children().index($(seat_parent));
        var row = $(seat_parent).parent().parent().children().index($(seat_parent).parent());
        $.ajax({
            url: "/admin/editSeat",
            type: "POST",
            data: {
                seatID: seat_id,
                memberID: member_id,
                tableNumber: table_number,
                seatMasterID: seatmaster_id,
                row: row,
                column: col,
            },
        }).done(function () {
            var res = JSON.parse(getSeat(selected));
            sendDataSocket(res);
        })
    });

    // event seat-trashed when move member from seat to department
    $(".seattab-content .child").on("seat-trashed", function (e, selected) {
        var seat_id = $(selected).attr("data-seat-id");
        var department_id = $(selected).parent().attr("data-department-id");
        var member_id = $(selected).attr("data-member-id");
        var seatmaster_id = $(selected).closest(".seattab").attr("data-seatmaster-id");
        $.ajax({
            url: "/admin/deleteSeat",
            type: "POST",
            data: {
                seatID: seat_id,
            },
        }).done(function () {
            socket.send(
                JSON.stringify({
                    DepartmentID: parseInt(department_id),
                    MemberID: parseInt(member_id),
                    SeatMasterID: parseInt(seatmaster_id),
                    MoveType: 1,
                })
            );
        })
    });

    // send data to server use socket
    function sendDataSocket(res) {
        socket.send(
            JSON.stringify({
                ID: res.ID,
                TableNumber: res.TableNumber,
                SeatmasterID: res.SeatMasterID,
                MemberID: res.MemberID,
                Row: res.Row,
                Col: res.Col,
                MoveType: 1,
            })
        );
    }

    // append member to seat at another clients socket
    function appendMemberToSeat(data_seat) {
        var seatmaster_id = data_seat.SeatMasterID;
        var member_id = data_seat.MemberID;
        var seat_id = data_seat.ID;
        var table_id = data_seat.TableNumber;
        var row = data_seat.Row;
        var col = data_seat.Col;
        var position = $("#table" + table_id + "-" + seatmaster_id + " tr:eq(" + row + ") td:eq(" + col + ")");
        var item = $(".seattab[data-seatmaster-id = " + seatmaster_id + "]").find($(".child[data-member-id = " + member_id + "]"));
        item.attr("data-seat-id", seat_id);
        item.appendTo($(position));
    }

    // handle move member
    function handleMove(that, move_ok) {
        $(document).mousemove(function (e) {
            $(that).css({ "position": "absolute", "top": e.pageY + 10, "left": e.pageX + 10 });
            $(that).css({ "display": "block" });
            if ($(that).parent().attr("data-box") == "seat") {
                var offset_parent = $(that).parent().offset();
                $(that).css({ "position": "absolute", "top": e.pageY - offset_parent.top + 10, "left": e.pageX - offset_parent.left + 10 });
            }
            $(document).off('mouseup').on('mouseup', function (e) {
                var container = $(".seat-parent");
                var container2 = $(".box-department");
                var container3 = $(".box-department .child");
                var container4 = $(".department-name");
                // if the target of the click isn't the container nor a descendant of the container
                if (!container.is(e.target) && container.has(e.target).length === 0 && !container2.is(e.target) && container2.has(e.target).length === 0 || container4.is(e.target)) {
                    $(that).css({ "position": "unset" });
                    original_position.append($(that).next());
                    check_click_child_box = false;
                    move_step2 = true;
                    move = true;
                    arr_seat = [];
                    $(document).off('mousemove');
                    $(document).off('mouseup');
                } else if (container3.is(e.target) || container3.has(e.target).length) {
                    check_click_child_box = true;
                } else if (container.is(e.target) || container2.is(e.target)) {
                    var data_color;
                    data_color = ($(e.target).attr("data-color-dp") ? $(e.target).attr("data-color-dp") : $(that).attr("data-color-seat"));
                    if ($(that).attr("data-color-seat") == data_color) {
                        $(that).css({ "position": "unset", "display": "block" });
                        arr_seat.push($(that));
                        // move member to department
                        if ($(e.target).attr("data-box") == "department") {
                            $(e.target).append($(that));
                            $(that).css({ "position": "unset" });
                            // member selected inside department
                            if ($(original_position).attr("data-box") == "department") {
                                if (swap_member) {
                                    $(that).trigger("seat-added", arr_seat[0]);
                                    $(that).trigger("seat-trashed", arr_seat[1]);
                                    swap_member = false;
                                }
                            }
                            // member selected inside seat
                            else if ($(original_position).attr("data-box") == "seat") {
                                if (!swap_member) {
                                    $(that).trigger("seat-trashed", $(that));
                                } else {
                                    $(that).trigger("seat-updated", arr_seat[0]);
                                    $(that).trigger("seat-trashed", arr_seat[1]);
                                    swap_member = false;
                                }
                            }
                        }
                        // move member to seat
                        else if ($(e.target).attr("data-box") == "seat") {
                            if ($(e.target).children().length == 0) {
                                $(e.target).append($(that));
                                // member selected inside department
                                if ($(original_position).attr("data-box") == "department") {
                                    if (!swap_member) {
                                        $(that).trigger("seat-added", arr_seat[0]);
                                    } else {
                                        $(that).trigger("seat-added", arr_seat[0]);
                                        $(that).trigger("seat-updated", arr_seat[1]);
                                        swap_member = false;
                                    }
                                }
                                // member selected inside seat
                                else if ($(original_position).attr("data-box") == "seat") {
                                    if (!swap_member) {
                                        $(that).trigger("seat-updated", $(that));
                                    } else {
                                        $(that).trigger("seat-updated", arr_seat[0]);
                                        $(that).trigger("seat-updated", arr_seat[1]);
                                        swap_member = false;
                                    }
                                }
                            } else {
                                original_position.append($(that).next());
                            }
                        }
                        arr_seat = [];
                        check_click_child_box = false;
                        move_step2 = true;
                        move = true;
                        $(document).off('mousemove');
                        $(document).off('mouseup');
                    }
                } else {
                    if (move_ok) {
                        if (move_step2) {
                            // move step 2
                            $(e.target).css({ "display": "none" });
                            $(e.target).parent().append($(that));
                            $(that).css({ "position": "unset" });
                            arr_seat.push($(that));
                            check_click_child_box = false;
                            $(document).off('mousemove');
                            $(document).off('mouseup');
                            move_step2 = true;
                            move = false;
                            swap_member = true;
                            check_click_child_seat = 0;
                        }
                    } else {
                        // when click child of seat after swap
                        move = false;
                        move_step2 = false;
                        $(document).off('mouseup');
                    }
                }
            });
        });
    }

    // getSeat get seat from db
    function getSeat(selected) {
        var seat_parent = $(selected).parent();
        var member_id = $(selected).attr("data-member-id");
        var seatmaster_id = $(seat_parent).closest(".seattab").attr("data-seatmaster-id");
        var res = $.ajax({
            url: "/admin/getSeat",
            type: "POST",
            async: false,
            data: {
                seatMasterID: seatmaster_id,
                memberID: member_id,
            },
        });

        return res.responseText;
    }

})
