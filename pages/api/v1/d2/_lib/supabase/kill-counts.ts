import { supabase } from '../../../../../../lib/server/supabase-client';

export async function getKills(membership_id: string, weapon_id: string) {
    const { data, error } = await supabase
        .from('destiny_weapon_kill_counts')
        .select(
            `
      pvp_kills,
      pve_kills
    `
        )
        .eq('membership_id', membership_id)
        .eq('weapon_id', weapon_id);

    if (error) {
        return [{ error: error.message }, { status: 500 }];
    }
    return data[0] || [];
}

export async function updateKills(
    membership_id: string,
    weapon_id: string,
    weapon_hash: string,
    pvp_kills?: number,
    pve_kills?: number,
    trials_kills?: number
) {
    if (pvp_kills == null && pve_kills == null && trials_kills == null) {
        console.log('No kill values provided, skipping update.');
        return;
    }

    const dataToUpdate: {
        membership_id: string;
        weapon_id: string;
        weapon_hash: string;
        pvp_kills?: number;
        pve_kills?: number;
        trials_kills?: number;
    } = {
        membership_id: membership_id,
        weapon_id: weapon_id,
        weapon_hash: weapon_hash,
    };

    if (pvp_kills !== undefined) {
        dataToUpdate.pvp_kills = pvp_kills;
    }
    if (pve_kills !== undefined) {
        dataToUpdate.pve_kills = pve_kills;
    }
    if (trials_kills !== undefined) {
        dataToUpdate.trials_kills = trials_kills;
    }

    const { error } = await supabase
        .from('destiny_weapon_kill_counts')
        .upsert(dataToUpdate, { onConflict: 'membership_id, weapon_id' });

    if (error) {
        console.error('Error upserting record:', error);
    }
}
