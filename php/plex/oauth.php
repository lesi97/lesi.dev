<?php

include '../../../private/loadEnv.php';
include_once 'buildpage.php';
include_once 'auth.php';

$code = $_GET["code"];

if(isset($code) && $code !== "") {
    handleOAuthToken($code);
} else {
    header('Location: login.php');
    exit();
}


function handleOAuthToken($accessToken) {

    $anilistClientId = $_ENV["ANILIST_CLIENT_ID"];
    $anilistClientSecret = $_ENV["ANILIST_CLIENT_SECRET"];

    $redirectUri = $_ENV["ANILIST_REDIRECT_URI"];
    $tokenUri = 'https://anilist.co/api/v2/oauth/token';
    
    $fields = [
        "grant_type"    => "authorization_code",
        "client_id"     => $anilistClientId,
        "client_secret" => $anilistClientSecret,
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
        $responseData = json_decode($response, true);
        if (isset($responseData["access_token"]) && 
            isset($responseData["refresh_token"]) &&
            isset($responseData["expires_in"])) {
                $pageMessage = '<h1>Success</h1>
                                <h2>You Can Now Close This Page</h2>';
                buildPage("Success", $pageMessage);
                storeTokens($responseData["access_token"], $responseData["refresh_token"], $responseData["expires_in"]);
        } else {
            buildPage("Invalid", '<h1>Invalid Response From Anilist</h1>' . $response);
        }
    } else {       
        buildPage("Failure", '<h1>CURL Failed</h1>');
    }
}

?>
