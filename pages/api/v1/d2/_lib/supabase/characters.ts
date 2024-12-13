import { supabase } from '../../../../../../lib/server/supabase-client';

export async function getCharacters(membership_id: string) {
    const { data, error } = await supabase
        .from('destiny_user_characters')
        .select(`character_id, character_type`)
        .eq('membership_id', membership_id);

    if (error) {
        return { data: null, error };
    }
    return { data: data || null, error: null };
}

export async function updateCharacters(membership_id: string, character_id: string, character_type: string) {
    const dataToUpdate: { membership_id: string; character_id: string; character_type?: string } = {
        membership_id: membership_id,
        character_id: character_id,
    };

    if (character_type !== undefined) {
        dataToUpdate.character_type = character_type;
    }

    const { error } = await supabase
        .from('destiny_user_characters')
        .upsert(dataToUpdate, { onConflict: 'membership_id, character_id' });

    if (error) {
        console.error('Error upserting record:', error);
    }
}
