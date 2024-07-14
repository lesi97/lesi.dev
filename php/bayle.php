<?php
header('Content-Type: application/json');

// Define your messages
$messages = [
    "Message 1: Curse you Bayle!",
    "Message 2: I hereby vow! You will rue this day!",
    "Message 3: Behold, a true drake warrior! And I, Igon!"
];

// Select a message based on the current timestamp
$messageIndex = time() % count($messages);
$response = [
    'message' => $messages[$messageIndex]
];

$delay = 5000; // 5 seconds

// URL encode the messages and create the response
$response = implode('/', array_map('urlencode', $messages)) . "?i=$delay&d=1";

echo $response;
?>