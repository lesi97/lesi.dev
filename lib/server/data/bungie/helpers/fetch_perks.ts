import { getApiDetails } from '../../database/database';
import type { ApiDetailsResponse } from '@/types/db-queries';
import { upsertSql } from '../../database/helpers';

export async function fetchPerks() {
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

    async function processPerks(jsonData) {
        const perks = [];

        for (const key of Object.keys(jsonData)) {
            const item = jsonData[key];
            // if (!isWeaponPerk(item)) return;
            // if (item?.itemTypeDisplayName?.includes('Enhanced')) console.log(item);
            // console.log(item);
            perks.push({
                name: item?.displayProperties?.name || 'Unknown',
                description: item?.displayProperties?.description || 'Unknown',
                item_type: item?.itemTypeDisplayName || 'Unknown',
                hash_id: item.hash,
            });
        }

        const table = 'destiny_weapon_perks';
        const conflict = 'hash_id';
        try {
            const [error] = await upsertSql(table, perks, conflict);
            if (error) {
                console.error('Error upserting perks:', error.message);
            } else {
                console.log('Perk data upserted successfully.');
            }
        } catch (err) {
            console.error('Error during database operation:', err.message);
        }
    }

    processPerks(data).catch((error) => console.error('Error processing weapons:', error));
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
        // manifestData.Response.jsonWorldComponentContentPaths.en.DestinySandboxPerkDefinition;
        manifestData.Response.jsonWorldComponentContentPaths.en.DestinyInventoryItemDefinition;
    // manifestData.Response.jsonWorldComponentContentPaths.en.DestinySocketTypeDefinition;
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
