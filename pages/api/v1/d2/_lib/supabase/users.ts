import { supabase } from '../../../../../../lib/server/supabase-client';
import { getCharacters } from './characters';

export async function getUser(bungie_id: string) {
    const { data, error } = await supabase
        .from('destiny_users')
        .select(
            `
            id,
            membership_id,
            bungie_id,
            preferred_platform,
            friendly_name
        `
        )
        .ilike('bungie_id', bungie_id);

    if (error || data.length === 0) {
        return { data: null, error };
    }

    const userCharacters = await getCharacters(data[0].membership_id);
    const profile = data[0];
    profile.characters = userCharacters.data;

    return { data: profile || null, error: null };
}

export async function updateUser(
    membership_id: string,
    bungie_id: string,
    friendly_name: string,
    preferred_platform: number
) {
    const dataToUpdate: {
        membership_id: string;
        bungie_id: string;
        friendly_name?: string;
        preferred_platform?: number;
    } = {
        membership_id: membership_id,
        bungie_id: bungie_id,
    };

    if (friendly_name !== undefined) {
        dataToUpdate.friendly_name = friendly_name;
    }
    if (preferred_platform !== undefined) {
        dataToUpdate.preferred_platform = preferred_platform;
    }

    const { error } = await supabase
        .from('destiny_users')
        .upsert(dataToUpdate, { onConflict: 'membership_id, bungie_id' });

    if (error) {
        console.error('Error upserting record:', error);
    }
}
