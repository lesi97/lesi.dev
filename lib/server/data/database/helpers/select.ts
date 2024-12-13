import { supabase } from '@/lib/server/supabase-client';
import * as Database from '../database.types';
import { PostgrestError } from '@supabase/supabase-js';

export async function selectSql(
    table: string,
    columns: string,
    filters: Database.FilterArray
): Promise<[Record<string, string>[], null] | [null, PostgrestError] | [null, string]> {
    const query = supabase.from(table).select(columns);

    if (filters) {
        Object.entries(filters).forEach(([key, { value, filter_type }]) => {
            if (filter_type === 'ilike') {
                query.ilike(key, value);
            } else if (filter_type === 'like') {
                query.like(key, value);
            } else if (filter_type === 'in') {
                query.in(key, value);
            } else if (filter_type === 'gte') {
                query.gte(key, value);
            } else if (filter_type === 'lte') {
                query.lte(key, value);
            } else {
                query.eq(key, value);
            }
        });
    }

    const { data, error } = await query;
    if (error) return [null, error];
    if (!data || data.length === 0) return [null, 'No Data'];
    return [data, null];
}
