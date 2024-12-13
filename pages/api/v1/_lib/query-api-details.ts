import { supabase } from '@/lib/server/supabase-client';
import { PostgrestError } from '@supabase/postgrest-js';
import { encrypt, decrypt } from '@/lib/server/encryption';

export async function getApiDetails(application: string) {
    const { data, error } = await supabase
        .from('api_keys')
        .select(`client_id, client_secret, access_token, refresh_token, refresh_token_expiry, base_url, redirect_url`)
        .ilike('name', application);
    if (error) {
        return { data: null, error };
    }
    const unencryptedData = data.map((record: any) => ({
        client_id: record.client_id ? decrypt(record.client_id) : null,
        client_secret: record.client_secret ? decrypt(record.client_secret) : null,
        access_token: record.access_token ? decrypt(record.access_token) : null,
        refresh_token: record.refresh_token ? decrypt(record.refresh_token) : null,
        refresh_token_expiry: record.refresh_token_expiry,
        base_url: record.base_url,
        redirect_url: record.redirect_url,
    }));

    return { data: unencryptedData, error: null };
}

export async function updateApiDetails(
    application: string,
    client_id: string | undefined,
    client_secret: string | undefined,
    access_token: string | undefined,
    refresh_token: string | undefined,
    refresh_token_expiry: number | undefined,
    base_url: string | undefined,
    redirect_url: string | undefined
) {
    const dataToUpdate: {
        name: string;
        client_id?: string;
        client_secret?: string;
        access_token?: string;
        refresh_token?: string;
        refresh_token_expiry?: number;
        base_url?: string;
        redirect_url?: string;
    } = {
        name: application,
    };

    if (access_token !== undefined) {
        dataToUpdate.access_token = encrypt(access_token);
    }
    if (refresh_token !== undefined) {
        dataToUpdate.refresh_token = encrypt(refresh_token);
    }
    if (refresh_token_expiry !== undefined) {
        dataToUpdate.refresh_token_expiry = refresh_token_expiry;
    }
    if (client_id !== undefined) {
        dataToUpdate.client_id = encrypt(client_id);
    }
    if (client_secret !== undefined) {
        dataToUpdate.client_secret = encrypt(client_secret);
    }
    if (base_url !== undefined) {
        dataToUpdate.base_url = base_url;
    }
    if (redirect_url !== undefined) {
        dataToUpdate.redirect_url = redirect_url;
    }

    const { error } = await supabase.from('api_keys').upsert(dataToUpdate, { onConflict: 'name' });

    if (error) {
        console.error('Error upserting record:', error);
    }
}
