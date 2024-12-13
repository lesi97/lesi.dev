import { getApiDetails } from '../../database/database';
import type { ApiDetailsResponse } from '@/types/db-queries';
import { upsertSql } from '../../database/helpers';

export async function fetchWeapons() {
    const [apiDetails, api_query_error] = await getApiDetails('Bungie');

    if (api_query_error) {
        return console.error('Error fetching API details:', api_query_error.message);
    }
    const apiKey = apiDetails.client_id;

    const headersList = {
        Accept: '*/*',
        'x-api-key': `${apiKey}`,
    };

    const url: string = await getManifestUrl();
    const response = await fetch(url, {
        method: 'GET',
        headers: headersList,
    });
    const data = await response.json();

    async function processWeapons(jsonData) {
        const weapons = [];

        Object.keys(jsonData).forEach((key) => {
            const item = jsonData[key];

            if (item.itemType === 3) {
                weapons.push({
                    id: item.hash,
                    display_name: item.displayProperties?.name || 'Unknown',
                    item_type_display_name: item.itemTypeDisplayName || '',
                    flavor_text: item.flavorText || '',
                    bucket_type_hash: item.inventory?.bucketTypeHash || null,
                    tier_type_hash: item.inventory?.tierTypeHash || null,
                    tier_type_name: item.inventory?.tierTypeName || '',
                    tier_type: item.inventory?.tierType || null,
                    talent_grid_hash: item.talentGrid?.talentGridHash || null,
                });
            }
        });

        const table = 'destiny_weapons';
        const conflict = 'id';
        await upsertSql(table, weapons, conflict);

        // console.log('Weapons data upserted successfully');
    }

    processWeapons(data).catch((error) => console.error('Error processing weapons:', error));
}

async function getManifestUrl(): string {
    const [apiDetails, api_query_error] = await getApiDetails('Bungie');

    if (api_query_error) {
        return console.error('Error fetching API details:', api_query_error.message);
    }

    const apiKey = apiDetails.client_id;
    const baseUrl = apiDetails.base_url;

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
