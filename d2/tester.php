<?php
ini_set('memory_limit', '256M');

ini_set('display_errors', 1);
ini_set('display_startup_errors', 1);
error_reporting(E_ALL);

include '../../private/loadEnv.php';

$api_key = $_ENV['BUNGIE_KEY'];
$user = $_GET['user'];
$platform = $_GET['platform'] ?? null;

$bungieEndpoint = 'https://www.bungie.net/Platform/';
$endpointType = 'Destiny2/';
$membershipType = null;
$destinyMembershipId = null;

$warlockHash = "2271682572";
$titanHash = "3655393761";
$hunterHash = "671679327";

$warlock = "";
$titan = "";
$hunter = "";

if ($user) {

    $user = rawurlencode($user);

    $platformMapping = [
        "xb" => "1",
        "ps" => "2",
        "pc" => "3",
        "bnet" => "4",
        "st" => "5",
        "demon" => "10"
    ];

    $platform = $platformMapping[$platform] ?? "-1";

    // Get Membership ID From Bungie
    $membershipIdUrl = "{$bungieEndpoint}{$endpointType}SearchDestinyPlayer/{$platform}/{$user}/";
    $membershipIdResponse = fetchFromBungieApi($membershipIdUrl, $api_key);

    if ($membershipIdResponse && isset($membershipIdResponse["Response"][0])) {
        if ($membershipIdResponse["Response"][0]["isPublic"] === false) {
            echo "$_GET[user] ur account is private dummy, go here https://www.bungie.net/7/en/User/Account/Privacy then gift a sub";
            return;
        }

        $destinyMembershipId = $membershipIdResponse["Response"][0]["membershipId"];
        $membershipType = $membershipIdResponse["Response"][0]["membershipType"];
    } else {
        echo 'bungie api is currently down, gift a sub instead';
        return;
    }

    // Get Character ID's
    $charactersIdsUrl = "{$bungieEndpoint}{$endpointType}{$membershipType}/Profile/{$destinyMembershipId}/?components=200";
    $characterIdsResponse = fetchFromBungieApi($charactersIdsUrl, $api_key);

    if ($characterIdsResponse && isset($characterIdsResponse["Response"]["characters"]["data"])) {
        $charactersData = $characterIdsResponse["Response"]["characters"]["data"];

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
        echo 'bungie api is currently down, gift a sub instead';
        return;
    }

    // Fetch character and equipment data in parallel
    $components = "200,205,302,305,309";
    $url = "{$bungieEndpoint}{$endpointType}{$membershipType}/Profile/{$destinyMembershipId}/?components={$components}";
    $bungieResponse = fetchFromBungieApi($url, $api_key);

    echo "Here " . $bungieResponse . "\n";

    if ($bungieResponse && isset($bungieResponse["Response"])) {
        $data = $bungieResponse;
        if (isset($data["ErrorStatus"]) && $data["ErrorStatus"] !== "SystemDisabled") {

            $mostRecentCharacter = getMostRecentCharacter([$warlock, $titan, $hunter], $data);
            $currentCharacter = $mostRecentCharacter['id'];

            $itemInstanceId = $data["Response"]["characterEquipment"]["data"][$currentCharacter]["items"]["0"]["itemInstanceId"];
            $itemHash = $data["Response"]["characterEquipment"]["data"][$currentCharacter]["items"]["0"]["itemHash"];
            $perkHashes = $data["Response"]["itemComponents"]["sockets"]["data"][$itemInstanceId]["sockets"];

            $perkHash = array_column($perkHashes, 'plugHash');

            $weaponData = fetchWeaponData($itemHash, $perkHash, $manifestUrl, $api_key, $data, $itemInstanceId);

            echo formatOutput($_GET['user'], $weaponData);

        } else {
            echo "bungie api is currently down, gift a sub instead ";
        }
    } else {
        echo " ";
    }

} else {
    echo "Test You have to include a valid user. Example below where the %23 is a replacement for #:<br>";
    echo "https://lesi.dev/d2/kinetic?user=Lesi%235934<br><br>";
    echo "You can also specify the platform by appending the URL with any of the following:<br>";
    echo "Xbox:&emsp;&emsp;&emsp; &platform=xb<br>";
    echo "PlayStation:&emsp;&platform=ps<br>";
    echo "Steam:&emsp;&emsp;&emsp;&platform=pc<br><br>";
    echo "Example below:<br>";
    echo "https://lesi.dev/d2/kinetic?user=Lesi%235934&platform=pc";
}

function fetchFromBungieApi($url, $api_key) {
    echo $url . "\n";
    $ch = curl_init();
    curl_setopt($ch, CURLOPT_URL, $url);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($ch, CURLOPT_FOLLOWLOCATION, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, array('x-api-key: ' . $api_key));
    $response = curl_exec($ch);
    if (curl_errno($ch)) {
        echo 'cURL error: ' . curl_error($ch);
        return false;
    }
    curl_close($ch);
    return json_decode($response, true);
}

