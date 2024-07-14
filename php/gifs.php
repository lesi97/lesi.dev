<?php

header("Cache-Control: no-store, no-cache, must-revalidate, max-age=0");
header("Cache-Control: post-check=0, pre-check=0", false);
header("Pragma: no-cache");
header("Expires: Sat, 26 Jul 1997 05:00:00 GMT"); 

    $gifs = [
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/disappointed-disappointed-fan.gif",
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/disappointed-hercules.gif",
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/disappointed-sigh.gif",
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/new-girl-disappointed.gif",
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/rainWilson_Disappointed.gif",
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/shikimori_Disappointed.gif",
        "https://lesi.dev/assets/images/gif/zendeskDisappointed/steveCarell_Disappointed.gif"
    ];
    
    $randomGif = $gifs[array_rand($gifs)];
    
/*
    header('Content-Type: application/json');
    echo $randomGif;
  */  
    header('Content-Type: image/gif');
    readfile($randomGif);
    

?>