$(document).ready(function () {
    var arr = [0];
    $("table[data-type=seat]").each(function() {
        arr.push(parseInt($(this).attr('data-id')));
    });
    var max = Math.max(...arr);
    $("#room_design").attr('data-max',max);
    var arr_obj = [0];
    $(".div_table[data-type=obj]").each(function() {
        arr_obj.push(parseInt($(this).attr('data-id')));
    });
    var max_obj = Math.max(...arr_obj);
    $("#input_obj").attr('data-max',max_obj);
    if($("#check_page[data-page=add]").length > 0){
        $("#room_design").resizable({
            containment: "#outside_room_design",
        });
    }
    if($("#check_page[data-page=detail]").length > 0){
        var id = $("#check_page").attr('data-id');
        $.ajax({
            type: "POST",
            url: "/seats/getSeats",
            data: {
                id: id
            },
            success: function (res) {
                res = JSON.parse(res);
                var seats = JSON.parse(res.SeatMaster);
                $("#room_design").css({
                    width: seats[3][0] + 'px',
                    height: seats[3][1] + 'px',
                    margin: 'auto'
                });
                for (i = 0; i < seats[0].length; i++) {
                    var seat = createStaticTable(seats[0][i][0],seats[0][i][2],seats[0][i][1],seats[0][i][3],seats[0][i][4]);
                    seat.appendTo($("#room_design"));
                }
                for (i = 0; i < seats[1].length; i++){
                    var obj = createStaticObj(seats[1][i][1],seats[1][i][2],seats[1][i][3],seats[1][i][4],seats[1][i][5]);
                    obj.appendTo($("#room_design"));
                }
                for (i= 0; i < seats[2].length; i++) {
                    var door = createStaticDoor(seats[2][i][0],seats[2][i][1],seats[2][i][2]);
                    door.appendTo($("#room_design"));
                }
            }
        });
    }
    if($("#check_page[data-page=edit]").length > 0){
        $("#room_design").resizable({
            containment: "#outside_room_design",
        });
        var id = $("#check_page").attr('data-id');
        $.ajax({
            type: "POST",
            url: "/seats/getSeats",
            data: {
                id: id
            },
            success: function (res) {
                res = JSON.parse(res);
                var seats = JSON.parse(res.SeatMaster);
                $("#room_design").css({
                    width: seats[3][0] + 'px',
                    height: seats[3][1] + 'px',
                    margin: 'auto'
                });
                var arr = [0];
                var arr_obj = [0];
                for (i = 0; i < seats[0].length; i++) {
                    var seat = createTable(seats[0][i][2],seats[0][i][1]);
                    seat.css({
                        top: seats[0][i][3] + 'px',
                        left: seats[0][i][4] + 'px',
                    });
                    seat.children('table').attr('data-id', seats[0][i][0]);
                    seat.children('table').children().children().eq(0).text(seats[0][i][0]);
                    seat.appendTo($("#room_design"));
                    absoluteDrag();
                    arr.push(seats[0][i][0]);
                }
                var max = Math.max(...arr);
                $("#room_design").attr('data-max',Math.max(...arr));
                for (i = 0; i < seats[1].length; i++){
                    var obj = createObj();
                    obj.css({
                        position: 'absolute',
                        width: seats[1][i][1]+'px',
                        height: seats[1][i][2]+'px',
                        top: seats[1][i][3]+'px',
                        left: seats[1][i][4]+'px',
                    });
                    obj.attr('data-id',seats[1][i][0]);
                    obj.children('b').text(seats[1][i][5]);
                    obj.appendTo($("#room_design"));
                    arr_obj.push(seats[1][i][0]);
                }
                var max_obj = Math.max(...arr_obj);
                $("#input_obj").attr('data-max',max_obj);
                for (i= 0; i < seats[2].length; i++) {
                    var door = createDoor(seats[2][i][0]);
                    if (seats[2][i][0] == 't' || seats[2][i][0] == 'b'){
                        door.css({
                            left: seats[2][i][2] + 'px',
                        })
                    }else{
                        door.css({
                            top: seats[2][i][1] + 'px',
                        })
                    }
                    door.attr('data-ed',seats[2][i][0]);
                    door.appendTo($("#room_design"));
                }
            }
        });
    }
    var rows = 10;
    var cols = 10;
    var grid = '<div class="grid">';
    for (var i = 0; i < rows; i++) {
        grid += '<div class="roW">';
        for (var c = 0; c < cols; c++) {
            grid += '<div class="square"><div class="inner"></div></div>';
        }
        grid += '</div>';
    }
    grid += '</div>';

    var $gridChooser = $('.grid-chooser');
    var $grid = $(grid)
        .height(rows * 27)
        .width(cols * 27)
        .insertAfter($gridChooser)
    $grid.find('.roW').css({
        height: 'calc(100%/' + rows + ')',
        'margin-left': '0px',
    })
    $grid.find('.square').css({
        width: 'calc(100%/' + cols + ')'
    })

    var $allSquares = $('.square');

    $grid.on('mouseover', '.square', function () {
        var $this = $(this);
        var col = $this.index() + 1;
        var row = $this.parent().index() + 1;
        $allSquares.removeClass('highlight');
        $('.roW:nth-child(-n+' + row + ') .square:nth-child(-n+' + col + ')')
            .addClass('highlight');
        $gridChooser.val(col + ' x ' + row);
        $("#number_col").val(col);
        $("#number_row").val(row);
    })
    absoluteDrag();

    if($("#check_page[data-page=edit]").length > 0 || $("#check_page[data-page=add]").length > 0){
        $.contextMenu({
            selector: '#room_design',
            callback: function (key, options) {
                if (key == 'paste') {
                    var type = $("#room_design").attr('data-type');
                    if (type == 'obj') {
                        $(".context-menu-list").each(function () {
                            if ($(this).css("display") !== "none") {
                                var data_id = $("#room_design").attr('data-id');
                                var childPos = $(this).offset();
                                var parentPos = $("#room_design").offset();
                                var id_new = idObj();
                                var childOffset = {
                                    top: childPos.top - parentPos.top,
                                    left: childPos.left - parentPos.left,
                                }
                                var obj = $(".div_table[data-type=obj][data-id="+data_id+"]").clone();
                                obj.children(".ui-resizable-handle").remove();
                                obj.attr('data-id',id_new);
                                obj.mouseover(function () {
                                    $(this).css({
                                        border: '2px solid blue'
                                    });
                                }).mouseout(function () {
                                    $(this).css({
                                        border: '1px solid black'
                                    });
                                });
                                obj.draggable({
                                    containment: '#room_design',
                                    cursor: 'move',
                                    snap: "#room_design,.div_table",
                                    grid: [10, 10],
                                    drag: function () {
                                        $(this).find($('.topline')).css('display', 'block');
                                        $(this).find($('.rightline')).css('display', 'block');
                                        $(this).find($('.botline')).css('display', 'block');
                                        $(this).find($('.leftline')).css('display', 'block');
                                    },
                                    start: function () {
                                        $(this).find($('.topline')).css('display', 'block');
                                        $(this).find($('.rightline')).css('display', 'block');
                                        $(this).find($('.botline')).css('display', 'block');
                                        $(this).find($('.leftline')).css('display', 'block');
                                    },
                                    stop: function () {
                                        $(this).find($('.topline')).css('display', 'none');
                                        $(this).find($('.rightline')).css('display', 'none');
                                        $(this).find($('.botline')).css('display', 'none');
                                        $(this).find($('.leftline')).css('display', 'none');
                                    }
                                }).resizable();
                                obj.css({
                                    top: Math.round(childOffset.top / 10) * 10,
                                    left: Math.round(childOffset.left / 10) * 10,
                                }).appendTo($("#room_design"));
                                absoluteDrag();
                            }
                        })
                    }
                    if (type == 'seat') {
                        $(".context-menu-list").each(function () {
                            if ($(this).css("display") !== "none") {
                                var data_id = $("#room_design").attr('data-id');
                                var childPos = $(this).offset();
                                var parentPos = $("#room_design").offset();
                                var childOffset = {
                                    top: childPos.top - parentPos.top,
                                    left: childPos.left - parentPos.left,
                                }
                                var obj = $('table[data-type=seat][data-id=' + data_id + ']').parent().clone();
                                obj.mouseover(function () {
                                    $(this).css({
                                        border: '2px solid blue'
                                    });
                                }).mouseout(function () {
                                    $(this).css({
                                        border: '1px solid black'
                                    });
                                });
                                var id_new = idSeat();
                                obj.children('table').attr('data-id', id_new);
                                obj.children('table').children().children().eq(0).text(id_new);
                                obj.draggable({
                                    containment: '#room_design',
                                    cursor: 'move',
                                    grid: [10, 10],
                                    snap: "#room_design,.div_table",
                                    drag: function () {
                                        $(this).find($('.topline')).css('display', 'block');
                                        $(this).find($('.rightline')).css('display', 'block');
                                        $(this).find($('.botline')).css('display', 'block');
                                        $(this).find($('.leftline')).css('display', 'block');
                                    },
                                    start: function () {
                                        $(this).find($('.topline')).css('display', 'block');
                                        $(this).find($('.rightline')).css('display', 'block');
                                        $(this).find($('.botline')).css('display', 'block');
                                        $(this).find($('.leftline')).css('display', 'block');
                                    },
                                    stop: function () {
                                        $(this).find($('.topline')).css('display', 'none');
                                        $(this).find($('.rightline')).css('display', 'none');
                                        $(this).find($('.botline')).css('display', 'none');
                                        $(this).find($('.leftline')).css('display', 'none');
                                    }
                                });
                                obj.css({
                                    top: Math.round(childOffset.top / 10) * 10,
                                    left: Math.round(childOffset.left / 10) * 10,
                                }).appendTo($("#room_design"));
                                absoluteDrag();
                            }
                        })
                    }
                }
            },
            items: {
                "paste": { name: 'Paste', icon: "paste" }
            }
        })
    }

    $.contextMenu({
        selector: '.context-menu-one',
        callback: function (key, options) {
            if (key == 'edit') {
                var col = $(this).attr('data-col');
                var row = $(this).attr('data-row');
                $("#number_col").val(col);
                $("#number_row").val(row);
                $(".grid-chooser").val(col + ' x ' + row);
                var $allSquares = $('.square');
                $allSquares.removeClass('highlight');
                $('.roW:nth-child(-n+' + row + ') .square:nth-child(-n+' + col + ')').addClass('highlight');
                $("#seat_add").addClass('d-none');
                $("#seat_update").attr('data-id', $(this).attr('data-id')).removeClass('d-none');
                $("#room_design").attr('data-delete',$(this).attr('data-id'));
            }
            if (key == 'delete') {
                if($(this).attr('data-id') !=$("#room_design").attr('data-delete')){
                    $(this).parent().remove();
                }else{
                    $.toast({
                        heading: 'Error',
                        text: 'Can not delete edit seat',
                        icon: 'error',
                        position: 'bottom-right',
                        loader: false
                    });
                }
            }
            if (key == 'copy') {
                var type = $(this).attr('data-type');
                var data_id = $(this).attr('data-id');
                $("#room_design").attr({
                    'data-type': type,
                    'data-id': data_id,
                });
            }
        },
        items: {
            "edit": { name: 'Edit', icon: "edit" },
            "copy": { name: 'Copy', icon: "copy" },
            "delete": { name: 'Delete', icon: "delete" },
        }
    });

    $.contextMenu({
        selector: '.context-menu-two',
        callback: function (key, options) {
            if (key == 'delete') {
                var index = $(this).attr('data-id');
                if(index != $("#room_design").attr('data-edit')){
                    $(this).remove();
                }else{
                    $.toast({
                        heading: 'Error',
                        text: 'Can not delete edit object',
                        icon: 'error',
                        position: 'bottom-right',
                        loader: false
                    });
                }
            }
            if (key == 'copy') {
                var type = $(this).attr('data-type');
                var data_id = $(this).attr('data-id');
                $("#room_design").attr({
                    'data-type': type,
                    'data-id': data_id,
                });
            }
            if (key =='edit'){
                var txt = $(this).children('b').text();
                var index = $(this).attr('data-id');
                $("#obj_add").addClass('d-none');
                $("#obj_update").removeClass('d-none').attr('data-id',index);
                $("#input_obj").val(txt);
                $("#room_design").attr("data-edit",index);
            }
        },
        items: {
            "delete": { name: 'Delete', icon: "delete" },
            'copy': { name: 'Copy', icon: "copy" },
            'edit':{ name:'Change name', icon:"edit"}
        }
    });
    $.contextMenu({
        selector: '.right-click-door',

        callback: function (key, options) {
            if (key == 'delete') {
                $(this).remove();
            }
        },
        items: {
            "delete": { name: 'Delete', icon: "delete" }
        }
    });
    $("#seat_update").click(function () {
        var id = $(this).attr('data-id');
        $(this).attr('data-id', '');
        var col = $("#number_col").val();
        var row = $("#number_row").val();
        if (col === '' || row === '' || col == 0 || row == 0) {
            $.toast({
                heading: 'Error',
                text: 'Please input data for seat',
                icon: 'error',
                position: 'bottom-right',
                loader: false
            });
        } else {
            var length = $("#room_design").find($("[data-type=seat]")).length;
            var childPos = $("table[data-id=" + id + "]").offset();
            $("table[data-id=" + id + "]").parent().remove();
            var seat = createTable(row, col, 'update');
            seat.children('table').attr('data-id', id);
            seat.children('table').children().children().eq(0).text(id);
            var parentPos = $("#room_design").offset();
            var childOffset = {
                top: childPos.top - parentPos.top - 1,
                left: childPos.left - parentPos.left - 1
            }
            seat.css({
                top: childOffset.top,
                left: childOffset.left
            })
            seat.appendTo($('#room_design'));
            absoluteDrag();
            $("#number_col").val('1');
            $("#number_row").val('1');
            $("#room_design").attr('data-delete','');
            $(".grid-chooser").val('1 x 1');
            var $allSquares = $('.square');
            $allSquares.removeClass('highlight');
            $('.roW:nth-child(-n+' + 1 + ') .square:nth-child(-n+' + 1 + ')').addClass('highlight');
            $(this).addClass('d-none');
            $("#seat_add").removeClass('d-none');
        }
    });
    $("#obj_add").click(function () {
        var obj = createObj();
        var txt = $("#input_obj").val();
        obj.children('b').text(txt);
        obj.appendTo($("#room_design"));
        $("#input_obj").val('');
        absoluteDrag();
    });
    $("#obj_update").click(function(){
        var index = $(this).attr('data-id');
        var txt = $("#input_obj").val();
        $(".div_table[data-type=obj][data-id="+index+"]").children('b').text(txt);
        $(this).addClass('d-none');
        $("#obj_add").removeClass('d-none');
        $("#input_obj").val('');
        $("#room_design").attr("data-edit",'');
    })
    $("#seat_add").click(function () {
        var col = $("#number_col").val();
        var row = $("#number_row").val();
        if (col === '' || row === '' || col == 0 || row == 0) {
            $.toast({
                heading: 'Error',
                text: 'Please input data for seat',
                icon: 'error',
                position: 'bottom-right',
                loader: false
            });
        } else {
            var seat = createTable(row, col);
            seat.appendTo($("#room_design"));
            absoluteDrag();
            $("#number_col").val('1');
            $("#number_row").val('1');
            $(".grid-chooser").val('1 x 1');
            var $allSquares = $('.square');
            $allSquares.removeClass('highlight');
            $('.roW:nth-child(-n+' + 1 + ') .square:nth-child(-n+' + 1 + ')').addClass('highlight');
        }
    });
    $("#door_add").click(function () {
        var str = $("#door_input").val()
        var door = createDoor(str);
        door.attr('data-ed', str);
        door.appendTo($("#room_design"));
    });
    $('[data-toggle="tooltip"]').tooltip();
    $("#btn_save").click(function(){
        var arr = [], ar_seat = [], ar_obj = [], ar_door=[];
        var title = $("#input_title").val();
        if($.trim(title) == '' || title.length > 40) {
            $.toast({
                heading: 'Error',
                text: 'Title is wrong, please check!',
                icon: 'error',
                position: 'bottom-right',
                loader: false
            });
        }else{
            var validate = true;
            var company = $("#select_company").val();
            var width_room = $("#room_design").width();
            var height_room = $("#room_design").height();
            var ar_room = [width_room, height_room];
            $("table[data-type=seat]").each(function(){
                var id = $(this).attr('data-id');
                var col = $(this).attr('data-col');
                var row = $(this).attr('data-row');
                var top = $(this).parent().css('top').replace('px','');
                var left = $(this).parent().css('left').replace('px','');
                if(!validateSeat($(this), top, left)){
                    validate = false;
                }
                var ar = [id, col, row, top, left];
                ar_seat.push(ar);
            })
            $(".div_table[data-type=obj]").each(function(){
                var id = $(this).attr('data-id');
                var width = $(this).css('width').replace('px','');
                var height = $(this).css('height').replace('px','');
                var top = $(this).css('top').replace('px','');
                var left = $(this).css('left').replace('px','');
                if(!validateSeat($(this), top, left)){
                    validate = false;
                }
                var txt = $(this).children('b').text();
                var ar = [id, width, height, top, left, txt];
                ar_obj.push(ar);
            })
            $(".right-click-door").each(function(){
                var e = $(this).attr('data-ed');
                var top = $(this).css('top').replace('px','');
                var left = $(this).css('left').replace('px','');
                if(!validateSeat($(this), top, left)){
                    validate = false;
                }
                var ar = [e, top, left];
                ar_door.push(ar);
            })
            if(validate){
                arr.push(ar_seat);
                arr.push(ar_obj);
                arr.push(ar_door);
                arr.push(ar_room);
                arr = JSON.stringify(arr)
                $.ajax({
                    type: "POST",
                    url: "add",
                    data: {
                        title: title,
                        company: company,
                        seats: arr
                    },
                    success: function () {
                        window.location.replace("/admin/seats");
                    }
                });
            }else{
                $.toast({
                    heading: 'Error',
                    text: 'Something wrong, please check!',
                    icon: 'error',
                    position: 'bottom-right',
                    loader: false
                });
            }
        }
    })
    $("#btn_edit_seat").click(function(){
        var id = $("#check_page").attr('data-id');
        var arr = [], ar_seat = [], ar_obj = [], ar_door=[];
        var title = $("#input_title").val();
        if($.trim(title) == '' || title.length > 40) {
            $.toast({
                heading: 'Error',
                text: 'Title is wrong, please check!',
                icon: 'error',
                position: 'bottom-right',
                loader: false
            });
        }else{
            var validate = true;
            var company = $("#select_company").val();
            var width_room = $("#room_design").width();
            var height_room = $("#room_design").height();
            var ar_room = [width_room, height_room];
            $("table[data-type=seat]").each(function(){
                var id = $(this).attr('data-id');
                var col = $(this).attr('data-col');
                var row = $(this).attr('data-row');
                var top = $(this).parent().css('top').replace('px','');
                var left = $(this).parent().css('left').replace('px','');
                if(!validateSeat($(this), top, left)){
                    validate = false;
                }
                var ar = [id, col, row, top, left];
                ar_seat.push(ar);
            })
            $(".div_table[data-type=obj]").each(function(){
                var id = $(this).attr('data-id');
                var width = $(this).css('width').replace('px','');
                var height = $(this).css('height').replace('px','');
                var top = $(this).css('top').replace('px','');
                var left = $(this).css('left').replace('px','');
                if(!validateSeat($(this), top, left)){
                    validate = false;
                }
                var txt = $(this).children('b').text();
                var ar = [id, width, height, top, left, txt];
                ar_obj.push(ar);
            })
            $(".right-click-door").each(function(){
                var e = $(this).attr('data-ed');
                var top = $(this).css('top').replace('px','');
                var left = $(this).css('left').replace('px','');
                if(!validateSeat($(this), top, left)){
                    validate = false;
                }
                var ar = [e, top, left];
                ar_door.push(ar);
            })
            if(validate){
                arr.push(ar_seat);
                arr.push(ar_obj);
                arr.push(ar_door);
                arr.push(ar_room);
                arr = JSON.stringify(arr)
                $.ajax({
                    type: "POST",
                    url: "/admin/seats/edit",
                    data: {
                        id: id,
                        title: title,
                        company: company,
                        seats: arr,
                        status: $("#door_status").val(),
                    },
                    success: function () {
                        window.location.replace("/admin/seats");
                    }
                });
            }else{
                $.toast({
                    heading: 'Error',
                    text: 'Something wrong, please check!',
                    icon: 'error',
                    position: 'bottom-right',
                    loader: false
                });
            }
        }
    })
});

