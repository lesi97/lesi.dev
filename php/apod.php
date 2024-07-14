<?php

/*
ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);
*/


date_default_timezone_set("America/Denver");

header("Cache-Control: no-store, no-cache, must-revalidate, max-age=0");
header("Cache-Control: post-check=0, pre-check=0", false);
header("Pragma: no-cache");

	include '../../private/loadEnv.php';

    $date = date('d-m-Y');
    $yesterday = (new DateTime())->modify('-1 day')->format('d-m-Y');

    $yesterdayCacheFile = "../assets/data/nasa_apod_cache-{$yesterday}.json";
    $yesterdayHDImagePath = "../assets/images/space/apodHD-{$yesterday}.jpg";
    $yesterdayImagePath = "../assets/images/space/apod-{$yesterday}.jpg";

    if (file_exists($yesterdayCacheFile)) unlink($yesterdayCacheFile);
    if (file_exists($yesterdayHDImagePath)) unlink($yesterdayHDImagePath);
    if (file_exists($yesterdayImagePath)) unlink($yesterdayImagePath);
    
    $cacheFile = "../assets/data/nasa_apod_cache-{$date}.json";
    $hdImagePath = "../assets/images/space/apodHD-{$date}.jpg";
    $imagePath = "../assets/images/space/apod-{$date}.jpg";
    $api_key = $_ENV["NASA_KEY"];
    $nasaDomain = "https://api.nasa.gov/planetary/apod";
    $parameters = "?api_key=";
    $url = $nasaDomain . $parameters . $api_key;
	
	$backupApodPhoto = "../assets/images/space/backup-apod/ngc6960_Pugh_960.jpg";
	$backupApodPhotoHD = "../assets/images/space/backup-apod/ngc6960_Pugh_2000.jpg";
    
    if (file_exists($cacheFile)) {
        $cache = json_decode(file_get_contents($cacheFile), true);
        $cacheDate = $cache["date"];
        
        if ($cacheDate == date("d-m-Y")) {
            echo json_encode([$cache]);
            exit;
        }
    }

    $curl = curl_init();
    
    curl_setopt_array($curl, [
        CURLOPT_URL => $url,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => "",
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 30,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => "GET",
        CURLOPT_HTTPHEADER => [
            "Accept: */*"
        ],
    ]);
    
    $response = curl_exec($curl);
    $err = curl_error($curl);
    
    curl_close($curl);
    
   if ($err) {
       echo "cURL Error #:" . $err;
   } else {
       $responseArray = json_decode($response, true);
	   
		if ($responseArray["media_type"] !== "video") {

			$apodImageUrl = $responseArray["url"];
			$apodHdImageUrl = $responseArray["hdurl"];

			$apodImageCurl = curl_init();
			curl_setopt_array($apodImageCurl, [
				CURLOPT_URL => $apodImageUrl,
				CURLOPT_RETURNTRANSFER => true,
				CURLOPT_ENCODING => "",
				CURLOPT_MAXREDIRS => 10,
				CURLOPT_TIMEOUT => 30,
				CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
				CURLOPT_CUSTOMREQUEST => "GET",
				CURLOPT_HTTPHEADER => [
					"Accept: image/jpeg"
				],
			]);
			$apodImageResponse = curl_exec($apodImageCurl);
			$apodImageErr = curl_error($apodImageCurl);
			curl_close($apodImageCurl);

			$apodHdImageCurl = curl_init();
			curl_setopt_array($apodHdImageCurl, [
				CURLOPT_URL => $apodHdImageUrl,
				CURLOPT_RETURNTRANSFER => true,
				CURLOPT_ENCODING => "",
				CURLOPT_MAXREDIRS => 10,
				CURLOPT_TIMEOUT => 30,
				CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
				CURLOPT_CUSTOMREQUEST => "GET",
				CURLOPT_HTTPHEADER => [
					"Accept: image/jpeg"
				],
			]);

			$apodHdImageResponse = curl_exec($apodHdImageCurl);
			$apodHdImageErr = curl_error($apodHdImageCurl);
			curl_close($apodHdImageCurl);
       
			file_put_contents($imagePath, $apodImageResponse);
			file_put_contents($hdImagePath, $apodHdImageResponse); 

      
			$cacheData = json_encode([
				'url' => $imagePath,
				'hdurl' => $hdImagePath,
				'date' => $date,
				'response' => json_encode($responseArray)
			]);

			if (file_exists($yesterdayCacheFile)) unlink($yesterdayCacheFile);
			if (file_exists($yesterdayHDImagePath)) unlink($yesterdayHDImagePath);
			if (file_exists($yesterdayImagePath)) unlink($yesterdayImagePath);
       
			file_put_contents($cacheFile, $cacheData);
			echo json_encode($responseArray);
	} else { 			
			$altApod = [
				"date" => $date,
				"explanation" => "Ten thousand years ago, before the dawn of recorded human history, a new light would have suddenly have appeared in the night sky and faded after a few weeks. Today we know this light was from a supernova, or exploding star, and record the expanding debris cloud as the Veil Nebula, a supernova remnant. This sharp telescopic view is centered on a western segment of the Veil Nebula cataloged as NGC 6960 but less formally known as the Witch's Broom Nebula. Blasted out in the cataclysmic explosion, the interstellar shock wave plows through space sweeping up and exciting interstellar material. Imaged with narrow band filters, the glowing filaments are like long ripples in a sheet seen almost edge on, remarkably well separated into atomic hydrogen (red) and oxygen (blue-green) gas. The complete supernova remnant lies about 1400 light-years away towards the constellation Cygnus. This Witch's Broom actually spans about 35 light-years. The bright star in the frame is 52 Cygni, visible with the unaided eye from a dark location but unrelated to the ancient supernova remnant.",
				"media_type" => "image",
				"service_version" => "v1",
				"title" => "NGC 6960: The Witch's Broom Nebula",
				"copyright" => "Martin Pugh (Heaven's Mirror Observatory)",
				"url" => "https://apod.nasa.gov/apod/image/1804/ngc6960_Pugh_960.jpg",
				"hdurl" => "https://apod.nasa.gov/apod/image/1804/ngc6960_Pugh_2000.jpg"
			];
			
			copy($backupApodPhoto, $imagePath);
			copy($backupApodPhotoHD, $hdImagePath);
	
			$cacheData = json_encode([
				'url' => $imagePath,
				'hdurl' => $hdImagePath,
				'date' => $date,
				'response' => json_encode($altApod)
			]);
			file_put_contents($cacheFile, $cacheData);
			echo json_encode($altApod);
	};
}


?>