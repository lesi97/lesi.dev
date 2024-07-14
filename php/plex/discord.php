<?php

include '../../../private/loadEnv.php';

function sendDiscordNotification($message) {

    $discord_webhook_url = $_ENV["DISCORD_WEBHOOK_URL"];
    $user_id = $_ENV["DISCORD_USER_ID"];
    
    $json_data = json_encode([
        "content" => "<@{$user_id}> ```{$message}```",
        "username" => "Anilist Notifications",
        "avatar_url" => "https://lesi.dev/php/plex/luffy.jpg"
    ]);

    $ch = curl_init($discord_webhook_url);
    curl_setopt($ch, CURLOPT_HTTPHEADER, array('Content-type: application/json'));
    curl_setopt($ch, CURLOPT_POST, 1);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $json_data);
    curl_setopt($ch, CURLOPT_FOLLOWLOCATION, 1);
    curl_setopt($ch, CURLOPT_HEADER, 0);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

    $response = curl_exec($ch);
    if($response === false) {
        echo 'Curl error: ' . curl_error($ch);
    }    
    curl_close($ch);

}

http_response_code(404);

?>