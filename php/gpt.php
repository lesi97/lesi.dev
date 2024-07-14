<?php

//$apiUrl = 'https://api.openai.com/v1/engines/text-davinci-003/completions';
$apiUrl = 'https://api.openai.com/v1/chat/completions';


include '../../private/loadEnv.php';
$apiKey = $_ENV['OPENAI'];

function callOpenAI($apiKey, $apiUrl, $prompt) {

    $predefinedResponseForTerror = "I don't know, gift a sub instead of asking all of these questions. Who are you the cops?";

    if (stripos($prompt, 'terror') !== false) {
        return ['choices' => [['message' => ['content' => $predefinedResponseForTerror]]]];
    }

    $predefinedContext = "I am a helpful assistant knowledgeable in gaming, streaming, and Twitch culture. My current location in Terror's twitch channel.";
    $data = [
        'model' => 'gpt-3.5-turbo', // Specify the model here
        'messages' => [['role' => 'user', 'content' => $prompt]]
    ];

    $ch = curl_init($apiUrl);
    curl_setopt($ch, CURLOPT_HTTPHEADER, [
        'Content-Type: application/json',
        'Authorization: Bearer ' . $apiKey
    ]);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));

    $response = curl_exec($ch);
    curl_close($ch);

    return json_decode($response, true);
}

if (isset($_GET['prompt'])) {
    $prompt = $_GET['prompt'];
    $response = callOpenAI($apiKey, $apiUrl, $prompt);

    // Check if response is valid and extract content
if (isset($response['choices'][0]['message']['content'])) {
    $content = $response['choices'][0]['message']['content'];
    echo substr($content, 0, 400); // Trim content to 400 characters
} else {
    echo 'No valid response received.';
}} else {
    echo 'No prompt provided.';
}?>