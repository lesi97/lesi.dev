import { supabase } from '../../../../../../lib/server/supabase-client';

export async function getPerkName(hash_id: string) {
    const { data, error } = await supabase
        .from('destiny_weapon_perks')
        .select(`name, item_type`)
        .eq('hash_id', hash_id);

    if (error || !data[0]) {
        // getWeaponPerksManifest();
        return null;
    }
    return data[0];
}

export async function updatePerk(name: string, description: string, item_type: string, hash_id: string) {
    const dataToUpdate: { name: string; description: string; item_type: string; hash_id: number } = {
        name: name,
        description: description,
        item_type: item_type,
        hash_id: hash_id,
    };
    const { error } = await supabase.from('destiny_weapon_perks').upsert(dataToUpdate, { onConflict: 'hash_id' });

    if (error) {
        console.error('Error upserting record:', error);
    }
}
