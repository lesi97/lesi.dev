import { updatePerk } from '../supabase/weapon_perks';
import { getApiDetails } from '../../../_lib/query-api-details';
import type { Api_Details_Response } from '@/lib/server/data/database/database.types';

export async function getWeaponPerksManifest() {
    const apiDetails: Api_Details_Response = await getApiDetails('Bungie');

    if (apiDetails.error) {
        return console.error('Error fetching API details:', apiDetails.error.message);
    }

    const apiKey = apiDetails.data[0].client_id;
    const baseUrl = apiDetails.data[0].base_url;

    const headersList = {
        Accept: '*/*',
        'x-api-key': `${apiKey}`,
    };

    const url = await getManifestUrl();
    const response = await fetch(url, {
        method: 'GET',
        headers: headersList,
    });
    const data = await response.json();

    Object.values(data).forEach((obj: any) => {
        // if (obj.displayProperties.name.includes('Pocket Ace')) console.log(obj);
        if (isWeaponPerk(obj)) {
            const perkName = obj.displayProperties.name.replace("'", '');
            const perkDescription = obj.displayProperties.description.replace("'", '');
            const perkType = obj.itemTypeDisplayName.replace("'", '');
            const perkHash = obj.hash;

            updatePerk(perkName, perkDescription, perkType, perkHash);
        }
    });
    return;
}

async function getManifestUrl() {
    const apiDetails: Api_Details_Response = await getApiDetails('Bungie');

    if (apiDetails.error) {
        return console.error('Error fetching API details:', apiDetails.error.message);
    }

    const apiKey = apiDetails.data[0].client_id;
    const baseUrl = apiDetails.data[0].base_url;

    const headersList = {
        Accept: '*/*',
        'x-api-key': `${apiKey}`,
    };

    const manifestUrl = `${baseUrl}/Destiny2/Manifest/`;
    const manifestResponse = await fetch(manifestUrl, {
        method: 'GET',
        headers: headersList,
    });
    const manifestData = await manifestResponse.json();
    const destinyInventoryItemDefinition =
        manifestData.Response.jsonWorldComponentContentPaths.en.DestinyInventoryItemDefinition;
    const url = `${baseUrl?.replace('Platform', '')}/${destinyInventoryItemDefinition}`;

    return url;
}

function isWeaponPerk(obj: any): boolean {
    if (obj.itemType !== 19) return false;

    const perkIdentifiers = ['frames', 'perks', 'magazines', 'shader'];
    const isPerk = obj.perks && obj.perks.length > 0;
    const hasPerkIdentifier = perkIdentifiers.includes(obj.plug?.plugCategoryIdentifier);
    const isPerkType =
        obj.itemTypeDisplayName?.includes('Trait') ||
        obj.itemTypeDisplayName?.includes('Perk') ||
        obj.itemTypeDisplayName?.includes('Shader') ||
        obj.itemTypeDisplayName?.includes('Weapon Mod') ||
        obj.itemTypeDisplayName?.includes('Weapon Ornament') ||
        obj.itemTypeDisplayName?.includes('Magazine');

    const exclusions = ['Armor', 'Grenade', 'Transmat', 'Ghost Mod', 'Clan Perk', 'Enriching Tonic', 'Artifact Perk'];
    const isExcludedType = exclusions.some((exclusion) => obj.itemTypeDisplayName?.includes(exclusion));

    return (isPerk || hasPerkIdentifier || isPerkType) && !isExcludedType;
}
