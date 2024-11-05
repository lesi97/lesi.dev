<?php

include_once 'auth.php';
include_once 'discord.php';

$plexPayload = $_POST['payload'];
if (isset($plexPayload)) plexSearch($plexPayload);

function plexSearch($plexPayload) {
    include '../../../private/loadEnv.php';
    include '../../../private/loadAni.php';

    $anilistClientId = $_ENV["ANILIST_CLIENT_ID"];
    $anilistClientSecret = $_ENV["ANILIST_CLIENT_SECRET"];
    $refreshToken = $_ANIENV["ANILIST_REFRESH_TOKEN"];
    $accessToken = $_ANIENV["ANILIST_ACCESS_TOKEN"];    
    $refreshTokenExpireTime = $_ANIENV["ANILIST_REFRESH_TOKEN_EXPIRY"]; 

    if (isRefreshTokenExpired($refreshTokenExpireTime)) {
        $message = "Anilist Refresh Token Expired! Refreshing Now. ";
        sendDiscordNotification($message);
        $newTokens = refreshToken($anilistClientId, $anilistClientSecret, $refreshToken);
        $accessToken = $newTokens["accessToken"];
        $refreshToken = $newTokens["refreshToken"];
    }

    $anilistUrl = "https://graphql.anilist.co";
    
    $plexData = json_decode($_POST['payload'], true);	
    $event = $plexData["event"];
    $mediaType = $plexData["Metadata"]["librarySectionTitle"];
    $mediaTypeRegex = "/Anime/i";
    $plexUsername = $plexData["Account"]["title"];
    
    if ($event === "media.scrobble") {
        #file_put_contents("scrobble.txt", $plexPayload);
    }
    
    if ($event === "media.scrobble" && 
        preg_match($mediaTypeRegex, $mediaType) &&
        $plexUsername === "C_Lesi") {	
            
            if (isset($plexData["Metadata"]["parentTitle"]) && strpos($plexData["Metadata"]["parentTitle"], "Season") === false) {
                $showName = $plexData["Metadata"]["parentTitle"];
            } elseif (isset($plexData["Metadata"]["parentTitle"]) && strpos($plexData["Metadata"]["parentTitle"], "Season") !== false) {
                $showName = $plexData["Metadata"]["grandparentTitle"] . " " . $plexData["Metadata"]["parentTitle"];
            } else {
                $showName = $plexData["Metadata"]["grandparentTitle"];
            }         
            
            if ($plexData["Metadata"]["type"] === "episode") {
                $episodeNumber = $plexData["Metadata"]["index"];
            } else {
                $episodeNumber = 1;
            }    
            
            if ($episodeNumber % 100 == 0) {
                $message = "Anilist updated! Episode {$episodeNumber} on {$showName} just watched.";
                sendDiscordNotification($message);
            }
            
                
            $searchQuery = '
                query ($title: String) {
                    Media (search: $title, type: ANIME) {
                        id
                        episodes
                        endDate {
                            year
                            month
                            day
                        }
                    }
                }';

            $searchVariables = ['title' => $showName];

            $payload = json_encode(['query' => $searchQuery, 'variables' => $searchVariables]);

            $payload_pretty = json_encode($payload, JSON_PRETTY_PRINT);
            $searchQueryPayload = "Search Query Payload:\n" . $payload_pretty;
            #sendDiscordNotification($searchQueryPayload);


            $ch = curl_init();
            curl_setopt($ch, CURLOPT_URL, $anilistUrl);
            curl_setopt($ch, CURLOPT_HTTPHEADER, [
                "Content-Type: application/json",
                "Authorization: Bearer {$accessToken}"
            ]);
            curl_setopt($ch, CURLOPT_POST, true);
            curl_setopt($ch, CURLOPT_POSTFIELDS, $payload);
            curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
                
            $searchResult = curl_exec($ch);

            $searchQueryResult = "Search Query Result:\n" . $searchResult;
            #sendDiscordNotification($searchQueryResult);

            if ($searchResult === false) {
                $error = curl_error($ch);
                $apiError = "API Request Failed:\n" . $error;
                sendDiscordNotification($apiError);
                curl_close($ch);
                return;
            } 
            
            $searchResponse = json_decode($searchResult, true);

            if (isset($searchResponse['data']['Media'])) {
                #sendDiscordNotification($searchResponse);
                $maxEpisodeCount = $searchResponse['data']['Media']['episodes'];
                $anilistID = $searchResponse['data']['Media']['id'];
                updateAnilist($anilistID, $episodeNumber, $accessToken, $showName, $maxEpisodeCount);                
            }

           curl_close($ch);

        } else {                
            echo "Not expected user or not episode or not scrobble";
        }
} 



function updateAnilist($animeId, $episodeNumber, $accessToken, $showName, $maxEpisodeCount) {
    
    $mutation = '
        mutation ($mediaId: Int, $progress: Int, $status: MediaListStatus) {
            SaveMediaListEntry (mediaId: $mediaId, progress: $progress, status: $status ) {
                id
                progress
                status
            }
        }';

    if ($episodeNumber !== $maxEpisodeCount) {
        $variables = [
            'mediaId' => $animeId,
            'progress' => $episodeNumber,
            'status' => "CURRENT"
        ];
    } else {
        $variables = [
        'mediaId' => $animeId,
        'progress' => $episodeNumber,
        'status' => "COMPLETED"
    ];

    }
    

    $payload = json_encode(['query' => $mutation, 'variables' => $variables]);

    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, 'https://graphql.anilist.co');
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        "Content-Type: application/json",
        "Authorization: Bearer {$accessToken}"
    ]);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $payload);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

    $result = curl_exec($ch);
    curl_close($ch);

    if ($result === false) {
        return false;
    }

    $response = json_decode($result, true);
    if (isset($response['errors'])) {
        sendDiscordNotification($result);
        return false;
    }

    $message = "Anilist updated! Episode {$episodeNumber} on {$showName} just watched.";
    #sendDiscordNotification($message);
    return true;
}

?>