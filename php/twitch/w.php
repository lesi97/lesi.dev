<?

include '../private/loadEnv.php';
$twitchClientId = $_ENV["TWITCH_ID"];
$redirectURI = $_ENV["TWITCH_REDIRECT_URI"];

echo $twitchClientId;


?>