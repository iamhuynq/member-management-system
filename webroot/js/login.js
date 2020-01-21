$(document).ready(function () {
    // check firefox's autofill
    if ($(".md-form input").val()) {
        $(".md-form input").each(function () {
            if ($(this).val().length > 0) {
                $(this).next().addClass("active");
            }
        });
    }

    $(".md-form input").on("change focus", function () {
        $(this).next().addClass("active");
        $(this).removeClass("is-invalid");
    });

    $(".md-form input").focusout(function () {
        var value = $(this).val();
        if (value.length == 0) {
            $(this).next().removeClass("active");
        }
    });
})
