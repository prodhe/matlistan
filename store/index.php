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

// exit if not found
//if ($listid === 0) die("Not found.");

?>
<!doctype html>
<html lang="sv">
<head>
    <!--[if lt IE 9]>
        <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <link rel="stylesheet" href="../css/store.css" />
    <script src="https://code.jquery.com/jquery-1.11.3.min.js"></script>
<?php if ($listid !== 0) { ?>
    <script type="text/javascript">
    <!--
        var data = JSON.parse('<?=$res?>');
        $(document).ready(function() {
            data["dishes"].forEach(function(dish) {
                $('#dishes').append($('<li>').append(dish));
            });
            data["ingredients"].forEach(function(ingredient) {
                $('#ingredients')
                    .append(
                        $('<li>')
                        .click(function() {
                            var a = "item-ignore";
                            if ($(this).hasClass(a)) {
                                $(this).removeClass(a);
                            } else {
                                $(this).addClass(a);
                            }
                        }).append(ingredient)
                    );
            });
        });
    //-->
    </script>
<?php } else { ?>
    <script type="text/javascript">
    <!--
        function load() {
            var lid = $("#lid").val();
            if ($.isNumeric(lid)) {
                location.href = "./?" + lid;
            }
        }
    //-->
    </script>
<?php } ?>
    <title>Store</title>
</head>
<body>
<?php if ($listid === 0) { ?>
    <input type="text" size="6" id="lid" /> <br />
    <button id="send" onclick="load();">Ladda</button>
<?php } else { ?>
    <main>
        <section>
            <ul id="ingredients"></ul>
        </section>
        <section>
            <ul id="dishes"></ul>
        </section>
    </main>
    <footer>
    </footer>
<?php } ?>
</body>
</html>
