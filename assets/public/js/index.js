function setBundles() {
  const plugins = [];
  $(".plugin-list-item")
    .children("th")
    .each((index, element) => {
      var current = $(element);
      var txt = current.text();
      if (txt != "") {
        plugins.push(txt.trim());
      }
    });
  $("#pluginsHidden").val(plugins.join(","));
}

$(document).ready(() => {
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

  $(document).on("click", "#pluginRemove", (e) => {
    e.preventDefault();
    $(e.target).parents(".plugin-list-item").remove();
    setBundles();
  });

  $(document).on("submit", "#ftpForm", (e) => {
    setBundles();
  });

  $("#pluginAdd").click((e) => {
    e.preventDefault();
    const val = $("#pluginInput").val();
    $("#pluginsList").append(
      `<tr class="plugin-list-item">
        <th class="plugin-name" scope="row">
        ${val}
        </th>
        <td>
          <button id="pluginRemove">
          <i class="bi bi-backspace"></i>
          </button>
        </td>
        </tr>`
    );
    $("#pluginInput").val("");
    setBundles();
  });

  const stripe = Stripe(
    "pk_test_51IrBnXFNgmFlFbhMgWJyzrYXOEo9E8OOV33RAhccfGygNZT6fwOUhvcTk4UOlPekZoW3zQ3l4W9kAnUwsUwvqHnF00lXVnaiLM"
  );

  $("#checkoutButton").click(() => {
    stripe.redirectToCheckout({
      sessionId: $("#checkoutSessionId").val(),
    });
    // If `redirectToCheckout` fails due to a browser or network
    // error, display the localized error message to your customer
    // using `error.message`.
  });
});
