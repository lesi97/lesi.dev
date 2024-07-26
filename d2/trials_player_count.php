<?php

function getPlayerCount() {
    $cacheFile = 'player_count_cache.json';
    // $cacheTime = 1800; // 30 minutes in seconds
    $cacheTime = 600; // 10 minutes in seconds

    if (file_exists($cacheFile)) {
        $cache = json_decode(file_get_contents($cacheFile), true);
        if ($cache && (time() - $cache['timestamp'] < $cacheTime)) {
            return $cache;
        }
    }

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
        return null;
    }

    $data = json_decode($response, true);
    if (!$data) {
        echo "Failed to decode JSON response";
        return null;
    }

    /*
    $playerCount = $data["hourly"]["breakdown"] ?? 0;

    $totalCount = 0;

    if (is_array($playerCount) && !empty($playerCount)) {
        foreach ($playerCount as $count) {
            if (is_array($count) && isset($count["count"])) {
                $totalCount += $count["count"];
            }
        }
    }
    
    //echo json_encode($data["platforms"]);
    */

    $totalCount = $data["platforms"]["0"]["recentStats"]["playerCount"];

    $cacheData = [
        'playerCount' => $totalCount,
        'timestamp' => time()
    ];
    file_put_contents($cacheFile, json_encode($cacheData));

    return $cacheData;
}

$cacheData = getPlayerCount();

if ($cacheData !== null) {
    $playerCount = $cacheData['playerCount'];
    $currentTime = time();
    $lastUpdated = floor(($currentTime - $cacheData['timestamp']) / 60);
    $nextUpdate = ceil(($cacheData['timestamp'] + 1800 - $currentTime) / 60);

    if ($playerCount !== 0) {
        echo "There are currently " . $playerCount . " players in Trials of Osiris across all platforms | Last updated: " . $lastUpdated . " minutes ago";
    } else {
        echo "No data yet, gift a sub then come back later";
    }
}
?>
