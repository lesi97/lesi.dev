<?php

include '../../../private/loadEnv.php';
include_once 'buildpage.php';
include_once 'auth.php';

include_once 'discord.php';

function handleOAuthToken($accessToken) {

    $clientId = $_ENV["SPOTIFY_ID"];
    $clientSecret = $_ENV["SPOTIFY_SECRET"];

    $redirectUri = "https://lesi.dev/php/spotify/oauth";
    $tokenUri = 'https://accounts.spotify.com/api/token';
    
    $fields = [
        "grant_type"    => "authorization_code",
        "client_id"     => $clientId,
        "client_secret" => $clientSecret,
        "redirect_uri"  => $redirectUri,
        "code"          => $accessToken,
    ];
    
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $tokenUri);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($fields));
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    
    $response = curl_exec($ch);
    curl_close($ch);

    if ($response) {
        #echo $response;
        sendDiscordNotification($response);
        $responseData = json_decode($response, true);
        if (isset($responseData["access_token"]) && 
            isset($responseData["refresh_token"]) &&
            isset($responseData["expires_in"])) {
                $pageMessage = '<h1>Success</h1>
                                <h2>You Can Now Close This Page</h2>';
                buildPage("Success", $pageMessage);
                storeTokens($responseData["access_token"], $responseData["refresh_token"], $responseData["expires_in"]);
        } else {
            buildPage("Invalid", '<h1>Invalid Response From Spotify</h1>');
        }
    } else {       
        buildPage("Failure", '<h1>CURL Failed</h1>');
    }
}


if(isset($_GET["code"])) {
    $accessToken = $_GET["code"];
    if ($accessToken !== "") {
        handleOAuthToken($accessToken);
    } else {
        header('Location: login.php');
        exit();
    }
} else {
    header('Location: login.php');
    exit();
}


?>
