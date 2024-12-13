import { supabase } from '@/lib/server/supabase-client';
import { PostgrestError } from '@supabase/supabase-js';

export async function insertSql(table: string, data: Record<string, string>): Promise<[PostgrestError] | [null]> {
    const { error } = await supabase.from(table).insert(data);
    if (error) return [error];
    return [null];
}
