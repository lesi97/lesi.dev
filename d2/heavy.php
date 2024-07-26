<?php

//set_time_limit(120);
/*
ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);
*/

//	Variables

	include '../../private/loadEnv.php';
	$api_key = $_ENV['BUNGIE_KEY'];
	$user = $_GET['user'];
	$platform = isset($_GET['platform']) ? $_GET['platform'] : null;
	$bungieEndpoint = 'https://www.bungie.net/Platform/';
	$endpointType = 'Destiny2/';
	$membershipType;
	$destinyMembershipId;

	$warlockHash = "2271682572";
	$titanHash = "3655393761";
	$hunterHash = "671679327";

	$warlock = "";
	$titan = "";
	$hunter = "";

	if (isset($user)) {

		$user = rawurlencode($user);
			
		if ($platform === "xb") {
			$platform = "1";
		} else if ($platform === "ps") {
			$platform = "2";
		} else if ($platform === "pc") {
			$platform = "3";
		} else if ($platform === "bnet") {
			$platform = "4";
		} else if ($platform === "st") {
			$platform = "5";
		} else if ($platform === "demon") {
			$platform = "10";
		} else {
			$platform = "-1";
		}

		// Get Membership ID From Bungie

		
		$membershipIdUrl = $bungieEndpoint . $endpointType . "SearchDestinyPlayer/" . $platform . "/" . $user . "/";

		$getMembershipId = curl_init();
		curl_setopt($getMembershipId, CURLOPT_URL, $membershipIdUrl);
		curl_setopt($getMembershipId, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($getMembershipId, CURLOPT_FOLLOWLOCATION, true);
		curl_setopt($getMembershipId, CURLOPT_HTTPHEADER, array(
			'x-api-key: ' . $api_key
		));
		$membershipIdResponse = curl_exec($getMembershipId);
		
		if (curl_errno($getMembershipId)) {
			echo 'cURL error: ' . curl_error($getMembershipId);
		}
		
		curl_close($getMembershipId);

		if (isset($membershipIdResponse["ErrorStatus"]) && $membershipIdResponse["ErrorStatus"] === "SystemDisabled") {
			echo 'bungie api is currently down, gift a sub instead ';
			return;
		}
		
		if ($membershipIdResponse !== false) {
			$membershipIdData = json_decode($membershipIdResponse, true);

			if ($membershipIdData["Response"][0]["isPublic"] === false) {
				echo $_GET['user'] . " ur account is private dummy, go here https://www.bungie.net/7/en/User/Account/Privacy then gift a sub";
				return;
			}

			if (json_last_error() == JSON_ERROR_NONE) {
				$destinyMembershipId = $membershipIdData["Response"][0]["membershipId"];
				$membershipType = $membershipIdData["Response"][0]["membershipType"];
			} else {
				echo 'bungie api is currently down, gift a sub instead ';
				return;
			}
		}
	


		// Get Character ID's

		$charactersIdsUrl = $bungieEndpoint . $endpointType . $membershipType . "/Profile/" . $destinyMembershipId . "/?components=200";

		$getCharacterIds = curl_init();
		curl_setopt($getCharacterIds, CURLOPT_URL, $charactersIdsUrl);
		curl_setopt($getCharacterIds, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($getCharacterIds, CURLOPT_FOLLOWLOCATION, true);
		curl_setopt($getCharacterIds, CURLOPT_HTTPHEADER, array(
			'x-api-key: ' . $api_key
		));
		$characterIdsResponse = curl_exec($getCharacterIds);
		
		if (curl_errno($getCharacterIds)) {
			echo 'cURL error: ' . curl_error($getCharacterIds);
		}
		
		curl_close($getCharacterIds);
		
		if ($characterIdsResponse !== false) {
			$characterIdsData = json_decode($characterIdsResponse, true);
			if (json_last_error() == JSON_ERROR_NONE) {
				if (isset($characterIdsData["Response"]["characters"]["data"])) {
					$charactersData = $characterIdsData["Response"]["characters"]["data"];
		
					$mostRecentCharacters = [
						'warlock' => ['id' => null, 'dateLastPlayed' => null],
						'titan' => ['id' => null, 'dateLastPlayed' => null],
						'hunter' => ['id' => null, 'dateLastPlayed' => null],
					];
		
					foreach ($charactersData as $character) {
						if (isset($character['classHash'], $character['characterId'], $character['dateLastPlayed'])) {
							$classHash = $character['classHash'];
							$characterId = $character['characterId'];
							$dateLastPlayed = $character['dateLastPlayed'];
		
							if ($classHash == $warlockHash) {
								if ($mostRecentCharacters['warlock']['dateLastPlayed'] === null || $dateLastPlayed > $mostRecentCharacters['warlock']['dateLastPlayed']) {
									$mostRecentCharacters['warlock'] = ['id' => $characterId, 'dateLastPlayed' => $dateLastPlayed];
								}
							} elseif ($classHash == $titanHash) {
								if ($mostRecentCharacters['titan']['dateLastPlayed'] === null || $dateLastPlayed > $mostRecentCharacters['titan']['dateLastPlayed']) {
									$mostRecentCharacters['titan'] = ['id' => $characterId, 'dateLastPlayed' => $dateLastPlayed];
								}
							} elseif ($classHash == $hunterHash) {
								if ($mostRecentCharacters['hunter']['dateLastPlayed'] === null || $dateLastPlayed > $mostRecentCharacters['hunter']['dateLastPlayed']) {
									$mostRecentCharacters['hunter'] = ['id' => $characterId, 'dateLastPlayed' => $dateLastPlayed];
								}
							}
						}
					}
		
					$warlock = $mostRecentCharacters['warlock']['id'];
					$titan = $mostRecentCharacters['titan']['id'];
					$hunter = $mostRecentCharacters['hunter']['id'];
				} else {
					echo 'No character data available';
					return;
				}
			} else {
				echo 'bungie api is currently down, gift a sub instead';
				return;
			}
		}
		


		

		$characters = '200';
		$characterEquipment = '205';
		$itemPerks = '302';
		$itemSockets = '305';
		$itemPlugObjectives = '309';
		$components = "?components=" . $characterEquipment . "," . $itemPlugObjectives . "," . $characters . "," . $itemSockets;
		
		$crucibleTracker = '38912240';
		
		$url = $bungieEndpoint . $endpointType . $membershipType . "/Profile/" . $destinyMembershipId . $components;
		
		$manifestUrl = $bungieEndpoint . $endpointType . "Manifest/DestinyInventoryItemDefinition/";
		
		$pvpTrackerHash = "";
		$pveTrackerHash = "";
		
		$weaponData = array(
			"weaponName" => "",
			"weaponPerks" => "",
			"exotic" => false,
			"pvpTracker" => array(
				"enabled" => false,
				"killCount" => "",
			),
			"pveTracker" => array(
				"enabled" => false,
				"killCount" => "",
			),
			"weaponMod" => array(
				"enabled" => false,
				"name" => "",
			),
			"weaponOrnament" => array(
				"enabled" => false,
				"name" => "",
			),
			"weaponShader" => array(
				"enabled" => false,
				"name" => "",
			)
		);
	
	
	
////////////////////////////////////////////////////////////////////////////////////

		$getData = curl_init();
		curl_setopt($getData, CURLOPT_URL, $url);
		curl_setopt($getData, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($getData, CURLOPT_FOLLOWLOCATION, true);
		curl_setopt($getData, CURLOPT_HTTPHEADER, array(
			'x-api-key: ' . $api_key
		));
		$bungieResponse = curl_exec($getData);
		
		if (curl_errno($getData)) {
			echo 'cURL error: ' . curl_error($getData);
		}
		
		curl_close($getData);
		
		if ($bungieResponse !== false) {
			$data = json_decode($bungieResponse, true);
			if (json_last_error() == JSON_ERROR_NONE) {
				
				
				if (isset($data["ErrorStatus"]) && $data["ErrorStatus"] !== "SystemDisabled") {
					
					
					$warlockLastPlayed = date($data["Response"]["characters"]["data"][$warlock]["dateLastPlayed"]);
					$hunterLastPlayed = date($data["Response"]["characters"]["data"][$hunter]["dateLastPlayed"]);
					$titanLastPlayed = date($data["Response"]["characters"]["data"][$titan]["dateLastPlayed"]);
					
					$mostRecentCharacter = $warlockLastPlayed;
					$currentCharacter = $warlock;
					
					if ($hunterLastPlayed > $mostRecentCharacter) {
						$mostRecentCharacter = $hunterLastPlayed;
						$currentCharacter = $hunter;
					}
					if ($titanLastPlayed > $mostRecentCharacter) {
						$mostRecentCharacter = $titanLastPlayed;
						$currentCharacter = $titan;
					}

					
					
					
					// Get Kinetic Gun Name Hash And Perk Hashes - 0 Is Kinetic, 1 For Energy - 2 For Heavy
					$itemInstanceId = $data["Response"]["characterEquipment"]["data"][$currentCharacter]["items"]["2"]["itemInstanceId"];	
					$itemHash = $data["Response"]["characterEquipment"]["data"][$currentCharacter]["items"]["2"]["itemHash"];			
					$perkHashes = $data["Response"]["itemComponents"]["sockets"]["data"][$itemInstanceId]["sockets"];
					
					// Get Perk Hashes From $perkHashes Array
					if (is_array($perkHashes)) {
						$perkHash = []; 
						for ($i = 0; $i < count($perkHashes); $i++) {
							if (isset($perkHashes[$i]["plugHash"])) {
								$perkHash[$i] = $perkHashes[$i]["plugHash"];
							} else {
								error_log("Error: Missing 'plugHash' key in \$perkHashes at index $i");
							}
						}
					} else {
						$perkHash = []; 
						error_log("Error: \$perkHashes is null or not an array");
					}


					
					// Get Kinetic Gun name
					$chWepUrl = $manifestUrl . $itemHash . "/";
					$chWep = curl_init();
					curl_setopt($chWep, CURLOPT_URL, $chWepUrl);
					curl_setopt($chWep, CURLOPT_RETURNTRANSFER, true);
					curl_setopt($chWep, CURLOPT_HTTPHEADER, array(
						'x-api-key: ' . $api_key
					));
					$chWepResponse = curl_exec($chWep);
					curl_close($chWep);
					if ($chWepResponse !== false) {
						$chWepData = json_decode($chWepResponse, true);
						$weaponData["weaponName"] = $chWepData["Response"]["displayProperties"]["name"];
					}			
					
					$newData = array();
					$perks = array();
					
					
					
					// Get Perk Names And Put Into $perks Array
					for ($i = 0; $i < count($perkHashes); $i++) {
						$manifestNewUrl[$i] = $manifestUrl . $perkHash[$i] . '/';
						$ch_i = curl_init();
						curl_setopt($ch_i, CURLOPT_URL, $manifestNewUrl[$i]);
						curl_setopt($ch_i, CURLOPT_RETURNTRANSFER, true);
						curl_setopt($ch_i, CURLOPT_HTTPHEADER, array(
							'x-api-key: ' . $api_key
						));
						$execResult = curl_exec($ch_i);	
						if ($execResult !== false) {		
							$newData[$i] = json_decode($execResult, true);
							$perks[$i] = $newData[$i]["Response"]["displayProperties"]["name"];
							if ($perks[$i] == "Crucible Tracker") {
								$pvpTrackerHash = $perkHash[$i];
								$weaponData["pvpTracker"]["enabled"] = true;
								$weaponData["pvpTracker"]["killCount"] = number_format($data["Response"]["itemComponents"]["plugObjectives"]["data"][$itemInstanceId]["objectivesPerPlug"][$pvpTrackerHash]["0"]["progress"]);
							} elseif ($perks[$i] == "Kill Tracker") {
								$pveTrackerHash = $perkHash[$i];
								$weaponData["pveTracker"]["enabled"] = true;
								$weaponData["pveTracker"]["killCount"] = number_format($data["Response"]["itemComponents"]["plugObjectives"]["data"][$itemInstanceId]["objectivesPerPlug"][$pveTrackerHash]["0"]["progress"]);
							}
							if ($newData[$i]["Response"]["itemTypeDisplayName"] === "Weapon Ornament") {
								$weaponData["weaponOrnament"]["enabled"] = true;
								$weaponData["weaponOrnament"]["name"] = $newData[$i]["Response"]["displayProperties"]["name"];
							}
							if ($newData[$i]["Response"]["itemTypeAndTierDisplayName"] === "Exotic Weapon Ornament") {
								$weaponData["exotic"] = true;
							}
							if ($newData[$i]["Response"]["itemTypeDisplayName"] === "Shader") {
								if ($newData[$i]["Response"]["displayProperties"]["name"] !== "Default Shader") {
									$weaponData["weaponShader"]["enabled"] = true;
									$weaponData["weaponShader"]["name"] = $newData[$i]["Response"]["displayProperties"]["name"];
								}
							}		
							if (isset($newData[$i]["Response"]["itemTypeDisplayName"]) && $newData[$i]["Response"]["itemTypeDisplayName"] === "Weapon Mod") {
								$weaponData["weaponMod"]["enabled"] = true;
								$weaponData["weaponMod"]["name"] = $newData[$i]["Response"]["displayProperties"]["name"];
							}
						} else {
							echo "Curl error: " . curl_error($ch_i);
						}	
						curl_close($ch_i);
					}

					

					

					// Only Get The Perks People Care About
					$selectedPerks = array_slice($perks, 0, 5);
					$weaponData["weaponPerks"] = implode(", ", $selectedPerks);
					
					$echoString = $_GET['user'] . " " . $weaponData["weaponName"] . " | Perks: " .  $weaponData["weaponPerks"];
					
					if ($weaponData["weaponMod"]["enabled"]) {
						$echoString .= " | Mod: " . $weaponData["weaponMod"]["name"];
					}
					if ($weaponData["weaponOrnament"]["enabled"]) {
						$echoString .= " | Ornament: " . $weaponData["weaponOrnament"]["name"];
					}
					if ($weaponData["weaponShader"]["enabled"]) {
						$echoString .= " | Shader: " . $weaponData["weaponShader"]["name"];
					} elseif (!$weaponData["exotic"]) {
						$echoString .= " | Shader: None";
					}
					if ($weaponData["pvpTracker"]["enabled"]) {
						$echoString .= " | PVP Kill Count: " . $weaponData["pvpTracker"]["killCount"];
					}
					if ($weaponData["pveTracker"]["enabled"]) {
						$echoString .= " | PVE Kill Count: " . $weaponData["pveTracker"]["killCount"];
					}
					echo $echoString;
					
				} else {
					echo "bungie api is currently down, gift a sub instead ";
					return;
				}		
				
			} else {
				echo ' ';
			}
		}			

	} else {
		echo "You have to include a valid user. Example below where the %23 is a replacement for #:<br>";
		echo "https://lesi.dev/d2/heavy?user=Lesi%235934<br><br>";
		echo "You can also specific the platform by appending the URL with any of the following:<br>";
		echo "Xbox:&emsp;&emsp;&emsp; &platform=xb<br>";
		echo "Playstation:&emsp;&platform=ps<br>";
		echo "Steam:&emsp;&emsp;&emsp;&platform=pc<br><br>";
		echo "Example below:<br>";
		echo "https://lesi.dev/d2/heavy?user=Lesi%235934&platform=pc";

}

?>