var teamSelectNumber = $("select#team_input").length

$(document).ready(function () {
    $('[data-toggle="datepicker"]').datepicker({
        format: 'mm-dd-yyyy'
    });

    // default member's avatar
    $("#delete-avatar").click(function () {
        $('#avatar').attr('src', '/webroot/img/avatar_empty.png');
        $("#check-avt").val("0");
        $("#file").val('')
    });

    // default team's image
    $("#delete-team-icon").click(function () {
        $('#avatar').attr('src', '/webroot/img/team_default.png');
        $("#check-avt").val("0");
        $("#file").val('')
    });

    $("#upload").click(function () {
        $("#check-avt").val("1");
        $("#file").click();
    });

    // close - open sidebar
    $('.close-icon').on('click', function () {
        $('#sidebar').toggleClass('active');
        $('.close-icon').toggleClass('active');
        $('#collapse-icon').toggleClass("fa-angle-double-left fa-angle-double-right");
        $(".main-content").toggleClass("main-toggle")
    });

    // update avatar right after upload image
    function readURL(input) {
        if (input.files && input.files[0]) {
            var reader = new FileReader();

            reader.onload = function (e) {
                $('#avatar').attr('src', e.target.result);
            }

            reader.readAsDataURL(input.files[0]);
        }
    }

    $("#file").change(function () {
        readURL(this);
    });

    // init DataTable
    $('#myTable').DataTable({
        "order": [],
        "language": { search: '', searchPlaceholder: "Search" },
        "drawCallback": function () {
            $(".dataTables_paginate > span > .paginate_button").addClass("btn-paginate");
            $(".dataTables_paginate > .paginate_button").addClass("btn-paginate");
            $(".dataTables_filter").addClass("mb-2 mt-3");
            $(".dataTables_length").addClass("mb-2 mt-3");
            $(".dataTables_filter > label > input").addClass("table-search");
            $(".dataTables_info").addClass("mt-2");
            $(".paging_simple_numbers").addClass("mt-2");
        }
    });

    $("#add_team").click(function () {
        if (teamSelectNumber <5) {
            $("#team_input").removeClass("is-invalid");
            $(".team-feedback").hide();
            $("#team_select").clone().addClass("mb-2 offset-sm-3")
                .appendTo("#team_add").after('<div class="btn btn-normal btn-delete remove-team mb-2"  onclick="removeTeam($(this))"><i class="fa fa-times"></i></div>');
            teamSelectNumber++;
        }
        if(teamSelectNumber >=5) {
            $("#add_team").addClass("disabled");
        }
    });

    var rowDepartments = $('#department-table tr').length

    // if have only 1 department (2 rows in table), disable delete department button
    if (rowDepartments < 3) {
        $(".btn-delete-department").prop("disabled", true)
    }

    // add new department
    $("#add-department").on("click", function () {
        // add one new row to department table
        $('#department-table > tbody').append(
            "<tr>" +
                "<td>" +
                    "<input type='text' name='departmentID' class='d-none'>" +
                    "<input type='text' name='departmentName' class='form-control custom-input'>" +
                "</td>" +
                "<td>" +
                    "<div class='row w-100 m-0'>" +
                    "<input type='text' name='departmentColor' class='form-control color-input custom-input'> " +
                    "<div class='show-color'></div>" +
                    "</div>" +
                "</td>" +
                "<td class='col-btn-delete'>" +
                    "<button class='btn btn-normal btn-delete btn-delete-department' onclick='removeDepartment($(this))'><i class='fa fa-times'></i></button>" +
                "</td>" +
            "</tr>"
        );

        // add id to new input and show color button
        var newColorInput = "color-input-" + rowDepartments
        var newColor = "show-color-" + rowDepartments
        $("#department-table tr:last").find(".color-input").attr('id', newColorInput)
        $("#department-table tr:last").find(".show-color").attr('id', newColor)

        // add js color to new department input
        var newColorPicker = document.getElementById(newColorInput)
        newColorPicker.jscolor = new jscolor(newColorPicker, { styleElement: newColor, required: false, hash: true })

        // increase rowDepartments by 1
        rowDepartments += 1

        // enable delete department button
        $(".btn-delete-department").prop("disabled", false)
    });
});

// confirm before "back" button
function confirmBack(data) {
    var r = confirm("Are you sure you want to clear the input information?");
    if (r == true) {
        if (data == "member") {
            window.location.replace("/admin/members");
        }
        if (data == "teams") {
            window.location.replace("/admin/teams");
        }
        if (data == "company") {
            window.location.replace("/admin/companies");
        }
        if (data == "team") {
            window.location.replace("/team");
        }
        if (data == "my_page") {
            window.location.replace("/my_page");
        }
    }
}

function removeTeam(obj) {
    teamSelectNumber--;
    if(teamSelectNumber <5) {
        $("#add_team").removeClass("disabled");
    }
    obj.prev().remove();
    obj.remove();
}

function removeDepartment(obj) {
    // remove the department
    obj.parent().parent().remove();

    // count number of departments in department table
    var countDepartment = $('#department-table tr').length - 1

    // if company has 1 department, disable delete department button
    if (countDepartment < 2) {
        $(".btn-delete-department").prop("disabled", true)
    }
}

function selectedCompany() {
    // reset department option
    $("#department > option").css("display", "none");
    $("#department").prop('selectedIndex', 0);

    // show departments of company
    var companyID = "companyID-" + $("#company option:selected").val();
    $("." + companyID).css("display", "block");

    // allow member has no department
    $("#department > option:first").css("display", "block");
}
