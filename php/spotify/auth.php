<?php

include_once '../../../private/loadEnv.php';

function storeTokens($accessToken, $refreshToken, $expirySeconds) {
    $env = "../../../private/spotify.env";

    $envContent = file($env, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
    if ($envContent === false) {
        // Handle file read error
        echo "Error reading file";
        return;
    }

    $newEnvContent = [];
    $expiryTime = time() + $expirySeconds;
    $hasAccessToken = $hasRefreshToken = $hasExpiry = false;

    foreach ($envContent as $envLine) {
        if (strpos($envLine, 'SPOTIFY_ACCESS_TOKEN=') !== false) {
            $newEnvContent[] = "SPOTIFY_ACCESS_TOKEN=" . $accessToken;
            $hasAccessToken = true;
        } elseif (strpos($envLine, 'SPOTIFY_REFRESH_TOKEN=') !== false) {
            $newEnvContent[] = "SPOTIFY_REFRESH_TOKEN=" . $refreshToken;
            $hasRefreshToken = true;
        } elseif (strpos($envLine, 'SPOTIFY_REFRESH_TOKEN_EXPIRY=') !== false) {
            $newEnvContent[] = "SPOTIFY_REFRESH_TOKEN_EXPIRY=" . $expiryTime;
            $hasExpiry = true;
        } else {
            $newEnvContent[] = $envLine;
        }
    }

    if (!$hasAccessToken) {
        $newEnvContent[] = "SPOTIFY_ACCESS_TOKEN=" . $accessToken;
    }
    if (!$hasRefreshToken) {
        $newEnvContent[] = "SPOTIFY_REFRESH_TOKEN=" . $refreshToken;
    }
    if (!$hasExpiry) {
        $newEnvContent[] = "SPOTIFY_REFRESH_TOKEN_EXPIRY=" . $expiryTime;
    }

    $result = file_put_contents($env, implode(PHP_EOL, $newEnvContent));
    if ($result === false) {
        echo "Error writing to file";
        return;
    }
}


function refreshToken($clientId, $clientSecret, $refreshToken) {
    $tokenUri = 'https://accounts.spotify.com/api/token';

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

    return $responseData;
}

function isRefreshTokenExpired($refreshToken) {
    $expiryTime = $refreshToken;
    if (!$expiryTime) {
        return true;
    }
    $currentTime = time();
    return $currentTime > $expiryTime;
}

?>