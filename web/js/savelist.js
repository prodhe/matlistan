/* savelist.js */

// URL base for showing a correct link
var base_url = "https://vargklippan.se/matlistan/store/";

function saveIngredients() {
    // get chosen dishes and ingredients
    var store = {dishes: [], ingredients: []};
	$('#userlist > li').each(function(i, li) {
        var text = $(li).text();
        if (text !== "") {
            store["dishes"].push(text);
        }
	});
	$('#ingredientslist > li').each(function(i, li) {
        var text = $(li).text();
        if (text !== "" && !$(li).hasClass("item-ignore")) {
            store["ingredients"].push(text);
        }
	});
    // stringify the data and send to PHP script
    var str_store = JSON.stringify(store);
    if (store["ingredients"].length > 0) {
        $.ajax({
            type: "POST",
            url: "store/savelist.php",
            data: {data: str_store},
            cache: false,
            success: function(response) {
                // save the resulting ID and print
                // a link for viewing the saved list
                var store_id = response;
                var link = base_url + "?" + store_id;
                $('#popup')
                    .html('<a href="' + link + '">' + link + '</a>')
                    .dialog({
                        title: "Sparad: " + store_id,
                        width: 400,
                        height: 100
                    });
            }
        });
    }
}
