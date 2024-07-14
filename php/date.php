<?php

/*
error_reporting(E_ALL);
ini_set('display_errors', 1);
*/

header("Cache-Control: no-store, no-cache, must-revalidate, max-age=0");
header("Cache-Control: post-check=0, pre-check=0", false);
header("Pragma: no-cache");

$timezone = $_GET['timezone'];
echo $timezone;
    
if (in_array($timezone, timezone_identifiers_list())) {
   $dateTime = new DateTime('now', new DateTimeZone($timezone));
   echo $dateTime->format('l d F Y');
} else {
   echo "Invalid timezone.";
}
?>