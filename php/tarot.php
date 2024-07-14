<?php

    $json = "tarot.json";
    $jsonData = file_get_contents($json);
    $tarotRawData = json_decode($jsonData, true);	

    $tarotCards = array();
    foreach ($tarotRawData["cards"] as $card) {
        if (isset($card["name"])) {
            $tarotCards[] = [$card["name"], $card["desc"]];
        }
    }

    shuffle($tarotCards);

    $randomIndex = mt_rand(0, count($tarotCards) - 1);
    
    $message = $tarotCards[$randomIndex][0] . " | " . $tarotCards[$randomIndex][1];

    if (mb_strlen($message) > 400) {
        $message = mb_substr($message, 0, 400);
    }

    echo $message;

?>