/* printlist.js */

function printList() {
    // get chosen dishes and ingredients
    var store = {dishes: [], ingredients: []};
	$('#userlist > li').each(function(i, li) {
        if ($(li).text() !== "") {
            var dish = {
                name: $(li).children('.name').text(),
                category: $(li).children('.category').text(),
                ingredients: $(li).data("original-title")
            };
            store["dishes"].push(dish);
        }
	});
	$('#ingredientslist > li').each(function(i, li) {
        var text = $(li).text();
        if (text !== "" && !$(li).hasClass("item-ignore")) {
            store["ingredients"].push(text);
        }
	});

    if (store.ingredients.length > 0) {
        console.log(store);
        // hide all current sections
        $('section').hide();
        $('footer').hide();

        // add a new one displaying what to be printed, along
        $('main').append(
            $('<section>')
            .attr("id", "print")
            .dblclick(function () {
                this.remove();
                $('section').show();
                $('footer').show();
            })
            .append(function(){
                var obj = $('<dl>');
                store.dishes.forEach(function (dish) {
                    obj.append($('<dt>').append(dish.name));
                    obj.append($('<dd>').append(dish.ingredients));
                });
                return obj;
            })
            .append(function(){
                var obj = $('<p>');
                var str = "";
                store.ingredients.forEach(function (ing) {
                    str += ing + "<br />";
                });
                obj.append(str);
                return obj;
            })
        );
    }
}
