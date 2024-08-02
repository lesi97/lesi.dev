<?php

include '../../../private/loadEnv.php';

$id = $_ENV["SPOTIFY_ID"];
$secret = $_ENV["SPOTIFY_SECRET"];


$curl = curl_init();

curl_setopt_array($curl, [
  CURLOPT_URL => "https://accounts.spotify.com/api/token",
  CURLOPT_RETURNTRANSFER => true,
  CURLOPT_ENCODING => "",
  CURLOPT_MAXREDIRS => 10,
  CURLOPT_TIMEOUT => 30,
  CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
  CURLOPT_CUSTOMREQUEST => "GET",
  CURLOPT_POSTFIELDS => "grant_type=client_credentials&client_id=" . $id . "&client_secret=" . $secret,
  CURLOPT_HTTPHEADER => [
    "Content-Type: application/x-www-form-urlencoded"
  ],
]);

$response = curl_exec($curl);
$err = curl_error($curl);

function buildPage($pageTitle, $message) {
    $googleRecaptchaSiteKey = $_ENV["GOOGLE_RECAPTCHA_SITE_KEY"];
    echo "
    <!DOCTYPE html>
    <html>
        <head>
            <title>{$pageTitle}</title>
            <script src='https://www.google.com/recaptcha/api.js?render={$googleRecaptchaSiteKey}' async defer></script>
            
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #1c1c1c;
                    width: 100vw;
                    height: 100vh;
                    overflow: hidden;
                    color: white;
                    margin: 0;
                    font-size: 16px;
                }
                h1 {
                    font-size: 24px;
                }
                h2 {
                    font-size: 20px;
                }
                .message {
                    text-align: center;
                    position: absolute;
                    width: 50%;
                    height: 50%;
                    left: 50%;
                    top: 50%;
                    translate: -50% -50%;                            
                }
                a, button {
                    cursor: pointer;
                }
            </style>
        </head>
        <body>
                <div class='message'>
                    {$message}
                </div>
                <script>
                window.onload = function() {
                    grecaptcha.ready(function() {
                        grecaptcha.execute('{$googleRecaptchaSiteKey}', {action: 'login'}).then(function(token) {
                            var recaptchaResponse = document.getElementById('recaptchaResponse');
                            if(recaptchaResponse) {
                                recaptchaResponse.value = token;
                            }
                        });
                    });
                };                
            </script>
        </body>
    </html>
    ";
}

?>