<?php

// open database
$fn = "./db/store.sqlite";
$dbh = new SQLite3($fn) or die("Could not open database.");

// get id from URL
$listid = 0;
$res = "";
foreach(array_keys($_GET) as $key) {
    if (is_numeric($key)) {
        $q = "SELECT data FROM lists WHERE id='" . $key . "'";
        $res = $dbh->querySingle($q);
        if (!empty($res)) $listid = $key;
    }
}

// check for post if not found in URL
if ($listid === 0) {
}

?>
<!doctype html>
<html lang="sv">
<head>
    <!--[if lt IE 9]>
        <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
    <meta charset="UTF-8" />
    <meta name="viewport" content="initial-scale=1" />
    <title>Store</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
    <link rel="stylesheet" href="../css/store.css" />
    <script src="https://code.jquery.com/jquery-1.11.3.min.js"></script>
    <script src="../js/helpers.js"></script>
</head>
<body>
<?php if ($listid === 0) { ?>
<div class="container-fluid">
    <div class="row-fluid text-center">
        <div class="col-xs-8 col-xs-offset-2">
        <form action="./" method="get" class="form-inline">
            <div class="form-group">
                <div class="input-group">
                    <div class="input-group-addon">
                        <span class="glyphicon glyphicon-shopping-cart"></span>
                        <label for="lid" class="sr-only">ID:</label>
                    </div>
                    <input type="text" size="4" id="lid" class="form-control input-lg" placeholder="0" />
                    <div class="input-group-btn">
                        <button type="submit" id="send" class="btn btn-primary btn-lg">
                            <span class="glyphicon glyphicon-log-in"></span>
                            <span class="sr-only">Ladda</span>
                        </button>
                    </div>
                </div>
            </div>
        </form>
        </div>
    </div>
</div>
<script type="text/javascript">
$('form').submit(function() {
        var lid = $("#lid").val();
        if ($.isNumeric(lid)) {
            location.href = "./?" + lid;
        }
        return false;
});
</script>
<?php } else { ?>
<div id="header" class="btn btn-block btn-lg btn-primary">
    <span class="glyphicon glyphicon-shopping-cart"></span>
    <?=$listid?>
</div>
<ul id="ingredients" class="list-group"></ul>
<div class="well">
    <ul id="dishes" class="list-unstyled"></ul>
</div>
<?php
    // pre-parse for safer JSON parsing
    $res = str_replace("<","\\u003c",addslashes($res));
?>
<script type="text/javascript">
    var data = JSON.parse('<?=$res?>');
    $(document).ready(function() {
        data["dishes"].forEach(function(dish) {
            //$('#dishes').append($('<li>').addClass("list-group-item").append(dish));
            $('#dishes').append($('<li>').append(safe2html(dish)));
        });
        data["ingredients"].forEach(function(ingredient) {
            $('#ingredients')
                .append(
                    $('<li>')
                    .addClass("list-group-item")
                    .click(function() {
                        $(this).toggleClass("item-ignore");
                        $(this).toggleClass("text-muted");
                    }).append(safe2html(ingredient))
                );
        });
    });
</script>
<?php } ?>
</body>
</html>
