import { supabase } from '../../../../../../lib/server/supabase-client';
import { Filter } from '@/types/db-queries';

export async function selectSql(table: string, columns: string, filters: Record<string, Filter>) {
    const query = supabase.from(table).select(columns);

    Object.entries(filters).forEach(([key, { value, filter_type = 'eq' }]) => {
        if (filter_type === 'ilike') {
            query.ilike(key, value);
        } else if (filter_type === 'like') {
            query.like(key, value);
        } else {
            query.eq(key, value);
        }
    });

    const { data, error } = await query;

    if (error) {
        return { data: null, error };
    }

    return { data: data || null, error: null };
}
