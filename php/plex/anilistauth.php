<?php

include_once '../../../private/loadAni.php';

function storeTokens($accessToken, $refreshToken, $expirySeconds) {
    $anilistEnv = "../../../private/anilist.env";
    $envContent = file($anilistEnv, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
    $newEnvContent = [];
    $expiryTime = time() + $expirySeconds;

    foreach ($envContent as $envLine) {
        if (strpos($envLine, 'ANILIST_ACCESS_TOKEN=') !== false) {
            $newEnvContent[] = "ANILIST_ACCESS_TOKEN=" . $accessToken;
        } elseif (strpos($envLine, 'ANILIST_REFRESH_TOKEN=') !== false) {
            $newEnvContent[] = "ANILIST_REFRESH_TOKEN=" . $refreshToken;
        } elseif (strpos($envLine, 'ANILIST_REFRESH_TOKEN_EXPIRY=') !== false) {
            $newEnvContent[] = "ANILIST_REFRESH_TOKEN_EXPIRY=" . $expiryTime;
        } else {
            $newEnvContent[] = $envLine;
        }
    }

    if (!in_array('ANILIST_REFRESH_TOKEN_EXPIRY=' . $expiryTime, $newEnvContent)) {
        $newEnvContent[] = "ANILIST_REFRESH_TOKEN_EXPIRY=" . $expiryTime;
    }

    file_put_contents($anilistEnv, implode(PHP_EOL, $newEnvContent));
}

function refreshToken($clientId, $clientSecret, $refreshToken) {
    $tokenUri = 'https://anilist.co/api/v2/oauth/token';

    $fields = [
        'grant_type'    => 'refresh_token',
        'client_id'     => $clientId,
        'client_secret' => $clientSecret,
        'refresh_token' => $refreshToken,
    ];

    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $tokenUri);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($fields));
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);

    $response = curl_exec($ch);
    curl_close($ch);
    
    $responseData = json_decode($response, true);
    $accessToken = $responseData["access_token"];
    $refreshToken = $responseData["refresh_token"];
    $refreshExpires = $responseData["expires_in"];

    storeTokens($accessToken, $refreshToken, $refreshExpires);

    return $responseData['access_token'];
}

function isRefreshTokenExpired($refreshToken) {
    $expiryTime = $_ANIENV["ANILIST_REFRESH_TOKEN_EXPIRY"]; 
    if (!$expiryTime) {
        return true;
    }
    $currentTime = time();
    return $currentTime > $expiryTime;
}

?>