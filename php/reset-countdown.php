<?php

function getNextTuesdayAtFivePM()
{
    $currentDateTime = new DateTime();
    $currentDayOfWeek = $currentDateTime->format('w');
    $daysUntilNextTuesday = (9 - $currentDayOfWeek) % 7;
    if ($daysUntilNextTuesday == 0 && $currentDateTime->format('H:i') >= '17:00') {
        $daysUntilNextTuesday = 7;
    }
    $nextTuesdayDateTime = (clone $currentDateTime)->modify("+{$daysUntilNextTuesday} days")->setTime(17, 0, 0);
    return $nextTuesdayDateTime;
}

$futureDateTime = getNextTuesdayAtFivePM();
$currentDateTime = new DateTime();

$interval = $currentDateTime->diff($futureDateTime);

if ($interval->invert == 1) {
    echo "Reset's here! Maybe gift a sub 👉👈";
} else {
    $countdownParts = [];

    if ($interval->days > 0) {
        $countdownParts[] = "{$interval->days} days";
    }

    $countdownParts[] = "{$interval->h} hours";
    $countdownParts[] = "{$interval->i} minutes";
    $countdownParts[] = "{$interval->s} seconds";

    $countdownString = implode(', ', $countdownParts);
    $finalCountdown = $countdownString . " until reset!";
    echo $finalCountdown;
}

?>
