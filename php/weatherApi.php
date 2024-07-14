<?php

/*
error_reporting(E_ALL);
ini_set('display_errors', 1);
*/

header("Cache-Control: no-store, no-cache, must-revalidate, max-age=0");
header("Cache-Control: post-check=0, pre-check=0", false);
header("Pragma: no-cache");

include '../../private/loadEnv.php';

$api_key = $_ENV['WEATHER_APP'];

$latitude = $_GET['latitude'];
$longitude = $_GET['longitude'];

$url = 'http://api.weatherapi.com/v1/current.json?key=' . $api_key . '&q=' . $latitude . ',' . $longitude;

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($ch);

if (curl_errno($ch)) {
    echo 'Error:' . curl_error($ch);
}

curl_close($ch);

header('Content-Type: application/json');

echo $response;

?>