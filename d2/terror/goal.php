<?php

http_response_code(200);
header('Content-Type: text/plain');

$goal = $_GET['subs'];

if (is_numeric($goal) && (int)$goal == $goal) {
    echo $goal . " Gifted Remain To Secure 💫🖤 | https://x.com/ilyterror";
} else {
    echo "50 Gifted Remain To Secure 💫🖤 | https://x.com/ilyterror";
}

?>
