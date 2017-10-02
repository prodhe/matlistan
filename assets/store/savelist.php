<?php

/* savelist.php 
 * 
 * inserts a list of dishes and ingredients
 * in JSON format and print the resulting
 * ID from the post in the database
 */

// path to database
$fn = "./db/store.sqlite";

// no need for messy html
//header("Content-type: text/plain");

// exit if no data is sent
if (empty($_POST['data'])) die("No POST data sent.");

// open database
$dbh = new SQLite3($fn) or die("ERROR: open database");

// escape the JSON-string for safe SQL insert
$data = $dbh->escapeString($_POST['data']);

// insert data
$q = "INSERT INTO lists (data) VALUES ('" . $data . "')";
$dbh->exec($q) or die("ERROR: inserting data");

// get the id of the newly inserted data and print it
$q = "SELECT id FROM lists WHERE data='" . $data . "'";
$res = $dbh->querySingle($q);
printf("%s", $res);

// close database
$dbh->close();

?>
