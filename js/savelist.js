/* savelist.js */

// args: 1-level array of ingredients text string
function save_ingredients() {
    var ingredients = [];
	$('#ingredientslist > li').each(function(i, li) {
        var text = $(li).text();
        if (text !== "") {
            ingredients.push(text);
        }
	});
    console.log(ingredients);
    if (ingredients.length > 0) {
        console.log("time to save");
        $.ajax({
            type: "POST",
            url: "store/savelist.php",
            data: {data: ingredients},
            cache: false,
            success: function(response) {
                console.log("sent to cgi");
                console.log(response);
            }
        });
    }
}
