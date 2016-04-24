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
if ($listid === 0) die("Not found.");

?>
<!doctype html>
<html lang="sv">
<head>
    <!--[if lt IE 9]>
        <script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
    <meta charset="UTF-8" />
    <script src="https://code.jquery.com/jquery-1.11.3.min.js"></script>
    <script type="text/javascript">
    <!--
        var data = JSON.parse('<?=$res?>');
        console.log(data);
        $(document).ready(function() {
            data["dishes"].forEach(function(dish) {
                console.log(dish);
                $('#dishes').append($('<li>').append(dish));
            });
            data["ingredients"].forEach(function(ingredient) {
                console.log(ingredient);
                $('#ingredients').append($('<li>').append(ingredient));
            });
        });
    //-->
    </script>
    <link rel="stylesheet" href="../css/store.css" />
    <title>Store</title>
</head>
<body>
    <header>
        <h1><?=$listid?></h1>
    </header> 
    <main>
        <section>
            <ul id="dishes">
            </ul>
        </section>
        <section>
            <ul id="ingredients">
            </ul>
        </section>
    </main>
    <footer>
    </footer>
</body>
</html>
