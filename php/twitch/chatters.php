<?php



include 'auth.php';

/*
Command example:
!addcom !slap $(user) slaps tf out of $(eval a=decodeURIComponent(`$(querystring)`);a==``?`@$(urlfetch https://lesi.org.uk/php/twitch/chatters?channel=channel)`:a) TriHard 
*/


include 'discord.php';
$streamer = $_GET['channel'];
if (isset($streamer)) {
    getStreamerId($streamer);
}

function getStreamerId($streamer) {

    include '../../../private/loadTwitch.php';
    include '../../../private/loadEnv.php';

    $twitchClientId = $_ENV["TWITCH_ID"];
    $twitchClientSecret = $_ENV["TWITCH_SECRET"];
    $refreshToken = $_TWITCH["TWITCH_REFRESH_TOKEN"];
    $accessToken = $_TWITCH["TWITCH_ACCESS_TOKEN"];
    $refreshTokenExpireTime = $_TWITCH["TWITCH_REFRESH_TOKEN_EXPIRY"];

    if (isRefreshTokenExpired($refreshTokenExpireTime)) {
        $message = "Twitch Refresh Token Expired! Refreshing Now.";
        #sendDiscordNotification($message);
        $newTokens = refreshToken($twitchClientId, $twitchClientSecret, $refreshToken);
        $accessToken = $newTokens["access_token"];
        $refreshToken = $newTokens["refresh_token"];
    }


    $url = "https://api.twitch.tv/helix/users?login=" . $streamer;

    $apiCall = curl_init();
	curl_setopt($apiCall , CURLOPT_URL, $url);
	curl_setopt($apiCall , CURLOPT_RETURNTRANSFER, true);
	curl_setopt($apiCall , CURLOPT_FOLLOWLOCATION, true);
	curl_setopt($apiCall , CURLOPT_HTTPHEADER, array(
		'Client-ID: ' . $twitchClientId,
        'Authorization: Bearer ' . $accessToken
	));
	
    $apiCallResponse = curl_exec($apiCall);
	curl_close($apiCall);

    if ($apiCallResponse) {
        $decoded_data = json_decode($apiCallResponse, true);
    }

    $streamerID = $decoded_data["data"][0]["id"];
    if ($streamerID) {
        getChatters($twitchClientId, $accessToken, $streamerID, $streamer);   
    } 
}


function getChatters($twitchClientId, $accessToken, $streamerID, $streamerName) {   
    # Docs - https://dev.twitch.tv/docs/api/reference/#get-chatters
    #$streamerID = 477897; #terror
    $modID = 101129910; # me :)
    $streamUrl = "https://api.twitch.tv/helix/chat/chatters?first=1000&broadcaster_id=" . $streamerID . "&moderator_id=" . $modID;

    $curl = curl_init();

    curl_setopt_array($curl, [
        CURLOPT_URL => $streamUrl,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => "",
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 30,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => "GET",
        CURLOPT_HTTPHEADER => [
            "Accept: */*",
            "Authorization: Bearer " . $accessToken,
            "Client-ID: " . $twitchClientId
        ]
    ]);

    $response = curl_exec($curl);
    $err = curl_error($curl);
    
    curl_close($curl);

    if ($response) {
        $data = json_decode($response, true);
        $streamerArray = array(
            "user_id" => $streamerID,
            "user_login" => $streamerName,
            "user_name" => $streamerName
        );
        array_unshift($data['data'], $streamerArray);
        getRandomUser(json_encode($data), $streamerName);
    }

}


function getRandomUser($chatters, $streamer) {
    $data = json_decode($chatters, true);
    if (!isset($data['data']) || !is_array($data['data'])) {
        return "No users found";
    }

    $usernames = array();
    foreach ($data['data'] as $user) {
        if (isset($user['user_name'])) {
            $usernames[] = $user['user_name'];
        }
    }

    if (count($usernames) === 0) {
        return "No usernames found";
    }

    shuffle($usernames);

    $randomIndex = mt_rand(0, count($usernames) - 1);
    
    if (!$usernames[$randomIndex] == "") {
        echo $usernames[$randomIndex];
    } else {
        echo $streamer;
    }
    
}

?>