<?php

function getTrialsData() {
    $curl = curl_init();

    curl_setopt_array($curl, [
        CURLOPT_URL => "https://api.trialsofthenine.com/weeks/0",
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => "",
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 30,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => "GET",
    ]);

    $response = curl_exec($curl);
    $err = curl_error($curl);

    curl_close($curl);

    if ($err) {
        echo "cURL Error #:" . $err;
    } else {
        file_put_contents('trials_data.json', $response);
    }
}

function displayTrialsData() {
    $jsonData = @file_get_contents('trials_data.json');

    if ($jsonData === false) {
        echo "No data available.";
        getTrialsData();
        displayTrialsData();
        return;
    }

    $data = json_decode($jsonData, true);

    $currentDate = new DateTime();
    $startDate = new DateTime($data['startDate']);
    $endDate = new DateTime($data['endDate']);

    if ($currentDate > $endDate) {
        getTrialsData();
        displayTrialsData();
    } else {
        $flawlessLoot = $data['rewards']['flawless'];
        $mapName = $data['map']['name'];
        echo "Map: " . $mapName . " | Flawless loot: " . $flawlessLoot . " & Random 3, 5 or 7 win drop | Adept Mod: Random | Chance at Ship, Sparrow & ghost";
    }
}

displayTrialsData();

?>
