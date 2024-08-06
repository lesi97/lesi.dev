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

    $totalCount = $data["platforms"]["0"]["recentStats"]["playerCount"];
    $lastUpdated = $data["platforms"]["0"]["recentStats"]["updatedAt"];

    $cacheData = [
        'playerCount' => $totalCount,
        'updatedAt' => $lastUpdated,
        'timestamp' => time(),
        'startDate' => $data['startDate'],
        'endDate' => $data['endDate']
    ];
    file_put_contents($cacheFile, json_encode($cacheData));

    return $cacheData;
}

$cacheData = getPlayerCount();

if ($cacheData !== null) {
    $playerCount = $cacheData['playerCount'];
    $currentTime = time();
    
    $lastUpdatedTimestamp = strtotime($cacheData['updatedAt']);

    $twoHoursInSeconds = 2 * 60 * 60;
    
    $elapsedTime = $currentTime - $lastUpdatedTimestamp;
    $elapsedMinutes = floor($elapsedTime / 60);
    $nextUpdate = ceil(($cacheData['timestamp'] + 1800 - $currentTime) / 60);

    $currentDate = new DateTime("now", new DateTimeZone("UTC"));

    $week = $currentDate->format("W");
    $year = $currentDate->format("Y");
    $start = new DateTime("{$year}-W{$week}-5 17:00:00", new DateTimeZone("UTC")); // 5 PM UTC Friday
    $end = new DateTime("{$year}-W{$week}-2 17:00:00", new DateTimeZone("UTC")); // 5 PM UTC Tuesday

    if ($currentDate < $start) {
        $start->modify('-1 week');
    }


    if (($currentTime - $lastUpdatedTimestamp) > $twoHoursInSeconds) {
        if (file_exists('player_count_cache.json')) {
            unlink('player_count_cache.json');
        }
        echo "trials isn't here yet dummy, gift a sub and come back later";
    } else {
        if ($elapsedMinutes < 60) {
            echo "There are currently " . number_format($playerCount) . " players in Trials of Osiris across all platforms | Last updated: " . $elapsedMinutes . " minutes ago";
        } else {
            $hours = floor($elapsedMinutes / 60);
            $minutes = $elapsedMinutes % 60;
            $timeString = "Last updated: " . $hours . " hour" . ($hours != 1 ? "s" : "");
        
            if ($minutes > 0) {
                $timeString .= " and " . $minutes . " minute" . ($minutes != 1 ? "s" : "");
            }
            echo "There are currently " . number_format($playerCount) . " players in Trials of Osiris across all platforms | " . $timeString . " ago";
        }
    }
}
?>
