/* food.js */

/* global storage for loaded json data */
window.data = null;

/* helper functions */
function sortAlpha(a, b) {
	return a.toLowerCase() > b.toLowerCase() ? 1 : -1;
}

/* make items in userlist sortable/draggable
   and update handler on change */
$(function($) {
	$('#userlist')
		.sortable({
			revert: true,
			placeholder: "placeholder"
		});

	// need to bind event function after initialization for
	// the outside trigger function to work properly...
	$('#userlist').bind("sortupdate", function(e, ui) {
		$('#userlist').draggable({
			connectWith: "#mainlist",
			remove: function(e, ui) {
				ui.remove();
			}
		});
		$('#userlist > li').dblclick(function() {
			$(this).remove();
			updateIngredients();
		});
		updateIngredients();
	});
});

/* fill main list with items */
function loadList() {
	window.data.forEach(function(item) {
		$('#mainlist').append(
			$('<li>').addClass("item")
					 .draggable({
						 cursor: 'move',
						 helper: 'clone',
						 connectToSortable: "#userlist"
					 })
					 .dblclick(function() {
						 $(this).clone().appendTo($('#userlist'));
						 $('#userlist').trigger("sortupdate");
					 })
					 .append(item["name"])
		).droppable({
			accept: "#userlist > li",
			over: function(e, ui) {
				$(this).fadeTo(200, 0.3);
			},
			out: function(e, ui) {
				$(this).fadeTo(200, 1.0);
			},
			drop: function(e, ui) {
				ui.draggable.remove();
				$(this).fadeTo(200, 1.0);
			}
		});
	});
};

/* update ingredients list */
function updateIngredients() {
	// clear current list
	$('#ingredientslist').children().remove();
	
	// get the names for current selected
	var dishes = [];
	$('#userlist > li').each(function(i, li) {
		dishes.push($(li).text());
	});

	// find and add ingredients for current selected
	var ingredients = [];
	dishes.forEach(function(dish) {
		window.data.forEach(function(d) {
			if (d["name"] == dish) {
				ingredients = ingredients.concat(d["ingredients"]);
			}
		});
	});

	// sort the ingredients list
	ingredients.sort(sortAlpha);

	// count each ingredient
	var count = {};
	ingredients.forEach(function(item) {
		if (count[item]) {
			count[item] += 1;
		}
		else {
			count[item] = 1;
		}
	});

	// remove duplicates from ingredients list,
	// now that we have a proper count
	ingredients = ingredients.filter(function(el,i,a){
		return (i==a.indexOf(el)) ? 1 : 0;
	});

	// add ingredients to the visual HTML list
	for (i=0; i<ingredients.length; i++) {
		var ing_name = ingredients[i];
		var txt = (count[ing_name] > 1) ? ing_name + " x" + count[ing_name] : ing_name;
		$('#ingredientslist').append($('<li>').append(txt));
	} 

	// for visual purposes, make sure there's always 5 rows in the list
	while ($('#ingredientslist > li').length < 5) {
		$('#ingredientslist').append($('<li>'));
	}	
}

/* onload */
$(document).ready(function() {
	/* load food data and add to main list */
	$.ajax({
		url: "./data.json",
		success: function(file) {
			//console.log("loaded JSON data");
			window.data = file["dishes"];
			window.data.sort(function(a,b) {
				return sortAlpha(a["name"],b["name"]);
			});
			loadList();
		}
	});
});
