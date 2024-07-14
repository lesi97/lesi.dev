<?php

header("Cache-Control: no-store, no-cache, must-revalidate, max-age=0");
header("Cache-Control: post-check=0, pre-check=0", false);
header("Pragma: no-cache");

	include '../../private/loadEnv.php';

	$home = $_ENV['HOME_IP'];

	$ipv4 = $_SERVER['REMOTE_ADDR'];

	$counterFile = "../assets/data/home_page_counter.json";

	if (!file_exists($counterFile)) {
		file_put_contents($counterFile, json_encode(['count' => 0]));
	}
	
	$data = json_decode(file_get_contents($counterFile), true);
	
	if ($ipv4 !== $home) {
		$data['count']++;
		file_put_contents($counterFile, json_encode($data));
	}
	
	echo $data['count'];

?>
