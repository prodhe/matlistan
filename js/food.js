/*
 * food.js
 * */

// global storage for loaded json data
window.data = [];

// fill main list with items
function loadList(obj) {
    $('#mainlist').children().remove();
    var data = obj.dishes;
    window.data = data;
    data.forEach(function (item) {
        $('#mainlist')
            .append(
                $('<li>')
                .addClass("list-group-item")
                .draggable({
                    helper: 'clone',
                    addClasses: false,
                    stack: ".list-group-item"
                })
                .dblclick(function () {
                    addObjToList(this);
                })
                .append(safe2html(item.name))
            );
    });
}

// update ingredients list
function updateIngredients() {
    // clear current list
    $('#ingredientslist').children().remove();

    // get the names for current selected
    var dishes = [];
    $('#userlist > li').each(function (i, li) {
        dishes.push($(li).text());
    });

    // find and add ingredients for current selected
    var ingredients = [];
    dishes.forEach(function (dish) {
        window.data.forEach(function (d) {
            if (d.name === dish) {
                ingredients = ingredients.concat(d.ingredients);
            }
        });
    });

    // sort the ingredients list
    ingredients.sort(sortAlpha);

    // count each ingredient
    var count = {};
    ingredients.forEach(function (item) {
        if (count[item]) {
            count[item] += 1;
        } else {
            count[item] = 1;
        }
    });

    // remove duplicates from ingredients list, now that we have a proper count
    ingredients = ingredients.filter(function (el, i, a) {
        return (i === a.indexOf(el)) ? 1 : 0;
    });

    // add ingredients to the visual HTML list and add click function for
    // ignoring item before saving
    var i;
    for (i = 0; i < ingredients.length; i += 1) {
        var ing_name = ingredients[i];
        var txt = (count[ing_name] > 1) ? ing_name + " x" + count[ing_name] : ing_name;
        $('#ingredientslist').append($('<li>').click(function () {
            var a = "item-ignore";
            if ($(this).hasClass(a)) {
                $(this).removeClass(a);
            } else {
                $(this).addClass(a);
            }
        }).append(safe2html(txt)));
    }

    // for visual purposes, make sure there's always 5 rows in the list (or the
    // same as in the HTML file)
    while ($('#ingredientslist > li').length < 5) {
        $('#ingredientslist').append($('<li>'));
    }
}

// randomize the user list with n elements from main
function randomUserlist(n) {
    // abort if not numeric
    if (!($.isNumeric(n))) {
        return false;
    }

    clearUserlist();
    var mainlength = $('#mainlist > li').length;
    var rand = 0;
    var memory = [];
    var i;
    for (i = 0; i < n; i += 1) {
        // get random number that is between 1 and length of mainlist
        rand = Math.floor((Math.random() * mainlength) + 1);

        // if it's already in the list, add 1 (or circle around to beginning)
        // and check again
        while ($.inArray(rand, memory) !== -1) {
            // if we want more random items than length of mainlist, clear the
            // memory if it's already full
            if (memory.length === mainlength) {
                memory = [];
            }
            if (rand !== mainlength) {
                rand = rand + 1;
            } else {
                rand = 1;
            }
        }

        // finally, remember the item number for next run and add it to the
        // visual list
        addObjToList($('#mainlist > li:nth-child(' + rand + ')'));
        memory.push(rand);
    }
}

// onload
$(document).ready(function () {

    // // load food data and add to main list
    // $.ajax({
    //     url: "./data.json",
    //     success: function (file) {
    //         var a = file.dishes;
    //         for (var i=0; i<a.length; i++) {
    //             if ("name" in a[i] && "ingredients" in a[i]) {
    //                 // clear out empty string ingredients
    //                 a[i]["ingredients"] = a[i]["ingredients"].filter(function(el,i,a){
    //                     return (el != "") ? 1 : 0;
    //                 });
    //                 // add to global storage
    //                 window.data.push(a[i]);
    //             }
    //         }
    //         // sort main list
    //         window.data.sort(function(a,b) {
    //             return sortAlpha(a["name"],b["name"]);
    //         });
    //         // load to HTML
    //         loadList({"dishes":window.data});
    //     }
    // });

    // make items in userlist sortable
    $('#userlist')
        .sortable({
            revert: 100,
            forcePlaceholderSize: true
        })
        .droppable({
            accept: "#mainlist > li",
            over: function (e, ui) {
                $(this).fadeTo(200, 0.3);
            },
            out: function (e, ui) {
                $(this).fadeTo(200, 1.0);
            },
            drop: function (e, ui) {
                addObjToList(ui.draggable);
                ui.helper.remove();
                $(this).fadeTo(200, 1.0);
            }
        })
        .disableSelection();

    // make the mainlist a valid drop target for dragged items in userlist
    $('#mainlist')
        .droppable({
            accept: "#userlist > li",
            over: function (e, ui) {
                $(this).fadeTo(200, 0.3);
            },
            out: function (e, ui) {
                $(this).fadeTo(200, 1.0);
            },
            drop: function (e, ui) {
                ui.draggable.remove();
                $(this).fadeTo(200, 1.0);
            }
        })
        .disableSelection();

    // need to bind update event function after initialization for the outside
    // trigger function to work properly...
    $('#userlist').bind("sortupdate", function (e, ui) {
        $('#userlist > li').dblclick(function () {
            $(this).remove();
            updateIngredients();
        });
        updateIngredients();
    });

    // connect buttons
    $('#file_load_list').change(function (e) {
        loadExternalList(this);
    });
    $('#btn_random_userlist').click(function (e) {
        randomUserlist($('#random_dishes').val());
    });
    $('#btn_clear_userlist').click(clearUserlist);
    $('#btn_create_store_list').click(function (e) {
        saveIngredients(); //savelist.js
    });

    // activate bootstraps tooltip
    $('[data-toggle="tooltip"]').tooltip();
});
