import { updateUser } from '../supabase/users';
import { updateCharacters } from '../supabase/characters';
import { getCharacterTypeEnum } from './get-character-type-enum';
import { getApiDetails } from '../../../_lib/query-api-details';
import type { Api_Details_Response } from '@/lib/server/data/database/database.types';

export async function getBungieUser(bungie_id: string, platformEnum: number) {
    try {
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

        const endpointType = 'Destiny2';
        const id = encodeURIComponent(bungie_id);

        const url = new URL(`${baseUrl}/${endpointType}/SearchDestinyPlayer/-1/${id}/`);

        const response = await fetch(url, {
            method: 'GET',
            headers: headersList,
        });

        const data = await response.json();
        const membership_id = data.Response[0].membershipId.toString();
        const friendly_name = data.Response[0].displayName;
        const preferred_platform = data.Response[0].membershipType;

        updateUser(membership_id, bungie_id, friendly_name, preferred_platform).catch((error) =>
            console.error('Database update failed:', error)
        );

        return data.Response[0];
    } catch (error) {
        console.error(error);
        return { error: 'bungie api is currently down, gift a sub instead' };
    }
}

export async function getBungieUserCharacters(membership_id: string, platformEnum: number) {
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

    const endpointType = 'Destiny2';

    const url = new URL(
        `${baseUrl}/${endpointType}/${platformEnum}/Profile/${membership_id}/?components=200,205,302,305,309`
    );

    const response = await fetch(url, {
        method: 'GET',
        headers: headersList,
    });

    const data = await response.json();
    const characters = data.Response.characters.data;

    Object.values(characters).forEach((character: any) => {
        const character_type = getCharacterTypeEnum(character.classType);
        updateCharacters(membership_id, character.characterId, character_type);
    });

    return data;
}
