/* printlist.js */

// URL base for showing a correct link

function printList() {
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

    if (store["ingredients"].length > 0) {
        console.log(window.data);
        // $('main').html(store);
    }
}
