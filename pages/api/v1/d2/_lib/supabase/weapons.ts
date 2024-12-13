import { supabase } from '../../../../../../lib/server/supabase-client';

export async function getWeaponName(id: string) {
    const { data, error } = await supabase.from('destiny_weapons').select(`display_name, tier_type_name`).eq('id', id);
    if (error) {
        return [{ error: error.message }, { status: 500 }];
    }
    return data[0] || '';
}
