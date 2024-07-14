<?php

session_start();
if (!isset($_SESSION["csrf_token"])) {
    $_SESSION["csrf_token"] = bin2hex(random_bytes(32));
};

include '../../../private/loadEnv.php';
include_once 'buildpage.php';

$secretKey = $_ENV["ANILIST_LOGIN_SECRET"];
#$googleRecaptchaSiteKey = $_ENV["GOOGLE_RECAPTCHA_SITE_KEY"];
#$googleRecaptchaSecretKey = $_ENV["GOOGLE_RECAPTCHA_SECRET_KEY"];

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


#if ($_SERVER["REQUEST_METHOD"] == "POST") { 
#
#    if (!isset($_POST["csrf_token"]) || $_POST["csrf_token"] !== $_SESSION["csrf_token"]) {        
#        http_response_code(401);
#        $loginInput .= '<p>CSRF token validation failed</p>';
#        die(buildPage($denied, $loginInput));
#    }
/*
    $captcha = $_POST['g-recaptcha-response'];
    $captchaVerifyUrl = "https://www.google.com/recaptcha/api/siteverify";
    $captchaData = [
        "secret" => $googleRecaptchaSecretKey,
        "response" => $captcha
    ];
    $options = [
        "http" => [
            "header" => "Content-type: application/x-www-form-urlencoded\r\n",
            "method" => "POST",
            "content" => http_build_query($captchaData)
        ]
    ];
    $context = stream_context_create($options);
    $response = file_get_contents($captchaVerifyUrl, false, $context);  
    $responseKeys = json_decode($response, true);

    $userInputKey = isset($_POST["key"]) ? htmlspecialchars($_POST["key"], ENT_QUOTES, 'UTF-8') : '';
    
    if ($responseKeys["success"] && $responseKeys["score"] >= 0.5) {
        if ($userInputKey !== $secretKey) {
            http_response_code(401);
            $loginInput .= '<p>Invalid Password. Please try again.</p>';
            die(buildPage($denied, $loginInput));
        }
    } else {
        http_response_code(401);
        $loginInput .= '<p>Invalid CAPTCHA. Please try again.</p>';
        die(buildPage($denied, $loginInput));
    }
*/
#} else {
#    http_response_code(401);
#    die(buildPage($denied, $loginInput));
#}


$twitchClientId = $_ENV["TWITCH_ID"];
$redirectURI = $_ENV["TWITCH_REDIRECT_URI"];
$scopes = [
    "moderator:read:chatters"
];
$scopeString = implode(" ", $scopes);

$query = [
    "client_id" => $twitchClientId,
	"redirect_uri" => $redirectURI,
    "response_type" => "code",
    'scope' => $scopeString   
];

$url = "https://id.twitch.tv/oauth2/authorize?" . http_build_query($query);

$pageTitle = "Login";
$loginButton = "
    <a href='{$url}'>
    <button>Login with Twitch</button>
    ";

buildPage($pageTitle, $loginButton)

?>