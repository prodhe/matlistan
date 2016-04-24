<?php

/* one-time use script for creating the sqlite3 database */

$fn = "./db/store.sqlite";

// no need for messy html
header("Content-type: text/plain");

// abort if already created
if (file_exists($fn)) die("The database $fn already exists.");

printf("%s\n\n", "Filename: $fn");

printf("%s\n", "Creating and opening database.");
$dbh = new SQLite3($fn) or die("ERROR: open");

printf("%s\n", "Creating table 'lists' (id, data).");
$q = "CREATE TABLE lists (id INTEGER PRIMARY KEY, data TEXT)";
$dbh->exec($q) or die("ERROR: create table");

printf("%s\n", "Closing database.");
$dbh->close();

printf("\n%s\n", "Done.");
?>
