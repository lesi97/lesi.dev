<?php
header("Cache-Control: no-store, no-cache, must-revalidate, max-age=0");
header("Cache-Control: post-check=0, pre-check=0", false);
header("Pragma: no-cache");

    $timezone = $_GET['timezone'];
    
    if (in_array($timezone, timezone_identifiers_list())) {
        $dateTime = new DateTime('now', new DateTimeZone($timezone));
        $dateTime->modify('-38 seconds'); // -38 seconds as it's 38 seconds too fast on the server
        echo $dateTime->format('H:i:s');
    } else {
        echo "Invalid timezone.";
    }

?>