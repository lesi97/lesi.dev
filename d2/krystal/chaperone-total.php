<?php

//	Variables

	include '../../../private/loadEnv.php';

	$api_key = $_ENV['BUNGIE_KEY'];

	$bungieEndpoint = 'https://www.bungie.net/Platform/';
	$endpointType = 'Destiny2/';

	$membershipType = '3/';
	$destinyMembershipId = '4611686018467506425/';
	$warlock = '2305843009410653169/';
	$hunter = '2305843009299388895/';
	$titan = '2305843009377320147/';
	
	$url = $bungieEndpoint . $endpointType . $membershipType . "Account/" . $destinyMembershipId . "Character/" . $hunter . "Stats/UniqueWeapons/";

	$weaponName = "chaperone";
	
	$jsonKeyName = $weaponName . "_total_kills";
	$jsonKeyNamePrecisionKillCount = $weaponName . "_total_headshots";
	$jsonKeyNamePrecisionPercentage = $weaponName . "_total_headshots_percentage";
	$jsonFileName = 'kill-counts.json';

	$jsonString = file_get_contents(__DIR__.'/../exotics.json');
	$data = json_decode($jsonString, true);
	$weaponId = $data['specialWeapons']['Shotguns']['chaperone'];


////////////////////////////////////////////////////////////////////////////////////

	$ch = curl_init();
	curl_setopt($ch, CURLOPT_URL, $url);
	curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
	curl_setopt($ch, CURLOPT_HTTPHEADER, array(
		'x-api-key: ' . $api_key
	));
	$response = curl_exec($ch);
	curl_close($ch);
	$current_dateTime = date("Y-m-d H:i:s");
	
	$uniqueWeaponKills = null;
	$uniqueWeaponPrecisionKills = null;
	$uniqueWeaponKillsPrecisionPercentage = null;
	
	if ($response !== false) {
		$data1 = json_decode($response, true);
		$uniqueWeaponsList = $data1["Response"]["weapons"];
		foreach ($uniqueWeaponsList as $weapon) {
			if ($weapon['referenceId'] == $weaponId) {
				$uniqueWeaponKills = $weapon['values']['uniqueWeaponKills']['basic']['value'];
				$uniqueWeaponPrecisionKills = $weapon['values']['uniqueWeaponPrecisionKills']['basic']['value'];
				$uniqueWeaponKillsPrecisionPercentage = $weapon['values']['uniqueWeaponKillsPrecisionKills']['basic']['displayValue'];
				$jsonData = file_get_contents($jsonFileName);
				$data2 = json_decode($jsonData, true);
				$data2[$jsonKeyName] = $uniqueWeaponKills;
				//$data2[$jsonKeyNamePrecisionKillCount] = $uniqueWeaponPrecisionKills;
				//$data2[$jsonKeyNamePrecisionPercentage] = $uniqueWeaponKillsPrecisionPercentage;
				$jsonData = json_encode($data2);
				file_put_contents($jsonFileName, $jsonData);
				echo number_format($uniqueWeaponKills);
			}
		}
	}
?>