$(window).on("scroll load", function () {
  if ($(window).scrollTop() > 1000) {
    $("#scrollTop").addClass("active");
  } else {
    $("#scrollTop").removeClass("active");
  }
});

$("#scrollTop").on("click", function (e) {
  e.preventDefault();
  $("html, body").animate({ scrollTop: 0 }, 1000);
});
