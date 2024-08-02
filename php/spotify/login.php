<?php

session_start();
if (!isset($_SESSION["csrf_token"])) {
    $_SESSION["csrf_token"] = bin2hex(random_bytes(32));
};

include '../../../private/loadEnv.php';
include_once 'buildpage.php';


$denied = "Access Denied";
$loginInput = "
    <h1>Input Code</h1>
    <form action='login' method='post'>
        <input type='password' name='key'>
        <input type='hidden' name='g-recaptcha-response' id='recaptchaResponse'>
        <input type='hidden' name='csrf_token' value='" . $_SESSION['csrf_token'] . "'>
        <input type='submit' value='Sumbit'>
    </form>
";   


$id = $_ENV["SPOTIFY_ID"];
$secret = $_ENV["SPOTIFY_SECRET"];
$redirectURI = "https://lesi.dev/php/spotify/oauth";
$scopes = [
    "user-read-currently-playing"
];
$scopeString = implode(" ", $scopes);

$query = [
    "client_id" => $id,
	"redirect_uri" => $redirectURI,
    "response_type" => "code",
    'scope' => $scopeString   
];

$url = "https://accounts.spotify.com/authorize?" . http_build_query($query);

$pageTitle = "Login";
$loginButton = "
    <a href='{$url}'>
    <button>Login with Spotify</button>
    ";

buildPage($pageTitle, $loginButton);
?>