function getMostRecentCharacter($characterIds, $data) {
    $mostRecent = ['id' => null, 'dateLastPlayed' => null];
    foreach ($characterIds as $characterId) {
        if (isset($data["Response"]["characters"]["data"][$characterId]["dateLastPlayed"])) {
            $dateLastPlayed = $data["Response"]["characters"]["data"][$characterId]["dateLastPlayed"];
            if ($mostRecent['dateLastPlayed'] === null || $dateLastPlayed > $mostRecent['dateLastPlayed']) {
                $mostRecent = ['id' => $characterId, 'dateLastPlayed' => $dateLastPlayed];
            }
        }
    }
    return $mostRecent;
}

function fetchWeaponData($itemHash, $perkHashes, $manifestUrl, $api_key, $data, $itemInstanceId) {
    $weaponData = [
        "weaponName" => "",
        "weaponPerks" => "",
        "exotic" => false,
        "pvpTracker" => ["enabled" => false, "killCount" => ""],
        "pveTracker" => ["enabled" => false, "killCount" => ""],
        "weaponMod" => ["enabled" => false, "name" => ""],
        "weaponOrnament" => ["enabled" => false, "name" => ""],
        "weaponShader" => ["enabled" => false, "name" => ""]
    ];

    // Fetch weapon name
    $chWepUrl = "{$manifestUrl}{$itemHash}/";
    $chWepResponse = fetchFromBungieApi($chWepUrl, $api_key);
    if ($chWepResponse) {
        $weaponData["weaponName"] = $chWepResponse["Response"]["displayProperties"]["name"];
    }

    // Fetch perk names
    $perks = [];
    foreach ($perkHashes as $perkHash) {
        $manifestNewUrl = "{$manifestUrl}{$perkHash}/";
        $newData = fetchFromBungieApi($manifestNewUrl, $api_key);
        if ($newData) {
            $perkName = $newData["Response"]["displayProperties"]["name"];
            $perks[] = $perkName;
            updateWeaponData($newData, $perkName, $perkHash, $weaponData, $data, $itemInstanceId);
        }
    }
    $weaponData["weaponPerks"] = implode(", ", array_slice($perks, 0, 5));
    return $weaponData;
}

function updateWeaponData($newData, $perkName, $perkHash, &$weaponData, $data, $itemInstanceId) {
    if ($perkName == "Crucible Tracker") {
        $weaponData["pvpTracker"]["enabled"] = true;
        $weaponData["pvpTracker"]["killCount"] = number_format($data["Response"]["itemComponents"]["plugObjectives"]["data"][$itemInstanceId]["objectivesPerPlug"][$perkHash]["0"]["progress"]);
    } elseif ($perkName == "Kill Tracker") {
        $weaponData["pveTracker"]["enabled"] = true;
        $weaponData["pveTracker"]["killCount"] = number_format($data["Response"]["itemComponents"]["plugObjectives"]["data"][$itemInstanceId]["objectivesPerPlug"][$perkHash]["0"]["progress"]);
    } elseif ($newData["Response"]["itemTypeDisplayName"] === "Weapon Ornament") {
        $weaponData["weaponOrnament"]["enabled"] = true;
        $weaponData["weaponOrnament"]["name"] = $newData["Response"]["displayProperties"]["name"];
    } elseif ($newData["Response"]["itemTypeAndTierDisplayName"] === "Exotic Weapon Ornament") {
        $weaponData["exotic"] = true;
    } elseif ($newData["Response"]["itemTypeDisplayName"] === "Shader" && $newData["Response"]["displayProperties"]["name"] !== "Default Shader") {
        $weaponData["weaponShader"]["enabled"] = true;
        $weaponData["weaponShader"]["name"] = $newData["Response"]["displayProperties"]["name"];
    } elseif ($newData["Response"]["itemTypeDisplayName"] === "Weapon Mod") {
        $weaponData["weaponMod"]["enabled"] = true;
        $weaponData["weaponMod"]["name"] = $newData["Response"]["displayProperties"]["name"];
    }
}

function formatOutput($user, $weaponData) {
    $output = "{$user} {$weaponData["weaponName"]} | Perks: {$weaponData["weaponPerks"]}";
    if ($weaponData["weaponMod"]["enabled"]) {
        $output .= " | Mod: {$weaponData["weaponMod"]["name"]}";
    }
    if ($weaponData["weaponOrnament"]["enabled"]) {
        $output .= " | Ornament: {$weaponData["weaponOrnament"]["name"]}";
    }
    if ($weaponData["weaponShader"]["enabled"]) {
        $output .= " | Shader: {$weaponData["weaponShader"]["name"]}";
    } elseif (!$weaponData["exotic"]) {
        $output .= " | Shader: None";
    }
    if ($weaponData["pvpTracker"]["enabled"]) {
        $output .= " | PVP Kill Count: {$weaponData["pvpTracker"]["killCount"]}";
    }
    if ($weaponData["pveTracker"]["enabled"]) {
        $output .= " | PVE Kill Count: {$weaponData["pveTracker"]["killCount"]}";
    }
    return $output;
}
?>
