import { supabase } from '@/lib/server/supabase-client';
import { PostgrestError } from '@supabase/supabase-js';
import * as Database from '../database.types';

export async function upsertSql(
    table: string,
    data: Record<string, string | number | null>[] | Database.Api_Details,
    conflict: string
): Promise<[PostgrestError] | [null]> {
    const { error } = await supabase.from(table).upsert(data, { onConflict: conflict });

    if (error) return [error];
    return [null];
}
