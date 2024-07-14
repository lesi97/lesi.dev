<?php

$date = '2024-06-04 17:00:00';

$futureDateTime = new DateTime($date);
$currentDateTime = new DateTime();

$interval = $currentDateTime->diff($futureDateTime);

if ($interval->invert == 1) {
    //echo "The date $futureDate has already passed.\n";
    echo "The Final Shape is out! Maybe gift a sub 👉👈";
} else {

    $countdownParts = [];

    if ($interval->days > 0) {
        $countdownParts[] = "{$interval->days} days";
    }

    $countdownParts[] = "{$interval->h} hours";
    $countdownParts[] = "{$interval->i} minutes";
    $countdownParts[] = "{$interval->s} seconds";

    $countdownString = implode(', ', $countdownParts);

    //$finalCountdown = $countdownString . " until " . $futureDateTime->format('Y-m-d H:i:s');

    $finalCountdown = $countdownString . " until The Final Shape!";
    echo $finalCountdown;
}

//echo $currentDateTime->format('Y-m-d H:i:s');

?>