function createTable(row, col, type) {
    var x = 65;
    var y = 40;
    var temp = 0;
    var div = $("<div>").addClass('div_table').css({
        height: 'auto',
        width: 'auto',
        border: '1px solid black'
    });
    div.mouseover(function () {
        $(this).css({
            border: '2px solid blue'
        });
    }).mouseout(function () {
        $(this).css({
            border: '1px solid black'
        });
    });
    var span1 = $("<span>").addClass("topline").appendTo($(div));
    var span2 = $("<span>").addClass("rightline").appendTo($(div));
    var span3 = $("<span>").addClass("botline").appendTo($(div));
    var span4 = $("<span>").addClass("leftline").appendTo($(div));
    var table = $('<table>').addClass("context-menu-one").attr({
        'data-col': col,
        'data-row': row,
        'data-type': 'seat'
    }).appendTo($(div));
    for (m = 0; m < row; m++) {
        tr = $('<tr>').appendTo($(table));
        for (n = 0; n < col; n++) {
            td = $('<td>').css({
                height: y + 'px',
                width: x + 'px',
                border: '1px solid black',
            }).appendTo($(tr));
            var div1 = $('<div>').css({
                width: '100%',
                height: '100%',
            }).appendTo($(td));
            temp++;
        }
    }
    table.css({
        height: row * y + 'px',
        width: col * x + 'px',
    });
    div.draggable({
        containment: '#room_design',
        cursor: 'move',
        grid: [10, 10],
        snap: "#room_design,.div_table",
        drag: function () {
            $(this).find($('.topline')).css('display', 'block');
            $(this).find($('.rightline')).css('display', 'block');
            $(this).find($('.botline')).css('display', 'block');
            $(this).find($('.leftline')).css('display', 'block');
        },
        start: function () {
            $(this).find($('.topline')).css('display', 'block');
            $(this).find($('.rightline')).css('display', 'block');
            $(this).find($('.botline')).css('display', 'block');
            $(this).find($('.leftline')).css('display', 'block');
        },
        stop: function () {
            $(this).find($('.topline')).css('display', 'none');
            $(this).find($('.rightline')).css('display', 'none');
            $(this).find($('.botline')).css('display', 'none');
            $(this).find($('.leftline')).css('display', 'none');
        }
    });
    if (type != 'update') {
        var id = idSeat();
        table.attr('data-id', id);
        table.children().children().eq(0).text(id);
    }
    return div;
}
function createStaticTable(id, row, col, top, left, seatmaster_id){
    var x = 65;
    var y = 40;
    var temp = 0;
    var div = $("<div>").css({
        height: 'auto',
        width: 'auto',
        border: '1px solid black',
        top: top + 'px',
        left: left + 'px',
        position: 'absolute',
    });
    var table = $('<table id = "table' + id + '-' + seatmaster_id + '"class="bs-unset" data-table-id="'+ id +'">').appendTo($(div));
    for (m = 0; m < row; m++) {
        tr = $('<tr class="bs-unset">').appendTo($(table));
        for (n = 0; n < col; n++) {
            td = $('<td class="seat-parent bs-unset" data-box="seat">').css({
                height: y + 'px',
                width: x + 'px',
                border: '1px solid black',
            }).appendTo($(tr));
            temp++;
        }
    }

    return div;
}
function createObj() {
    var id = idObj();
    var div1 = $("<div>").css({
        'width': '60px',
        'height': '60px',
        'text-align': 'center',
        display: 'flex',
        border: '2px solid black',
        'background-color' : 'white'
    }).attr({
        'data-type': 'obj',
        'data-id': id,
    }).addClass('context-menu-two div_table');
    var span1 = $("<span>").addClass("topline").appendTo($(div1));
    var span2 = $("<span>").addClass("rightline").appendTo($(div1));
    var span3 = $("<span>").addClass("botline").appendTo($(div1));
    var span4 = $("<span>").addClass("leftline").appendTo($(div1));
    var div2 = $("<b>").css({
        margin: '0 auto',
        'align-self': 'center'
    }).appendTo($(div1));
    div1.mouseover(function () {
        $(this).css({
            border: '2px solid blue'
        });
    }).mouseout(function () {
        $(this).css({
            border: '1px solid black'
        });
    });
    div1.draggable({
        containment: '#room_design',
        cursor: 'move',
        snap: "#room_design,.div_table",
        grid: [10, 10],
        drag: function () {
            $(this).find($('.topline')).css('display', 'block');
            $(this).find($('.rightline')).css('display', 'block');
            $(this).find($('.botline')).css('display', 'block');
            $(this).find($('.leftline')).css('display', 'block');
        },
        start: function () {
            $(this).find($('.topline')).css('display', 'block');
            $(this).find($('.rightline')).css('display', 'block');
            $(this).find($('.botline')).css('display', 'block');
            $(this).find($('.leftline')).css('display', 'block');
        },
        stop: function () {
            $(this).find($('.topline')).css('display', 'none');
            $(this).find($('.rightline')).css('display', 'none');
            $(this).find($('.botline')).css('display', 'none');
            $(this).find($('.leftline')).css('display', 'none');
        }
    }).resizable({
        containment: "#outside_room_design",
    });
    return div1;
}
function createStaticObj(w, h, top, left, txt){
    var div1 = $("<div>").css({
        'width': w + 'px',
        'height': h + 'px',
        'text-align': 'center',
        position: 'absolute',
        top: top+ 'px',
        left: left + 'px',
        border: '2px solid black',
        display: 'flex',
        'background-color' : 'white',
    });
    var div2 = $("<b>").css({
        margin: '0 auto',
        'align-self': 'center'
    }).text(txt).appendTo($(div1));
    return div1;
}
function createDoor(e) {
    var rotate, top, bottom, right, left, axis, height, width, background;
    if (e == 'l') {
        left = '0';
        axis = 'y';
        height = '60px';
        width = '30px';
        background = "url('/webroot/img/door.png')";
    }
    if (e == 't') {
        top = '0';
        axis = 'x';
        height = '30px';
        width = '60px';
        background = "url('/webroot/img/door_2.png')";
    }
    if (e == 'r') {
        rotate = 'rotate(180deg)';
        right = '0';
        axis = 'y';
        height = '60px';
        width = '30px';
        background = "url('/webroot/img/door.png')";
    }
    if (e == 'b') {
        rotate = 'rotate(180deg)';
        bottom = '0';
        axis = 'x';
        height = '30px';
        width = '60px';
        background = "url('/webroot/img/door_2.png')";
    }
    var div1 = $("<div>").css({
        height: height,
        width: width,
        'background-image': background,
        'background-size': 'contain',
        transform: rotate,
        position: 'absolute',
        top: top,
        left: left,
        bottom: bottom,
        right: right
    }).addClass('right-click-door');
    div1.draggable({
        containment: '#room_design',
        cursor: 'move',
        snap: "#room_design,.div_table",
        axis: axis
    });
    return div1;
}
function createStaticDoor(e, top, left){
    var rotate, top, bottom, right, left, axis, height, width, background;
    if (e == 'l') {
        left = '0';
        axis = 'y';
        height = '60px';
        width = '30px';
        background = "url('/webroot/img/door.png')";
    }
    if (e == 't') {
        top = '0';
        axis = 'x';
        height = '30px';
        width = '60px';
        background = "url('/webroot/img/door_2.png')";
    }
    if (e == 'r') {
        rotate = 'rotate(180deg)';
        right = '0';
        axis = 'y';
        height = '60px';
        width = '30px';
        background = "url('/webroot/img/door.png')";
    }
    if (e == 'b') {
        rotate = 'rotate(180deg)';
        bottom = '0';
        axis = 'x';
        height = '30px';
        width = '60px';
        background = "url('/webroot/img/door_2.png')";
    }
    var div1 = $("<div>").css({
        height: height,
        width: width,
        position: 'absolute',
        top: top + 'px',
        left: left + 'px',
        'background-image': background,
        'background-size': 'contain',
        transform: rotate,
    });
    return div1;
}
function confirmBack() {
    var r = confirm("Are you sure you want to clear the input information?");
    if (r == true) {
        window.location.replace("/admin/seats");
    }
}
function absoluteDrag() {
    $('.div_table').each(function () {
        var top = $(this).position().top + 'px';
        var left = $(this).position().left + 'px';
        $(this).css({ top: top, left: left });
    }).css({ position: 'absolute' });
}
function idSeat() {
    var max = parseInt($("#room_design").attr('data-max')) + 1;
    $("#room_design").attr('data-max', max);
    return max;
}
function idObj(){
    var max = parseInt($("#input_obj").attr('data-max')) + 1;
    $("#input_obj").attr('data-max', max);
    return max;
}
function validateSeat(obj, top, left){
    var height = obj.height();
    var width = obj.width();
    var height_room = $("#room_design").height();
    var width_room = $("#room_design").width();
    if((height+parseInt(top) - 2 > height_room)||(width+parseInt(left) - 2 > width_room) || width > width_room){
        return false;
    }
    return true;
}
